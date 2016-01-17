package drive

import (
    "fmt"
    "mime"
    "os"
    "path/filepath"
    "google.golang.org/api/drive/v3"
    "golang.org/x/net/context"
)

type UploadFileArgs struct {
    Path string
    Name string
    Parent string
    Mime string
    Recursive bool
    Stdin bool
    Share bool
}

func (self *Drive) Upload(args UploadFileArgs) {
    //if args.Stdin {
    //    self.uploadStdin()
    //}

    srcFile, err := os.Open(args.Path)
    errorF(err, "Failed to open file: %s", err)

    srcFileInfo, err := srcFile.Stat()
    errorF(err, "Failed to read file metadata: %s", err)

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

    // Set parent folder if provided
    if args.Parent != "" {
        dstFile.Parents = []string{args.Parent}
    }

    f, err := self.service.Files.Create(dstFile).ResumableMedia(context.Background(), srcFile, srcFileInfo.Size(), dstFile.MimeType).Do()
    errorF(err, "Failed to upload file: %s", err)

    fmt.Printf("Uploaded '%s' at %s, total %d\n", f.Name, "x/s", f.Size)
    //if args.Share {
    //    self.Share(TODO)
    //}
}
