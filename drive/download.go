package drive

import (
    "fmt"
    "io"
    "io/ioutil"
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

    // Discard other output if file is written to stdout
    out := args.Out
    if args.Stdout {
        out = ioutil.Discard
    }

    // Path to file
    fpath := filepath.Join(args.Path, f.Name)

    fmt.Fprintf(out, "Downloading %s -> %s\n", f.Name, fpath)

    bytes, rate, err := self.saveFile(saveFileArgs{
        out: args.Out,
        body: res.Body,
        contentLength: res.ContentLength,
        fpath: fpath,
        force: args.Force,
        stdout: args.Stdout,
        progress: args.Progress,
    })

    if err != nil {
        return err
    }

    fmt.Fprintf(out, "Download complete, rate: %s/s, total size: %s\n", formatSize(rate, false), formatSize(bytes, false))
    return nil
}

type saveFileArgs struct {
    out io.Writer
    body io.Reader
    contentLength int64
    fpath string
    force bool
    stdout bool
    progress io.Writer
}

func (self *Drive) saveFile(args saveFileArgs) (int64, int64, error) {
    // Wrap response body in progress reader
    srcReader := getProgressReader(args.body, args.progress, args.contentLength)

    if args.stdout {
        // Write file content to stdout
        _, err := io.Copy(args.out, srcReader)
        return 0, 0, err
    }

    // Check if file exists
    if !args.force && fileExists(args.fpath) {
        return 0, 0, fmt.Errorf("File '%s' already exists, use --force to overwrite", args.fpath)
    }

    // Ensure any parent directories exists
    if err := mkdir(args.fpath); err != nil {
        return 0, 0, err
    }

    // Download to tmp file
    tmpPath := args.fpath + ".incomplete"

    // Create new file
    outFile, err := os.Create(tmpPath)
    if err != nil {
        return 0, 0, fmt.Errorf("Unable to create new file: %s", err)
    }

    started := time.Now()

    // Save file to disk
    bytes, err := io.Copy(outFile, srcReader)
    if err != nil {
        outFile.Close()
        os.Remove(tmpPath)
        return 0, 0, fmt.Errorf("Failed saving file: %s", err)
    }

    // Calculate average download rate
    rate := calcRate(bytes, started, time.Now())

    //if deleteSourceFile {
    //    self.Delete(args.Id)
    //}

    // Close File
    outFile.Close()

    // Rename tmp file to proper filename
    return bytes, rate, os.Rename(tmpPath, args.fpath)
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
