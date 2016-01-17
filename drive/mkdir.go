package drive

import (
    "google.golang.org/api/drive/v3"
    "fmt"
)

const DirectoryMimeType = "application/vnd.google-apps.folder"

type MkdirArgs struct {
    Name string
    Parent string
    Share bool
}

func (self *Drive) Mkdir(args MkdirArgs) (err error) {
    dstFile := &drive.File{Name: args.Name, MimeType: DirectoryMimeType}

    // Set parent folder if provided
    if args.Parent != "" {
        dstFile.Parents = []string{args.Parent}
    }

    // Create folder
    f, err := self.service.Files.Create(dstFile).Do()
    if err != nil {
        return fmt.Errorf("Failed to create folder: %s", err)
    }

    PrintFileInfo(PrintFileInfoArgs{File: f})

    //if args.Share {
    //    self.Share(TODO)
    //}
    return
}
