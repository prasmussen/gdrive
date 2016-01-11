package drive

import (
    "net/http"
    "google.golang.org/api/drive/v3"
)

type Client interface {
    Service() *drive.Service
    Http() *http.Client
}

type Drive struct {
    service *drive.Service
    http *http.Client
}

func NewDrive(client Client) *Drive {
    return &Drive{
        service: client.Service(),
        http: client.Http(),
    }
}

type ListFilesArgs struct {
    MaxFiles int64
    Query string
    SkipHeader bool
    SizeInBytes bool
}

type DownloadFileArgs struct {
    Id string
    Force bool
    NoProgress bool
    Stdout bool
}
