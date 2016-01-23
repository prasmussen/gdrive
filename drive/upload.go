package drive

import (
    "fmt"
    "mime"
    "os"
    "io"
    "path/filepath"
    "google.golang.org/api/drive/v3"
    "golang.org/x/net/context"
)

type UploadFileArgs struct {
    Out io.Writer
    Path string
    Name string
    Parents []string
    Mime string
    Recursive bool
    Stdin bool
    Share bool
    NoProgress bool
}

func (self *Drive) Upload(args UploadFileArgs) (err error) {
    //if args.Stdin {
    //    self.uploadStdin()
    //}

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

    f, err := self.service.Files.Create(dstFile).ResumableMedia(context.Background(), srcFile, srcFileInfo.Size(), dstFile.MimeType).Do()
    if err != nil {
        return fmt.Errorf("Failed to upload file: %s", err)
    }

    fmt.Fprintf(args.Out, "Uploaded '%s' at %s, total %d\n", f.Name, "x/s", f.Size)
    //if args.Share {
    //    self.Share(TODO)
    //}
    return
}
