package drive

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"io"
	"text/tabwriter"
)

type ListChangesArgs struct {
	Out        io.Writer
	PageToken  string
	MaxChanges int64
	Now        bool
	NameWidth  int64
	SkipHeader bool
}

func (self *Drive) ListChanges(args ListChangesArgs) error {
	if args.Now {
		pageToken, err := self.GetChangesStartPageToken()
		if err != nil {
			return err
		}

		fmt.Fprintf(args.Out, "Page token: %s\n", pageToken)
		return nil
	}

	changeList, err := self.service.Changes.List(args.PageToken).SupportsTeamDrives(true).PageSize(args.MaxChanges).RestrictToMyDrive(true).Fields("newStartPageToken", "nextPageToken", "changes(fileId,removed,time,file(id,name,md5Checksum,mimeType,createdTime,modifiedTime))").Do()
	if err != nil {
		return fmt.Errorf("Failed listing changes: %s", err)
	}

	PrintChanges(PrintChangesArgs{
		Out:        args.Out,
		ChangeList: changeList,
		NameWidth:  int(args.NameWidth),
		SkipHeader: args.SkipHeader,
	})

	return nil
}

func (self *Drive) GetChangesStartPageToken() (string, error) {
	res, err := self.service.Changes.GetStartPageToken().SupportsTeamDrives(true).Do()
	if err != nil {
		return "", fmt.Errorf("Failed getting start page token: %s", err)
	}

	return res.StartPageToken, nil
}

type PrintChangesArgs struct {
	Out        io.Writer
	ChangeList *drive.ChangeList
	NameWidth  int
	SkipHeader bool
}

func PrintChanges(args PrintChangesArgs) {
	w := new(tabwriter.Writer)
	w.Init(args.Out, 0, 0, 3, ' ', 0)

	if !args.SkipHeader {
		fmt.Fprintln(w, "Id\tName\tAction\tTime")
	}

	for _, c := range args.ChangeList.Changes {
		var name string
		var action string

		if c.Removed {
			action = "remove"
		} else {
			name = c.File.Name
			action = "update"
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			c.FileId,
			truncateString(name, args.NameWidth),
			action,
			formatDatetime(c.Time),
		)
	}

	if len(args.ChangeList.Changes) > 0 {
		w.Flush()
		pageToken, hasMore := nextChangesPageToken(args.ChangeList)
		fmt.Fprintf(args.Out, "\nToken: %s, more: %t\n", pageToken, hasMore)
	} else {
		fmt.Fprintln(args.Out, "No changes")
	}
}

func nextChangesPageToken(cl *drive.ChangeList) (string, bool) {
	if cl.NextPageToken != "" {
		return cl.NextPageToken, true
	}

	return cl.NewStartPageToken, false
}
