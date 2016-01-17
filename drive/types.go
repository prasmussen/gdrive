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
    NameWidth int64
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

type UploadFileArgs struct {
    Path string
    Name string
    Parent string
    Mime string
    Recursive bool
    Stdin bool
    Share bool
}

type FileInfoArgs struct {
    Id string
    SizeInBytes bool
}

type MkdirArgs struct {
    Name string
    Parent string
    Share bool
}

type PrintFileListArgs struct {
    Files []*drive.File
    NameWidth int
    SkipHeader bool
    SizeInBytes bool
}

type PrintFileInfoArgs struct {
    File *drive.File
    SizeInBytes bool
}

type kv [2]string

func (self kv) key() string {
    return self[0]
}

func (self kv) value() string {
    return self[1]
}
