package drive

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"path/filepath"
)

func (self *Drive) newPathfinder() *remotePathfinder {
	return &remotePathfinder{
		service: self.service.Files,
		files:   make(map[string]*drive.File),
	}
}

type remotePathfinder struct {
	service *drive.FilesService
	files   map[string]*drive.File
}

func (self *remotePathfinder) absPath(f *drive.File) (string, error) {
	name := f.Name

	if len(f.Parents) == 0 {
		return name, nil
	}

	var path []string

	for {
		parent, err := self.getParent(f.Parents[0])
		if err != nil {
			return "", err
		}

		// Stop when we find the root dir
		if len(parent.Parents) == 0 {
			break
		}

		path = append([]string{parent.Name}, path...)
		f = parent
	}

	path = append(path, name)
	return filepath.Join(path...), nil
}

func (self *remotePathfinder) getParent(id string) (*drive.File, error) {
	// Check cache
	if f, ok := self.files[id]; ok {
		return f, nil
	}

	// Fetch file from drive
	f, err := self.service.Get(id).SupportsTeamDrives(true).Fields("id", "name", "parents").Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to get file: %s", err)
	}

	// Save in cache
	self.files[f.Id] = f

	return f, nil
}
