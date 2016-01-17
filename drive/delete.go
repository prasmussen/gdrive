package drive

import (
    "fmt"
)

type DeleteArgs struct {
    Id string
}

func (self *Drive) Delete(args DeleteArgs) (err error) {
    f, err := self.service.Files.Get(args.Id).Fields("name").Do()
    if err != nil {
        return fmt.Errorf("Failed to get file: %s", err)
    }

    err = self.service.Files.Delete(args.Id).Do()
    if err != nil {
        return fmt.Errorf("Failed to delete file", err)
    }

    fmt.Printf("Removed file '%s'\n", f.Name)
    return
}
