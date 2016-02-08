package drive

import (
    "fmt"
    "io"
    "os"
    "time"
    "path/filepath"
    "google.golang.org/api/drive/v3"
    "google.golang.org/api/googleapi"
)

type DownloadArgs struct {
    Out io.Writer
    Progress io.Writer
    Id string
    Path string
    Force bool
    Recursive bool
    Stdout bool
}

func (self *Drive) Download(args DownloadArgs) error {
    return self.download(args)
}

func (self *Drive) download(args DownloadArgs) error {
    f, err := self.service.Files.Get(args.Id).Fields("id", "name", "size", "mimeType", "md5Checksum").Do()
    if err != nil {
        return fmt.Errorf("Failed to get file: %s", err)
    }

    if isDir(f) && !args.Recursive {
        return fmt.Errorf("'%s' is a directory, use --recursive to download directories", f.Name)
    } else if isDir(f) && args.Recursive {
        return self.downloadDirectory(f, args)
    } else if isBinary(f) {
        return self.downloadBinary(f, args)
    } else if !args.Recursive {
        return fmt.Errorf("'%s' is a google document and must be exported, see the export command", f.Name)
    }

    return nil
}

func (self *Drive) downloadBinary(f *drive.File, args DownloadArgs) error {
    res, err := self.service.Files.Get(f.Id).Download()
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

    filename := filepath.Join(args.Path, f.Name)

    // Check if file exists
    if !args.Force && fileExists(filename) {
        return fmt.Errorf("File '%s' already exists, use --force to overwrite", filename)
    }

    // Ensure any parent directories exists
    if err = mkdir(filename); err != nil {
        return err
    }

    // Create new file
    outFile, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("Unable to create new file: %s", err)
    }

    // Close file on function exit
    defer outFile.Close()

    fmt.Fprintf(args.Out, "\nDownloading %s...\n", f.Name)
    started := time.Now()

    // Save file to disk
    bytes, err := io.Copy(outFile, srcReader)
    if err != nil {
        return fmt.Errorf("Failed saving file: %s", err)
    }

    // Calculate average download rate
    rate := calcRate(f.Size, started, time.Now())

    fmt.Fprintf(args.Out, "Downloaded '%s' at %s/s, total %s\n", filename, formatSize(rate, false), formatSize(bytes, false))

    //if deleteSourceFile {
    //    self.Delete(args.Id)
    //}
    return nil
}

func (self *Drive) downloadDirectory(parent *drive.File, args DownloadArgs) error {
    listArgs := listAllFilesArgs{
        query: fmt.Sprintf("'%s' in parents", parent.Id),
        fields: []googleapi.Field{"nextPageToken", "files(id,name)"},
    }
    files, err := self.listAllFiles(listArgs)
    if err != nil {
        return fmt.Errorf("Failed listing files: %s", err)
    }

    newPath := filepath.Join(args.Path, parent.Name)

    for _, f := range files {
        // Copy args and update changed fields
        newArgs := args
        newArgs.Path = newPath
        newArgs.Id = f.Id
        newArgs.Stdout = false

        err = self.download(newArgs)
        if err != nil {
            return err
        }
    }

    return nil
}

func isDir(f *drive.File) bool {
    return f.MimeType == DirectoryMimeType
}

func isBinary(f *drive.File) bool {
    return f.Md5Checksum != ""
}
