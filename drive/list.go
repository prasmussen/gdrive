package drive

import (
    "fmt"
    "io"
    "text/tabwriter"
    "google.golang.org/api/drive/v3"
)

type ListFilesArgs struct {
    Out io.Writer
    MaxFiles int64
    NameWidth int64
    Query string
    SortOrder string
    SkipHeader bool
    SizeInBytes bool
}

func (self *Drive) List(args ListFilesArgs) (err error) {
    fileList, err := self.service.Files.List().PageSize(args.MaxFiles).Q(args.Query).OrderBy(args.SortOrder).Fields("files(id,name,md5Checksum,mimeType,size,createdTime)").Do()
    if err != nil {
        return fmt.Errorf("Failed listing files: %s", err)
    }

    PrintFileList(PrintFileListArgs{
        Out: args.Out,
        Files: fileList.Files,
        NameWidth: int(args.NameWidth),
        SkipHeader: args.SkipHeader,
        SizeInBytes: args.SizeInBytes,
    })

    return
}

type PrintFileListArgs struct {
    Out io.Writer
    Files []*drive.File
    NameWidth int
    SkipHeader bool
    SizeInBytes bool
}

func PrintFileList(args PrintFileListArgs) {
    w := new(tabwriter.Writer)
    w.Init(args.Out, 0, 0, 3, ' ', 0)

    if !args.SkipHeader {
        fmt.Fprintln(w, "Id\tName\tType\tSize\tCreated")
    }

    for _, f := range args.Files {
        fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
            f.Id,
            truncateString(f.Name, args.NameWidth),
            filetype(f),
            formatSize(f.Size, args.SizeInBytes),
            formatDatetime(f.CreatedTime),
        )
    }

    w.Flush()
}

func filetype(f *drive.File) string {
    if isDir(f) {
        return "dir"
    } else if isBinary(f) {
        return "bin"
    }
    return "doc"
}
