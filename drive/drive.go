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
