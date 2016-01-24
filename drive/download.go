package drive

import (
    "fmt"
    "io"
    "os"
)

type DownloadFileArgs struct {
    Out io.Writer
    Progress io.Writer
    Id string
    Force bool
    Stdout bool
}

func (self *Drive) Download(args DownloadFileArgs) (err error) {
    getFile := self.service.Files.Get(args.Id)

    f, err := getFile.Do()
    if err != nil {
        return fmt.Errorf("Failed to get file: %s", err)
    }

    res, err := getFile.Download()
    if err != nil {
        return fmt.Errorf("Failed to download file: %s", err)
    }

    // Close body on function exit
    defer res.Body.Close()

    // Wrap response body in progress reader
    srcReader := getProgressReader(res.Body, args.Progress, res.ContentLength)

    if args.Stdout {
        // Write file content to stdout
        _, err := io.Copy(args.Out, srcReader)
        return err
    }

    // Check if file exists
    if !args.Force && fileExists(f.Name) {
        return fmt.Errorf("File '%s' already exists, use --force to overwrite", f.Name)
    }

    // Create new file
    outFile, err := os.Create(f.Name)
    if err != nil {
        return fmt.Errorf("Unable to create new file: %s", err)
    }

    // Close file on function exit
    defer outFile.Close()

    // Save file to disk
    bytes, err := io.Copy(outFile, srcReader)
    if err != nil {
        return fmt.Errorf("Failed saving file: %s", err)
    }

    fmt.Fprintf(args.Out, "Downloaded '%s' at %s, total %d\n", f.Name, "x/s", bytes)

    //if deleteSourceFile {
    //    self.Delete(args.Id)
    //}
    return
}
