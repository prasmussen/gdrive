package drive

import (
    "fmt"
)

type DeleteArgs struct {
    Id string
}

func (self *Drive) Delete(args DeleteArgs) {
    f, err := self.service.Files.Get(args.Id).Fields("name").Do()
    errorF(err, "Failed to get file: %s", err)

    err = self.service.Files.Delete(args.Id).Do()
    errorF(err, "Failed to delete file")
    fmt.Printf("Removed file '%s'\n", f.Name)
}
