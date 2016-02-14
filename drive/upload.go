package drive

import (
    "fmt"
    "mime"
    "os"
    "io"
    "time"
    "path/filepath"
    "google.golang.org/api/googleapi"
    "google.golang.org/api/drive/v3"
)

type UploadArgs struct {
    Out io.Writer
    Progress io.Writer
    Path string
    Name string
    Parents []string
    Mime string
    Recursive bool
    Share bool
    ChunkSize int64
}

func (self *Drive) Upload(args UploadArgs) error {
    if args.ChunkSize > intMax() - 1 {
        return fmt.Errorf("Chunk size is to big, max chunk size for this computer is %d", intMax() - 1)
    }

    return self.upload(args)
}

func (self *Drive) upload(args UploadArgs) error {
    info, err := os.Stat(args.Path)
    if err != nil {
        return fmt.Errorf("Failed stat file: %s", err)
    }

    if info.IsDir() && !args.Recursive {
        return fmt.Errorf("'%s' is a directory, use --recursive to upload directories", info.Name())
    } else if info.IsDir() {
        args.Name = ""
        return self.uploadDirectory(args)
    } else {
        _, err := self.uploadFile(args)
        return err
    }
}

func (self *Drive) uploadDirectory(args UploadArgs) error {
    srcFile, srcFileInfo, err := openFile(args.Path)
    if err != nil {
        return err
    }

    // Close file on function exit
    defer srcFile.Close()

    // Make directory on drive
    f, err := self.mkdir(MkdirArgs{
        Out: args.Out,
        Name: srcFileInfo.Name(),
        Parents: args.Parents,
        Share: args.Share,
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

        // Upload
        err = self.upload(newArgs)
        if err != nil {
            return err
        }
    }

    return nil
}

func (self *Drive) uploadFile(args UploadArgs) (*drive.File, error) {
    srcFile, srcFileInfo, err := openFile(args.Path)
    if err != nil {
        return nil, err
    }

    // Close file on function exit
    defer srcFile.Close()

    // Instantiate empty drive file
    dstFile := &drive.File{}

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
    srcReader := getProgressReader(srcFile, args.Progress, srcFileInfo.Size())

    fmt.Fprintf(args.Out, "\nUploading %s...\n", args.Path)
    started := time.Now()

    f, err := self.service.Files.Create(dstFile).Fields("id", "name", "size", "md5Checksum").Media(srcReader, chunkSize).Do()
    if err != nil {
        return nil, fmt.Errorf("Failed to upload file: %s", err)
    }

    // Calculate average upload rate
    rate := calcRate(f.Size, started, time.Now())

    fmt.Fprintf(args.Out, "[file] id: %s, md5: %s, name: %s\n", f.Id, f.Md5Checksum, f.Name)
    fmt.Fprintf(args.Out, "Uploaded '%s' at %s/s, total %s\n", f.Name, formatSize(rate, false), formatSize(f.Size, false))
    return f, nil
}

type UploadStreamArgs struct {
    Out io.Writer
    In io.Reader
    Name string
    Parents []string
    Mime string
    Share bool
    ChunkSize int64
}

func (self *Drive) UploadStream(args UploadStreamArgs) (err error) {
    if args.ChunkSize > intMax() - 1 {
        return fmt.Errorf("Chunk size is to big, max chunk size for this computer is %d", intMax() - 1)
    }

    // Instantiate empty drive file
    dstFile := &drive.File{Name: args.Name}

    // Set mime type if provided
    if args.Mime != "" {
        dstFile.MimeType = args.Mime
    }

    // Set parent folders
    dstFile.Parents = args.Parents

    // Chunk size option
    chunkSize := googleapi.ChunkSize(int(args.ChunkSize))

    f, err := self.service.Files.Create(dstFile).Media(args.In, chunkSize).Do()
    if err != nil {
        return fmt.Errorf("Failed to upload file: %s", err)
    }

    fmt.Fprintf(args.Out, "Uploaded '%s' at %s, total %d\n", f.Name, "x/s", f.Size)
    //if args.Share {
    //    self.Share(TODO)
    //}
    return
}

func openFile(path string) (*os.File, os.FileInfo, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, nil, fmt.Errorf("Failed to open file: %s", err)
    }

    info, err := f.Stat()
    if err != nil {
        return nil, nil, fmt.Errorf("Failed getting file metadata: %s", err)
    }

    return f, info, nil
}
