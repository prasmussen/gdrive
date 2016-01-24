package drive

import (
    "fmt"
    "io"
    "os"
)

type DownloadRevisionArgs struct {
    Out io.Writer
    Progress io.Writer
    FileId string
    RevisionId string
    Force bool
    Stdout bool
}

func (self *Drive) DownloadRevision(args DownloadRevisionArgs) (err error) {
    getRev := self.service.Revisions.Get(args.FileId, args.RevisionId)

    rev, err := getRev.Fields("originalFilename").Do()
    if err != nil {
        return fmt.Errorf("Failed to get file: %s", err)
    }

    if rev.OriginalFilename == "" {
        return fmt.Errorf("Download is not supported for this file type")
    }

    res, err := getRev.Download()
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
    if !args.Force && fileExists(rev.OriginalFilename) {
        return fmt.Errorf("File '%s' already exists, use --force to overwrite", rev.OriginalFilename)
    }

    // Create new file
    outFile, err := os.Create(rev.OriginalFilename)
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

    fmt.Fprintf(args.Out, "Downloaded '%s' at %s, total %d\n", rev.OriginalFilename, "x/s", bytes)

    //if deleteSourceFile {
    //    self.Delete(args.Id)
    //}
    return
}
