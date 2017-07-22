package drive

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type DeleteArgs struct {
	Out       io.Writer
	Id        string
	Recursive bool
	OutJSON   bool
}

func (self *Drive) Delete(args DeleteArgs) error {
	f, err := self.service.Files.Get(args.Id).Fields("name", "mimeType").Do()
	if err != nil {
		return fmt.Errorf("Failed to get file: %s", err)
	}

	if isDir(f) && !args.Recursive {
		return fmt.Errorf("'%s' is a directory, use the 'recursive' flag to delete directories", f.Name)
	}

	err = self.service.Files.Delete(args.Id).Do()
	if err != nil {
		return fmt.Errorf("Failed to delete file: %s", err)
	}
	if args.OutJSON {
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(f); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Fprintf(args.Out, "Deleted '%s'\n", f.Name)
	}
	return nil
}

func (self *Drive) deleteFile(fileId string) error {
	err := self.service.Files.Delete(fileId).Do()
	if err != nil {
		return fmt.Errorf("Failed to delete file: %s", err)
	}
	return nil
}
