package drive

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
	"text/tabwriter"
)

type ListRevisionsArgs struct {
	Out         io.Writer
	Id          string
	NameWidth   int64
	SkipHeader  bool
	SizeInBytes bool
}

func (self *Drive) ListRevisions(args ListRevisionsArgs) (err error) {
	revList, err := self.service.Revisions.List(args.Id).Fields("revisions(id,keepForever,size,modifiedTime,originalFilename)").Do()
	if err != nil {
		return fmt.Errorf("Failed listing revisions: %s", err)
	}

	PrintRevisionList(PrintRevisionListArgs{
		Out:         args.Out,
		Revisions:   revList.Revisions,
		NameWidth:   int(args.NameWidth),
		SkipHeader:  args.SkipHeader,
		SizeInBytes: args.SizeInBytes,
	})

	return
}

type PrintRevisionListArgs struct {
	Out         io.Writer
	Revisions   []*drive.Revision
	NameWidth   int
	SkipHeader  bool
	SizeInBytes bool
}

func PrintRevisionList(args PrintRevisionListArgs) {
	w := new(tabwriter.Writer)
	w.Init(args.Out, 0, 0, 3, ' ', 0)

	if !args.SkipHeader {
		fmt.Fprintln(w, "Id\tName\tSize\tModified\tKeepForever")
	}

	for _, rev := range args.Revisions {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			rev.Id,
			truncateString(rev.OriginalFilename, args.NameWidth),
			formatSize(rev.Size, args.SizeInBytes),
			formatDatetime(rev.ModifiedTime),
			formatBool(rev.KeepForever),
		)
	}

	w.Flush()
}
