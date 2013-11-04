package cli

import (
    "fmt"
    "os"
    "io"
    "path/filepath"
    "strings"
    "mime"
    "code.google.com/p/google-api-go-client/drive/v2"
    "../util"
    "../gdrive"
)

func List(d *gdrive.Drive, query, titleFilter string, maxResults int, sharedStatus bool) {
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
        fmt.Println(err)
        os.Exit(1)
    }

    items := make([]map[string]string, 0, 0)

    for _, f := range list.Items {
        // Skip files that dont have a download url (they are not stored on google drive)
        if f.DownloadUrl == "" {
            continue
        }

        items = append(items, map[string]string{
            "Id": f.Id,
            "Title": util.TruncateString(f.Title, 40),
            "Size": util.FileSizeFormat(f.FileSize),
            "Created": util.ISODateToLocal(f.CreatedDate),
        })
    }

    columnOrder := []string{"Id", "Title", "Size", "Created"}

    if sharedStatus {
        addSharedStatus(d, items)
        columnOrder = append(columnOrder, "Shared")
    }

    util.PrintColumns(items, columnOrder, 3)
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

func Info(d *gdrive.Drive, fileId string) {
    info, err := d.Files.Get(fileId).Do()
    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }
    printInfo(d, info)
}

func printInfo(d *gdrive.Drive, f *drive.File) {
    fields := map[string]string{
        "Id": f.Id,
        "Title": f.Title,
        "Description": f.Description,
        "Size": util.FileSizeFormat(f.FileSize),
        "Created": util.ISODateToLocal(f.CreatedDate),
        "Modified": util.ISODateToLocal(f.ModifiedDate),
        "Owner": strings.Join(f.OwnerNames, ", "),
        "Md5sum": f.Md5Checksum,
        "Shared": util.FormatBool(isShared(d, f.Id)),
    }

    order := []string{"Id", "Title", "Description", "Size", "Created", "Modified", "Owner", "Md5sum", "Shared"}
    util.Print(fields, order)
}

// Upload file to drive
func Upload(d *gdrive.Drive, input io.ReadCloser, title string, share bool) {
    // Use filename or 'untitled' as title if no title is specified
    if title == "" {
        if f, ok := input.(*os.File); ok && input != os.Stdin {
            title = filepath.Base(f.Name())
        } else {
            title = "untitled"
        }
    }

    var mimeType = mime.TypeByExtension(filepath.Ext(title))
    metadata := &drive.File{
      Title: title, 
      MimeType: mimeType,
    }
    getRate := util.MeasureTransferRate()

    info, err := d.Files.Insert(metadata).Media(input).Do()
    if err != nil {
        fmt.Printf("An error occurred uploading the document: %v\n", err)
        os.Exit(1)
    }

    // Total bytes transferred
    bytes := info.FileSize

    // Print information about uploaded file
    printInfo(d, info)
    fmt.Printf("MIME Type: %s\n", mimeType)
    fmt.Printf("Uploaded '%s' at %s, total %s\n", info.Title, getRate(bytes), util.FileSizeFormat(bytes))

    // Share file if the share flag was provided
    if share {
        Share(d, info.Id)
    }
}

func DownloadLatest(d *gdrive.Drive, stdout bool) {
    list, err := d.Files.List().Do()

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if len(list.Items) == 0 {
        fmt.Println("No files found")
        os.Exit(1)
    }

    latestId := list.Items[0].Id
    Download(d, latestId, stdout, true)
}

// Download file from drive
func Download(d *gdrive.Drive, fileId string, stdout, deleteAfterDownload bool) {
    // Get file info
    info, err := d.Files.Get(fileId).Do()
    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }

    if info.DownloadUrl == "" {
        // If there is no DownloadUrl, there is no body
        fmt.Println("An error occurred: File is not downloadable")
        os.Exit(1)
    }

    // Measure transfer rate
    getRate := util.MeasureTransferRate()

    // GET the download url
    res, err := d.Client().Get(info.DownloadUrl)
    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }

    // Close body on function exit
    defer res.Body.Close()

    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }

    // Write file content to stdout
    if stdout {
        io.Copy(os.Stdout, res.Body)
        return
    }

    // Check if file exists
    if util.FileExists(info.Title) {
        fmt.Printf("An error occurred: '%s' already exists\n", info.Title)
        os.Exit(1)
    }

    // Create a new file
    outFile, err := os.Create(info.Title) 
    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }

    // Close file on function exit
    defer outFile.Close()

    // Save file to disk
    bytes, err := io.Copy(outFile, res.Body)
    if err != nil {
        fmt.Printf("An error occurred: %v")
        os.Exit(1)
    }

    fmt.Printf("Downloaded '%s' at %s, total %s\n", info.Title, getRate(bytes), util.FileSizeFormat(bytes))

    if deleteAfterDownload {
        Delete(d, fileId)
    }
}

// Delete file with given file id
func Delete(d *gdrive.Drive, fileId string) {
    info, err := d.Files.Get(fileId).Do()
    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }

    if err = d.Files.Delete(fileId).Do(); err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Removed file '%s'\n", info.Title)
}

// Make given file id readable by anyone -- auth not required to view/download file
func Share(d *gdrive.Drive, fileId string) {
    info, err := d.Files.Get(fileId).Do()
    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }

    perm := &drive.Permission{
        Value: "me",
        Type: "anyone",
        Role: "reader",
    }

    _, err = d.Permissions.Insert(fileId, perm).Do()
    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("File '%s' is now readable by everyone @ %s\n", info.Title, util.PreviewUrl(fileId))
}

// Removes the 'anyone' permission -- auth will be required to view/download file
func Unshare(d *gdrive.Drive, fileId string) {
    info, err := d.Files.Get(fileId).Do()
    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }

    err = d.Permissions.Delete(fileId, "anyone").Do()
    if err != nil {
        fmt.Printf("An error occurred: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("File '%s' is no longer shared to 'anyone'\n", info.Title)
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
