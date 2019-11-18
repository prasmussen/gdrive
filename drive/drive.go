package drive

import (
	"net/http"
	"sync"

	"google.golang.org/api/drive/v3"
)

type Drive struct {
	service     *drive.Service
	waitGroup   sync.WaitGroup
	downloadErr error
}

func New(client *http.Client) (*Drive, error) {
	service, err := drive.New(client)
	if err != nil {
		return nil, err
	}

	return &Drive{service, sync.WaitGroup{}, nil}, nil
}
