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

const DirectoryMimeType = "application/vnd.google-apps.folder"


func (self *Drive) List(args ListFilesArgs) {
    fileList, err := self.service.Files.List().PageSize(args.MaxFiles).Q(args.Query).Fields("nextPageToken", "files(id,name,size,createdTime)").Do()
    errorF(err, "Failed listing files: %s\n", err)

    PrintFileList(PrintFileListArgs{
        Files: fileList.Files,
        NameWidth: int(args.NameWidth),
        SkipHeader: args.SkipHeader,
        SizeInBytes: args.SizeInBytes,
    })
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

func (self *Drive) Upload(args UploadFileArgs) {
    //if args.Stdin {
    //    self.uploadStdin()
    //}

    srcFile, err := os.Open(args.Path)
    errorF(err, "Failed to open file: %s", err)

    srcFileInfo, err := srcFile.Stat()
    errorF(err, "Failed to read file metadata: %s", err)

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
    errorF(err, "Failed to upload file: %s", err)

    fmt.Printf("Uploaded '%s' at %s, total %d\n", f.Name, "x/s", f.Size)
    //if args.Share {
    //    self.Share(TODO)
    //}
}

func (self *Drive) Info(args FileInfoArgs) {
    f, err := self.service.Files.Get(args.Id).Fields("id", "name", "size", "createdTime", "modifiedTime", "md5Checksum", "mimeType", "parents", "shared", "description").Do()
    errorF(err, "Failed to get file: %s", err)

    PrintFileInfo(PrintFileInfoArgs{
        File: f,
        SizeInBytes: args.SizeInBytes,
    })
}

func (self *Drive) Mkdir(args MkdirArgs) {
    dstFile := &drive.File{Name: args.Name, MimeType: DirectoryMimeType}

    // Set parent folder if provided
    if args.Parent != "" {
        dstFile.Parents = []string{args.Parent}
    }

    // Create folder
    f, err := self.service.Files.Create(dstFile).Do()
    errorF(err, "Failed to create folder: %s", err)

    PrintFileInfo(PrintFileInfoArgs{File: f})

    //if args.Share {
    //    self.Share(TODO)
    //}
}
