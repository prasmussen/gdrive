package drive

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

type AboutArgs struct {
	Out         io.Writer
	SizeInBytes bool
	OutJSON     bool
}

func (self *Drive) About(args AboutArgs) (err error) {
	about, err := self.service.About.Get().Fields("maxImportSizes", "maxUploadSize", "storageQuota", "user").Do()
	if err != nil {
		return fmt.Errorf("Failed to get about: %s", err)
	}

	user := about.User
	quota := about.StorageQuota

	if args.OutJSON {
		jsonabout := map[string]string{
			"displayname":     user.DisplayName,
			"email":           user.EmailAddress,
			"used":            formatSize(quota.Usage, args.SizeInBytes),
			"free":            formatSize(quota.Limit-quota.Usage, args.SizeInBytes),
			"total":           formatSize(quota.Limit, args.SizeInBytes),
			"max_upload_size": formatSize(about.MaxUploadSize, args.SizeInBytes),
		}
		enc := json.NewEncoder(os.Stdout)
		if err := enc.Encode(jsonabout); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Fprintf(args.Out, "User: %s, %s\n", user.DisplayName, user.EmailAddress)
		fmt.Fprintf(args.Out, "Used: %s\n", formatSize(quota.Usage, args.SizeInBytes))
		fmt.Fprintf(args.Out, "Free: %s\n", formatSize(quota.Limit-quota.Usage, args.SizeInBytes))
		fmt.Fprintf(args.Out, "Total: %s\n", formatSize(quota.Limit, args.SizeInBytes))
		fmt.Fprintf(args.Out, "Max upload size: %s\n", formatSize(about.MaxUploadSize, args.SizeInBytes))
	}
	return
}

type AboutImportArgs struct {
	Out io.Writer
}

func (self *Drive) AboutImport(args AboutImportArgs) (err error) {
	about, err := self.service.About.Get().Fields("importFormats").Do()
	if err != nil {
		return fmt.Errorf("Failed to get about: %s", err)
	}
	printAboutFormats(args.Out, about.ImportFormats)
	return
}

type AboutExportArgs struct {
	Out io.Writer
}

func (self *Drive) AboutExport(args AboutExportArgs) (err error) {
	about, err := self.service.About.Get().Fields("exportFormats").Do()
	if err != nil {
		return fmt.Errorf("Failed to get about: %s", err)
	}
	printAboutFormats(args.Out, about.ExportFormats)
	return
}

func printAboutFormats(out io.Writer, formats map[string][]string) {
	w := new(tabwriter.Writer)
	w.Init(out, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "From\tTo")

	for from, toFormats := range formats {
		fmt.Fprintf(w, "%s\t%s\n", from, formatList(toFormats))
	}

	w.Flush()
}
