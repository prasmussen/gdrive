package drive

import (
    "fmt"
)

type UrlArgs struct {
    FileId string    
    DownloadUrl bool    
}

func (self *Drive) Url(args UrlArgs) {
    if args.DownloadUrl {
        fmt.Println(downloadUrl(args.FileId))
        return
    }
    fmt.Println(previewUrl(args.FileId))
}

func previewUrl(id string) string {
    return fmt.Sprintf("https://drive.google.com/uc?id=%s", id)
}

func downloadUrl(id string) string {
    return fmt.Sprintf("https://drive.google.com/uc?id=%s&export=download", id)
}
