package drive

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"io"
	"sort"
	"text/tabwriter"
)

type ListSyncArgs struct {
	Out        io.Writer
	SkipHeader bool
}

func (self *Drive) ListSync(args ListSyncArgs) error {
	listArgs := listAllFilesArgs{
		query:  "appProperties has {key='syncRoot' and value='true'}",
		fields: []googleapi.Field{"nextPageToken", "files(id,name,mimeType,createdTime)"},
	}
	files, err := self.listAllFiles(listArgs)
	if err != nil {
		return err
	}
	printSyncDirectories(files, args)
	return nil
}

type ListRecursiveSyncArgs struct {
	Out         io.Writer
	RootId      string
	SkipHeader  bool
	PathWidth   int64
	SizeInBytes bool
	SortOrder   string
}

func (self *Drive) ListRecursiveSync(args ListRecursiveSyncArgs) error {
	rootDir, err := self.getSyncRoot(args.RootId)
	if err != nil {
		return err
	}

	files, err := self.prepareRemoteFiles(rootDir, args.SortOrder)
	if err != nil {
		return err
	}

	printSyncDirContent(files, args)
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

func printSyncDirContent(files []*RemoteFile, args ListRecursiveSyncArgs) {
	if args.SortOrder == "" {
		// Sort files by path
		sort.Sort(byRemotePath(files))
	}

	w := new(tabwriter.Writer)
	w.Init(args.Out, 0, 0, 3, ' ', 0)

	if !args.SkipHeader {
		fmt.Fprintln(w, "Id\tPath\tType\tSize\tModified")
	}

	for _, rf := range files {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			rf.file.Id,
			truncateString(rf.relPath, int(args.PathWidth)),
			filetype(rf.file),
			formatSize(rf.file.Size, args.SizeInBytes),
			formatDatetime(rf.file.ModifiedTime),
		)
	}

	w.Flush()
}
