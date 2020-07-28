package drive

import (
	"net/http"
	"sync"
	"time"

	"google.golang.org/api/drive/v3"
)

type Drive struct {
	service           *drive.Service
	downloadStartUnix int64
	downlaodCount     int64
	waitGroup         sync.WaitGroup
	downloadErr       error
}

func New(client *http.Client) (*Drive, error) {
	service, err := drive.New(client)
	if err != nil {
		return nil, err
	}

	return &Drive{service, 0, 0, sync.WaitGroup{}, nil}, nil
}

func (d *Drive) ResetDownloadTime() {
	d.downloadStartUnix = time.Now().Unix()
	d.downlaodCount = 0
}
