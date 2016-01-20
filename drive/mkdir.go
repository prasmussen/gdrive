package drive

import (
    "google.golang.org/api/drive/v3"
    "io"
    "fmt"
)

const DirectoryMimeType = "application/vnd.google-apps.folder"

type MkdirArgs struct {
    Out io.Writer
    Name string
    Parents []string
    Share bool
}

func (self *Drive) Mkdir(args MkdirArgs) (err error) {
    dstFile := &drive.File{Name: args.Name, MimeType: DirectoryMimeType}

    // Set parent folders
    dstFile.Parents = args.Parents

    // Create folder
    f, err := self.service.Files.Create(dstFile).Do()
    if err != nil {
        return fmt.Errorf("Failed to create folder: %s", err)
    }

    PrintFileInfo(PrintFileInfoArgs{Out: args.Out, File: f})

    //if args.Share {
    //    self.Share(TODO)
    //}
    return
}
