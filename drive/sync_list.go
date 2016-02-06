package drive

import (
    "fmt"
    "io"
    "text/tabwriter"
    "google.golang.org/api/googleapi"
    "google.golang.org/api/drive/v3"
)

type ListSyncArgs struct {
    Out io.Writer
    SkipHeader bool
}

func (self *Drive) ListSync(args ListSyncArgs) error {
    query := fmt.Sprintf("appProperties has {key='isSyncRoot' and value='true'}")
    fields := []googleapi.Field{"nextPageToken", "files(id,name,mimeType,createdTime)"}
    files, err := self.listAllFiles(query, fields)
    if err != nil {
        return err
    }
    printSyncDirectories(files, args)
    return nil
}

func printSyncDirectories(files []*drive.File, args ListSyncArgs) {
    w := new(tabwriter.Writer)
    w.Init(args.Out, 0, 0, 3, ' ', 0)

    if !args.SkipHeader {
        fmt.Fprintln(w, "Id\tName\tCreated")
    }

    for _, f := range files {
        fmt.Fprintf(w, "%s\t%s\t%s\n",
            f.Id,
            f.Name,
            formatDatetime(f.CreatedTime),
        )
    }

    w.Flush()
}
