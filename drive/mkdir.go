package drive

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
)

const DirectoryMimeType = "application/vnd.google-apps.folder"

type MkdirArgs struct {
	Out         io.Writer
	Name        string
	Description string
	Parents     []string
}

func (self *Drive) Mkdir(args MkdirArgs) error {
	f, err := self.mkdir(args)
	if err != nil {
		return err
	}
	fmt.Fprintf(args.Out, "Directory %s created\n", f.Id)
	return nil
}

func (self *Drive) mkdir(args MkdirArgs) (*drive.File, error) {
	dstFile := &drive.File{
		Name:        args.Name,
		Description: args.Description,
		MimeType:    DirectoryMimeType,
	}

	// Set parent folders
	dstFile.Parents = args.Parents

	// Create directory
	f, err := self.service.Files.Create(dstFile).SupportsTeamDrives(true).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to create directory: %s", err)
	}

	return f, nil
}
