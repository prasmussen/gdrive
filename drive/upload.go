package drive

import (
    "fmt"
    "mime"
    "os"
    "io"
    "path/filepath"
    "google.golang.org/api/googleapi"
    "google.golang.org/api/drive/v3"
)

type UploadFileArgs struct {
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

func (self *Drive) Upload(args UploadFileArgs) (err error) {
    if args.ChunkSize > intMax() - 1 {
        return fmt.Errorf("Chunk size is to big, max chunk size for this computer is %d", intMax() - 1)
    }

    srcFile, err := os.Open(args.Path)
    if err != nil {
        return fmt.Errorf("Failed to open file: %s", err)
    }

    srcFileInfo, err := srcFile.Stat()
    if err != nil {
        return fmt.Errorf("Failed to read file metadata: %s", err)
    }

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

    f, err := self.service.Files.Create(dstFile).Media(srcReader, chunkSize).Do()
    if err != nil {
        return fmt.Errorf("Failed to upload file: %s", err)
    }

    fmt.Fprintf(args.Out, "Uploaded '%s' at %s, total %d\n", f.Name, "x/s", f.Size)
    //if args.Share {
    //    self.Share(TODO)
    //}
    return
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
