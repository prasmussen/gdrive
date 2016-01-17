package drive

import (
    "fmt"
    "io"
    "os"
)

type DownloadFileArgs struct {
    Id string
    Force bool
    NoProgress bool
    Stdout bool
}

func (self *Drive) Download(args DownloadFileArgs) {
    getFile := self.service.Files.Get(args.Id)

    f, err := getFile.Do()
    errorF(err, "Failed to get file: %s", err)

    res, err := getFile.Download()
    errorF(err, "Failed to download file: %s", err)

    // Close body on function exit
    defer res.Body.Close()

    if args.Stdout {
        // Write file content to stdout
        io.Copy(os.Stdout, res.Body)
        return
    }

    // Check if file exists
    if !args.Force && fileExists(f.Name) {
        exitF("File '%s' already exists, use --force to overwrite", f.Name)
    }

    // Create new file
    outFile, err := os.Create(f.Name)
    errorF(err, "Unable to create new file: %s", err)

    // Close file on function exit
    defer outFile.Close()

    // Save file to disk
    bytes, err := io.Copy(outFile, res.Body)
    errorF(err, "Failed saving file: %s", err)

    fmt.Printf("Downloaded '%s' at %s, total %d\n", f.Name, "x/s", bytes)

    //if deleteSourceFile {
    //    self.Delete(args.Id)
    //}
}
