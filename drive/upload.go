package drive

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"io"
	"mime"
	"os"
	"path/filepath"
	"time"
)

type UploadArgs struct {
	Out         io.Writer
	Progress    io.Writer
	Path        string
	Name        string
	Description string
	Parents     []string
	Mime        string
	Recursive   bool
	Share       bool
	Delete      bool
	ChunkSize   int64
	Timeout     time.Duration
}

func (self *Drive) Upload(args UploadArgs) error {
	if args.ChunkSize > intMax()-1 {
		return fmt.Errorf("Chunk size is to big, max chunk size for this computer is %d", intMax()-1)
	}

	// Ensure that none of the parents are sync dirs
	for _, parent := range args.Parents {
		isSyncDir, err := self.isSyncFile(parent)
		if err != nil {
			return err
		}

		if isSyncDir {
			return fmt.Errorf("%s is a sync directory, use 'sync upload' instead", parent)
		}
	}

	if args.Recursive {
		return self.uploadRecursive(args)
	}

	info, err := os.Stat(args.Path)
	if err != nil {
		return fmt.Errorf("Failed stat file: %s", err)
	}

	if info.IsDir() {
		return fmt.Errorf("'%s' is a directory, use --recursive to upload directories", info.Name())
	}

	f, rate, err := self.uploadFile(args)
	if err != nil {
		return err
	}
	fmt.Fprintf(args.Out, "Uploaded %s at %s/s, total %s\n", f.Id, formatSize(rate, false), formatSize(f.Size, false))

	if args.Share {
		err = self.shareAnyoneReader(f.Id)
		if err != nil {
			return err
		}

		fmt.Fprintf(args.Out, "File is readable by anyone at %s\n", f.WebContentLink)
	}

	if args.Delete {
		err = os.Remove(args.Path)
		if err != nil {
			return fmt.Errorf("Failed to delete file: %s", err)
		}
		fmt.Fprintf(args.Out, "Removed %s\n", args.Path)
	}

	return nil
}

func (self *Drive) uploadRecursive(args UploadArgs) error {
	info, err := os.Stat(args.Path)
	if err != nil {
		return fmt.Errorf("Failed stat file: %s", err)
	}

	if info.IsDir() {
		args.Name = ""
		return self.uploadDirectory(args)
	} else if info.Mode().IsRegular() {
		_, _, err := self.uploadFile(args)
		return err
	}

	return nil
}

func (self *Drive) uploadDirectory(args UploadArgs) error {
	srcFile, srcFileInfo, err := openFile(args.Path)
	if err != nil {
		return err
	}

	// Close file on function exit
	defer srcFile.Close()

	fmt.Fprintf(args.Out, "Creating directory %s\n", srcFileInfo.Name())
	// Make directory on drive
	f, err := self.mkdir(MkdirArgs{
		Out:         args.Out,
		Name:        srcFileInfo.Name(),
		Parents:     args.Parents,
		Description: args.Description,
	})
	if err != nil {
		return err
	}

	// Read files from directory
	names, err := srcFile.Readdirnames(0)
	if err != nil && err != io.EOF {
		return fmt.Errorf("Failed reading directory: %s", err)
	}

	for _, name := range names {
		// Copy args and set new path and parents
		newArgs := args
		newArgs.Path = filepath.Join(args.Path, name)
		newArgs.Parents = []string{f.Id}
		newArgs.Description = ""

		// Upload
		err = self.uploadRecursive(newArgs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Drive) uploadFile(args UploadArgs) (*drive.File, int64, error) {
	srcFile, srcFileInfo, err := openFile(args.Path)
	if err != nil {
		return nil, 0, err
	}

	// Close file on function exit
	defer srcFile.Close()

	// Instantiate empty drive file
	dstFile := &drive.File{Description: args.Description}

	// Use provided file name or use filename
	if args.Name == "" {
		dstFile.Name = filepath.Base(srcFileInfo.Name())
	} else {
		dstFile.Name = args.Name
	}

	// Set provided mime type or get type based on file extension
	if args.Mime == "" {
		dstFile.MimeType = mime.TypeByExtension(filepath.Ext(dstFile.Name))
	} else {
		dstFile.MimeType = args.Mime
	}

	// Set parent folders
	dstFile.Parents = args.Parents

	// Chunk size option
	chunkSize := googleapi.ChunkSize(int(args.ChunkSize))

	// Wrap file in progress reader
	progressReader := getProgressReader(srcFile, args.Progress, srcFileInfo.Size())

	// Wrap reader in timeout reader
	reader, ctx := getTimeoutReaderContext(progressReader, args.Timeout)

	fmt.Fprintf(args.Out, "Uploading %s\n", args.Path)
	started := time.Now()

	f, err := self.service.Files.Create(dstFile).SupportsTeamDrives(true).Fields("id", "name", "size", "md5Checksum", "webContentLink").Context(ctx).Media(reader, chunkSize).Do()
	if err != nil {
		if isTimeoutError(err) {
			return nil, 0, fmt.Errorf("Failed to upload file: timeout, no data was transferred for %v", args.Timeout)
		}
		return nil, 0, fmt.Errorf("Failed to upload file: %s", err)
	}

	// Calculate average upload rate
	rate := calcRate(f.Size, started, time.Now())

	return f, rate, nil
}

type UploadStreamArgs struct {
	Out         io.Writer
	In          io.Reader
	Name        string
	Description string
	Parents     []string
	Mime        string
	Share       bool
	ChunkSize   int64
	Progress    io.Writer
	Timeout     time.Duration
}

func (self *Drive) UploadStream(args UploadStreamArgs) error {
	if args.ChunkSize > intMax()-1 {
		return fmt.Errorf("Chunk size is to big, max chunk size for this computer is %d", intMax()-1)
	}

	// Instantiate empty drive file
	dstFile := &drive.File{Name: args.Name, Description: args.Description}

	// Set mime type if provided
	if args.Mime != "" {
		dstFile.MimeType = args.Mime
	}

	// Set parent folders
	dstFile.Parents = args.Parents

	// Chunk size option
	chunkSize := googleapi.ChunkSize(int(args.ChunkSize))

	// Wrap file in progress reader
	progressReader := getProgressReader(args.In, args.Progress, 0)

	// Wrap reader in timeout reader
	reader, ctx := getTimeoutReaderContext(progressReader, args.Timeout)

	fmt.Fprintf(args.Out, "Uploading %s\n", dstFile.Name)
	started := time.Now()

	f, err := self.service.Files.Create(dstFile).SupportsTeamDrives(true).Fields("id", "name", "size", "webContentLink").Context(ctx).Media(reader, chunkSize).Do()
	if err != nil {
		if isTimeoutError(err) {
			return fmt.Errorf("Failed to upload file: timeout, no data was transferred for %v", args.Timeout)
		}
		return fmt.Errorf("Failed to upload file: %s", err)
	}

	// Calculate average upload rate
	rate := calcRate(f.Size, started, time.Now())

	fmt.Fprintf(args.Out, "Uploaded %s at %s/s, total %s\n", f.Id, formatSize(rate, false), formatSize(f.Size, false))
	if args.Share {
		err = self.shareAnyoneReader(f.Id)
		if err != nil {
			return err
		}

		fmt.Fprintf(args.Out, "File is readable by anyone at %s\n", f.WebContentLink)
	}
	return nil
}
