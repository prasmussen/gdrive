package drive

import (
    "fmt"
    "google.golang.org/api/drive/v3"
)

type FileInfoArgs struct {
    Id string
    SizeInBytes bool
}

func (self *Drive) Info(args FileInfoArgs) (err error) {
    f, err := self.service.Files.Get(args.Id).Fields("id", "name", "size", "createdTime", "modifiedTime", "md5Checksum", "mimeType", "parents", "shared", "description").Do()
    if err != nil {
        return fmt.Errorf("Failed to get file: %s", err)
    }

    PrintFileInfo(PrintFileInfoArgs{
        File: f,
        SizeInBytes: args.SizeInBytes,
    })

    return
}

type PrintFileInfoArgs struct {
    File *drive.File
    SizeInBytes bool
}

func PrintFileInfo(args PrintFileInfoArgs) {
    f := args.File

    items := []kv{
        kv{"Id", f.Id},
        kv{"Name", f.Name},
        kv{"Description", f.Description},
        kv{"Mime", f.MimeType},
        kv{"Size", formatSize(f.Size, args.SizeInBytes)},
        kv{"Created", formatDatetime(f.CreatedTime)},
        kv{"Modified", formatDatetime(f.ModifiedTime)},
        kv{"Md5sum", f.Md5Checksum},
        kv{"Shared", formatBool(f.Shared)},
        kv{"Parents", formatList(f.Parents)},
    }

    for _, item := range items {
        if item.value() != "" {
            fmt.Printf("%s: %s\n", item.key(), item.value())
        }
    }
}
