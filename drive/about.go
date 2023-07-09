package drive

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type AboutArgs struct {
	Out         io.Writer
	SizeInBytes bool
}

func (self *Drive) About(args AboutArgs) (err error) {
	about, err := self.service.About.Get().Fields("maxImportSizes", "maxUploadSize", "storageQuota", "user").Do()
	if err != nil {
		return fmt.Errorf("Failed to get about: %s", err)
	}

	user := about.User
	quota := about.StorageQuota

	fmt.Fprintf(args.Out, "User: %s, %s\n", user.DisplayName, user.EmailAddress)
	fmt.Fprintf(args.Out, "Used: %s\n", formatSize(quota.Usage, args.SizeInBytes))
	if quota.Limit > 0 {
		fmt.Fprintf(args.Out, "Free: %s\n", formatSize(quota.Limit-quota.Usage, args.SizeInBytes))
		fmt.Fprintf(args.Out, "Total: %s\n", formatSize(quota.Limit, args.SizeInBytes))
	} else {
		fmt.Println("Total: Unlimited")
	}
	fmt.Fprintf(args.Out, "Max upload size: %s\n", formatSize(about.MaxUploadSize, args.SizeInBytes))
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
