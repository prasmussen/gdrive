package drive

import (
    "io"
    "fmt"
)

type UrlArgs struct {
    Out io.Writer
    FileId string    
    DownloadUrl bool    
}

func (self *Drive) Url(args UrlArgs) {
    if args.DownloadUrl {
        fmt.Fprintln(args.Out, downloadUrl(args.FileId))
        return
    }
    fmt.Fprintln(args.Out, previewUrl(args.FileId))
}

func previewUrl(id string) string {
    return fmt.Sprintf("https://drive.google.com/uc?id=%s", id)
}

func downloadUrl(id string) string {
    return fmt.Sprintf("https://drive.google.com/uc?id=%s&export=download", id)
}
