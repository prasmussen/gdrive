package cli

import (
	"fmt"
	"github.com/prasmussen/gdrive/gdrive"
	"github.com/prasmussen/gdrive/util"
	"github.com/prasmussen/google-api-go-client/drive/v2"
	"golang.org/x/net/context"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

func List(d *gdrive.Drive, query, titleFilter string, maxResults int, sharedStatus bool, noHeader bool) error {
	caller := d.Files.List()

	if maxResults > 0 {
		caller.MaxResults(int64(maxResults))
	}

	if titleFilter != "" {
		q := fmt.Sprintf("title contains '%s'", titleFilter)
		caller.Q(q)
	}

	if query != "" {
		caller.Q(query)
	}

	list, err := caller.Do()
	if err != nil {
		return err
	}

	items := make([]map[string]string, 0, 0)

	for _, f := range list.Items {
		// Skip files that dont have a download url (they are not stored on google drive)
		if f.DownloadUrl == "" {
			if f.MimeType != "application/vnd.google-apps.folder" {
				continue
			}
		}
		if f.Labels.Trashed {
			continue
		}

		items = append(items, map[string]string{
			"Id":      f.Id,
			"Title":   util.TruncateString(f.Title, 40),
			"Size":    util.FileSizeFormat(f.FileSize),
			"Created": util.ISODateToLocal(f.CreatedDate),
		})
	}

	columnOrder := []string{"Id", "Title", "Size", "Created"}

	if sharedStatus {
		addSharedStatus(d, items)
		columnOrder = append(columnOrder, "Shared")
	}

	util.PrintColumns(items, columnOrder, 3, noHeader)
	return nil
}

// Adds the key-value-pair 'Shared: True/False' to the map
func addSharedStatus(d *gdrive.Drive, items []map[string]string) {
	// Limit to 10 simultaneous requests
	active := make(chan bool, 10)
	done := make(chan bool)

	// Closure that performs the check
	checkStatus := func(item map[string]string) {
		// Wait for an empty spot in the active queue
		active <- true

		// Perform request
		shared := isShared(d, item["Id"])
		item["Shared"] = util.FormatBool(shared)

		// Decrement the active queue and notify that we are done
		<-active
		done <- true
	}

	// Go, go, go!
	for _, item := range items {
		go checkStatus(item)
	}

	// Wait for all goroutines to finish
	for i := 0; i < len(items); i++ {
		<-done
	}
}

func Info(d *gdrive.Drive, fileId string) error {
	info, err := d.Files.Get(fileId).Do()
	if err != nil {
		return fmt.Errorf("An error occurred: %v\n", err)
	}
	printInfo(d, info)
	return nil
}

func printInfo(d *gdrive.Drive, f *drive.File) {
	fields := map[string]string{
		"Id":          f.Id,
		"Title":       f.Title,
		"Description": f.Description,
		"Size":        util.FileSizeFormat(f.FileSize),
		"Created":     util.ISODateToLocal(f.CreatedDate),
		"Modified":    util.ISODateToLocal(f.ModifiedDate),
		"Owner":       strings.Join(f.OwnerNames, ", "),
		"Md5sum":      f.Md5Checksum,
		"Shared":      util.FormatBool(isShared(d, f.Id)),
		"Parents":     util.ParentList(f.Parents),
	}

	order := []string{
		"Id",
		"Title",
		"Description",
		"Size",
		"Created",
		"Modified",
		"Owner",
		"Md5sum",
		"Shared",
		"Parents",
	}
	util.Print(fields, order)
}

// Create folder in drive
func Folder(d *gdrive.Drive, title string, parentId string, share bool) error {
	info, err := makeFolder(d, title, parentId, share)
	if err != nil {
		return err
	}
	printInfo(d, info)
	fmt.Printf("Folder '%s' created\n", info.Title)
	return nil
}

func makeFolder(d *gdrive.Drive, title string, parentId string, share bool) (*drive.File, error) {
	// File instance
	f := &drive.File{Title: title, MimeType: "application/vnd.google-apps.folder"}
	// Set parent (if provided)
	if parentId != "" {
		p := &drive.ParentReference{Id: parentId}
		f.Parents = []*drive.ParentReference{p}
	}
	// Create folder
	info, err := d.Files.Insert(f).Do()
	if err != nil {
		return nil, fmt.Errorf("An error occurred creating the folder: %v\n", err)
	}
	// Share folder if the share flag was provided
	if share {
		Share(d, info.Id)
	}
	return info, err
}

// Upload file to drive
func UploadStdin(d *gdrive.Drive, input io.ReadCloser, title string, parentId string, share bool, mimeType string, convert bool) error {
	// File instance
	f := &drive.File{Title: "untitled"}
	// Set parent (if provided)
	if parentId != "" {
		p := &drive.ParentReference{Id: parentId}
		f.Parents = []*drive.ParentReference{p}
	}
	getRate := util.MeasureTransferRate()

	if convert {
		fmt.Printf("Converting to Google Docs format enabled\n")
	}

	info, err := d.Files.Insert(f).Convert(convert).Media(input).Do()
	if err != nil {
		return fmt.Errorf("An error occurred uploading the document: %v\n", err)
	}

	// Total bytes transferred
	bytes := info.FileSize

	// Print information about uploaded file
	printInfo(d, info)
	fmt.Printf("MIME Type: %s\n", mimeType)
	fmt.Printf("Uploaded '%s' at %s, total %s\n", info.Title, getRate(bytes), util.FileSizeFormat(bytes))

	// Share file if the share flag was provided
	if share {
		err = Share(d, info.Id)
	}
	return err
}

func Upload(d *gdrive.Drive, input *os.File, title string, parentId string, share bool, mimeType string, convert bool) error {
	// Grab file info
	inputInfo, err := input.Stat()
	if err != nil {
		return err
	}

	if inputInfo.IsDir() {
		return uploadDirectory(d, input, inputInfo, title, parentId, share, mimeType, convert)
	} else {
		return uploadFile(d, input, inputInfo, title, parentId, share, mimeType, convert)
	}

	return nil
}

func uploadDirectory(d *gdrive.Drive, input *os.File, inputInfo os.FileInfo, title string, parentId string, share bool, mimeType string, convert bool) error {
	// Create folder
	folder, err := makeFolder(d, filepath.Base(inputInfo.Name()), parentId, share)
	if err != nil {
		return err
	}

	// Read all files in directory
	files, err := input.Readdir(0)
	if err != nil {
		return err
	}

	// Get current dir
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Go into directory
	err = input.Chdir()
	if err != nil {
		return err
	}

	// Change back to original directory when done
	defer func() {
		os.Chdir(currentDir)
	}()

	for _, fi := range files {
		f, err := os.Open(fi.Name())
		if err != nil {
			return err
		}

		if fi.IsDir() {
			err = uploadDirectory(d, f, fi, "", folder.Id, share, mimeType, convert)
		} else {
			err = uploadFile(d, f, fi, "", folder.Id, share, mimeType, convert)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func uploadFile(d *gdrive.Drive, input *os.File, inputInfo os.FileInfo, title string, parentId string, share bool, mimeType string, convert bool) error {
	if title == "" {
		title = filepath.Base(inputInfo.Name())
	}

	if mimeType == "" {
		mimeType = mime.TypeByExtension(filepath.Ext(title))
	}

	// File instance
	f := &drive.File{Title: title, MimeType: mimeType}
	// Set parent (if provided)
	if parentId != "" {
		p := &drive.ParentReference{Id: parentId}
		f.Parents = []*drive.ParentReference{p}
	}
	getRate := util.MeasureTransferRate()

	if convert {
		fmt.Printf("Converting to Google Docs format enabled\n")
	}

	info, err := d.Files.Insert(f).Convert(convert).ResumableMedia(context.Background(), input, inputInfo.Size(), mimeType).Do()
	if err != nil {
		return fmt.Errorf("An error occurred uploading the document: %v\n", err)
	}

	// Total bytes transferred
	bytes := info.FileSize

	// Print information about uploaded file
	printInfo(d, info)
	fmt.Printf("MIME Type: %s\n", mimeType)
	fmt.Printf("Uploaded '%s' at %s, total %s\n", info.Title, getRate(bytes), util.FileSizeFormat(bytes))

	// Share file if the share flag was provided
	if share {
		err = Share(d, info.Id)
	}
	return err
}

func DownloadLatest(d *gdrive.Drive, stdout bool) error {
	list, err := d.Files.List().Do()
	if err != nil {
		return err
	}

	if len(list.Items) == 0 {
		return fmt.Errorf("No files found")
	}

	latestId := list.Items[0].Id
	return Download(d, latestId, stdout, true)
}

// Download file from drive
func Download(d *gdrive.Drive, fileId string, stdout, deleteAfterDownload bool) error {
	// Get file info
	info, err := d.Files.Get(fileId).Do()
	if err != nil {
		return fmt.Errorf("An error occurred: %v\n", err)
	}

	if info.DownloadUrl == "" {
		// If there is no DownloadUrl, there is no body
		return fmt.Errorf("An error occurred: File is not downloadable")
	}

	// Measure transfer rate
	getRate := util.MeasureTransferRate()

	// GET the download url
	res, err := d.Client().Get(info.DownloadUrl)
	if err != nil {
		return fmt.Errorf("An error occurred: %v\n", err)
	}

	// Close body on function exit
	defer res.Body.Close()

	// Write file content to stdout
	if stdout {
		io.Copy(os.Stdout, res.Body)
		return nil
	}

	// Check if file exists
	if util.FileExists(info.Title) {
		return fmt.Errorf("An error occurred: '%s' already exists\n", info.Title)
	}

	// Create a new file
	outFile, err := os.Create(info.Title)
	if err != nil {
		return fmt.Errorf("An error occurred: %v\n", err)
	}

	// Close file on function exit
	defer outFile.Close()

	// Save file to disk
	bytes, err := io.Copy(outFile, res.Body)
	if err != nil {
		return fmt.Errorf("An error occurred: %s", err)
	}

	fmt.Printf("Downloaded '%s' at %s, total %s\n", info.Title, getRate(bytes), util.FileSizeFormat(bytes))

	if deleteAfterDownload {
		err = Delete(d, fileId)
	}
	return err
}

// Delete file with given file id
func Delete(d *gdrive.Drive, fileId string) error {
	info, err := d.Files.Get(fileId).Do()
	if err != nil {
		return fmt.Errorf("An error occurred: %v\n", err)
	}

	if err := d.Files.Delete(fileId).Do(); err != nil {
		return fmt.Errorf("An error occurred: %v\n", err)

	}

	fmt.Printf("Removed file '%s'\n", info.Title)
	return nil
}

// Make given file id readable by anyone -- auth not required to view/download file
func Share(d *gdrive.Drive, fileId string) error {
	info, err := d.Files.Get(fileId).Do()
	if err != nil {
		return fmt.Errorf("An error occurred: %v\n", err)
	}

	perm := &drive.Permission{
		Value: "me",
		Type:  "anyone",
		Role:  "reader",
	}

	if _, err := d.Permissions.Insert(fileId, perm).Do(); err != nil {
		return fmt.Errorf("An error occurred: %v\n", err)
	}

	fmt.Printf("File '%s' is now readable by everyone @ %s\n", info.Title, util.PreviewUrl(fileId))
	return nil
}

// Removes the 'anyone' permission -- auth will be required to view/download file
func Unshare(d *gdrive.Drive, fileId string) error {
	info, err := d.Files.Get(fileId).Do()
	if err != nil {
		return fmt.Errorf("An error occurred: %v\n", err)
	}

	if err := d.Permissions.Delete(fileId, "anyone").Do(); err != nil {
		return fmt.Errorf("An error occurred: %v\n", err)
	}

	fmt.Printf("File '%s' is no longer shared to 'anyone'\n", info.Title)
	return nil
}

func isShared(d *gdrive.Drive, fileId string) bool {
	r, err := d.Permissions.List(fileId).Do()
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		os.Exit(1)
	}

	for _, perm := range r.Items {
		if perm.Type == "anyone" {
			return true
		}
	}
	return false
}
