package client

import (
    "net/http"
    "google.golang.org/api/drive/v3"
)

type Client struct {
    service *drive.Service
    http *http.Client
}

func (self *Client) Service() *drive.Service {
    return self.service
}

func (self *Client) Http() *http.Client {
    return self.http
}

func NewClient(client *http.Client) (*Client, error) {
    service, err := drive.New(client)
    if err != nil {
        return nil, err
    }

    return &Client{service, client}, nil
}
