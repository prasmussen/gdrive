package drive

import (
	"fmt"
	"io"
	"google.golang.org/api/drive/v3"
)

type RenameArgs struct {
	Out       io.Writer
	Id        string
	NewName   string
}

func (self *Drive) Rename(args RenameArgs) error {
	f, err := self.service.Files.Get(args.Id).Fields("name", "mimeType").Do()
	if err != nil {
		return fmt.Errorf("Failed to get file: %s", err)
	}

	f2, err := self.service.Files.Update(args.Id,&drive.File{Name:args.NewName}).Do()
	if err != nil {
		return fmt.Errorf("Failed to rename file: %s", err)
	}

	fmt.Fprintf(args.Out, "Renamed '%s' -> '%s'\n", f.Name, f2.Name)
	return nil
}
