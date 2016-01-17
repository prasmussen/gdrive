package drive

import (
    "fmt"
    "os"
    "text/tabwriter"
    "google.golang.org/api/drive/v3"
)

type ListFilesArgs struct {
    MaxFiles int64
    NameWidth int64
    Query string
    SkipHeader bool
    SizeInBytes bool
}

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

type PrintFileListArgs struct {
    Files []*drive.File
    NameWidth int
    SkipHeader bool
    SizeInBytes bool
}

func PrintFileList(args PrintFileListArgs) {
    w := new(tabwriter.Writer)
    w.Init(os.Stdout, 0, 0, 3, ' ', 0)

    if !args.SkipHeader {
        fmt.Fprintln(w, "Id\tName\tSize\tCreated")
    }

    for _, f := range args.Files {
        fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
            f.Id,
            truncateString(f.Name, args.NameWidth),
            formatSize(f.Size, args.SizeInBytes),
            formatDatetime(f.CreatedTime),
        )
    }

    w.Flush()
}
