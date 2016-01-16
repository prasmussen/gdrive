package drive

import (
    "fmt"
    "io"
    "mime"
    "os"
    "path/filepath"
    "google.golang.org/api/drive/v3"
    "golang.org/x/net/context"
)

func (self *Drive) List(args ListFilesArgs) {
    fileList, err := self.service.Files.List().PageSize(args.MaxFiles).Q(args.Query).Fields("nextPageToken", "files(id,name,size,createdTime)").Do()
    if err != nil {
        exitF("Failed listing files: %s\n", err.Error())
    }

    for _, f := range fileList.Files {
        fmt.Printf("%s %s %d %s\n", f.Id, f.Name, f.Size, f.CreatedTime)
    }
}


func (self *Drive) Download(args DownloadFileArgs) {
    getFile := self.service.Files.Get(args.Id)

    f, err := getFile.Do()
    if err != nil {
        exitF("Failed to get file: %s", err.Error())
    }

    res, err := getFile.Download()
    if err != nil {
        exitF("Failed to download file: %s", err.Error())
    }

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
    if err != nil {
        exitF("Unable to create new file: %s", err.Error())
    }

    // Close file on function exit
    defer outFile.Close()

    // Save file to disk
    bytes, err := io.Copy(outFile, res.Body)
    if err != nil {
        exitF("Failed saving file: %s", err.Error())
    }

    fmt.Printf("Downloaded '%s' at %s, total %d\n", f.Name, "x/s", bytes)

    //if deleteSourceFile {
    //    self.Delete(args.Id)
    //}
}

func (self *Drive) Upload(args UploadFileArgs) {
    //if args.Stdin {
    //    self.uploadStdin()
    //}

    srcFile, err := os.Open(args.Path)
    if err != nil {
        exitF("Failed to open file: %s", err.Error())
    }

    srcFileInfo, err := srcFile.Stat()
    if err != nil {
        exitF("Failed to read file metadata: %s", err.Error())
    }

    // Instantiate empty drive file
    dstFile := &drive.File{}

    // Use provided file name or use filename
    if args.Name == "" {
        dstFile.Name = filepath.Base(srcFileInfo.Name())
    } else {
        dstFile.Name = args.Name
    }

    // Set provided mime type or get type based on file extension
    if args.Mime == "" {
        dstFile.MimeType = mime.TypeByExtension(filepath.Ext(dstFile.Name))
    } else {
        dstFile.MimeType = args.Mime
    }

    // Set parent folder if provided
    if args.Parent != "" {
        dstFile.Parents = []string{args.Parent}
    }

    f, err := self.service.Files.Create(dstFile).ResumableMedia(context.Background(), srcFile, srcFileInfo.Size(), dstFile.MimeType).Do()
    if err != nil {
        exitF("Failed to upload file: %s", err.Error())
    }

    fmt.Printf("Uploaded '%s' at %s, total %d\n", f.Name, "x/s", f.Size)
    //if args.Share {
    //    self.Share(TODO)
    //}
}

//func newFile(args UploadFileArgs) *drive.File {
//
//}
