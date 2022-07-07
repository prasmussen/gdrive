package drive

import (
	"context"
	"fmt"
	"io"
	"text/tabwriter"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

type ListFilesArgs struct {
	Out         io.Writer
	MaxFiles    int64
	NameWidth   int64
	Query       string
	SortOrder   string
	SkipHeader  bool
	SizeInBytes bool
	AbsPath     bool
}

func (self *Drive) List(args ListFilesArgs) (err error) {
	listArgs := listAllFilesArgs{
		query:     args.Query,
		fields:    []googleapi.Field{"nextPageToken", "files(id,name,md5Checksum,mimeType,size,createdTime,parents)"},
		sortOrder: args.SortOrder,
		maxFiles:  args.MaxFiles,
	}
	files, err := self.listAllFiles(listArgs)
	if err != nil {
		return fmt.Errorf("Failed to list files: %s", err)
	}

	pathfinder := self.newPathfinder()

	if args.AbsPath {
		// Replace name with absolute path
		for _, f := range files {
			f.Name, err = pathfinder.absPath(f)
			if err != nil {
				return err
			}
		}
	}

	PrintFileList(PrintFileListArgs{
		Out:         args.Out,
		Files:       files,
		NameWidth:   int(args.NameWidth),
		SkipHeader:  args.SkipHeader,
		SizeInBytes: args.SizeInBytes,
	})

	return
}

type listAllFilesArgs struct {
	query     string
	fields    []googleapi.Field
	sortOrder string
	maxFiles  int64
}

func (self *Drive) listAllFiles(args listAllFilesArgs) ([]*drive.File, error) {
	return self.listAllFilesWithContext(context.TODO(), &args)
}

type PrintFileListArgs struct {
	Out         io.Writer
	Files       []*drive.File
	NameWidth   int
	SkipHeader  bool
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
