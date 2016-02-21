package drive

import (
	"google.golang.org/api/drive/v3"
	"net/http"
)

type Drive struct {
	service *drive.Service
}

func New(client *http.Client) (*Drive, error) {
	service, err := drive.New(client)
	if err != nil {
		return nil, err
	}

	return &Drive{service}, nil
}
