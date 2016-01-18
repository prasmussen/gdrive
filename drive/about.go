package drive

import (
    "io"
    "fmt"
    "text/tabwriter"
)

type AboutArgs struct {
    Out io.Writer
    SizeInBytes bool
    ImportFormats bool
    ExportFormats bool
}

func (self *Drive) About(args AboutArgs) (err error) {
    about, err := self.service.About.Get().Fields("exportFormats", "importFormats", "maxImportSizes", "maxUploadSize", "storageQuota", "user").Do()
    if err != nil {
        return fmt.Errorf("Failed to get about: %s", err)
    }

    if args.ExportFormats {
        printSupportedFormats(args.Out, about.ExportFormats)
        return
    }

    if args.ImportFormats {
        printSupportedFormats(args.Out, about.ImportFormats)
        return
    }

    user := about.User
    quota := about.StorageQuota

    fmt.Fprintf(args.Out, "User: %s, %s\n", user.DisplayName, user.EmailAddress)
    fmt.Fprintf(args.Out, "Used: %s\n", formatSize(quota.UsageInDrive, args.SizeInBytes))
    fmt.Fprintf(args.Out, "Free: %s\n", formatSize(quota.Limit - quota.UsageInDrive, args.SizeInBytes))
    fmt.Fprintf(args.Out, "Total: %s\n", formatSize(quota.Limit, args.SizeInBytes))
    fmt.Fprintf(args.Out, "Max upload size: %s\n", formatSize(about.MaxUploadSize, args.SizeInBytes))
    return
}

func printSupportedFormats(out io.Writer, formats map[string][]string) {
    w := new(tabwriter.Writer)
    w.Init(out, 0, 0, 3, ' ', 0)

    fmt.Fprintln(w, "From\tTo")

    for from, toFormats := range formats {
        fmt.Fprintf(w, "%s\t%s\n", from, formatList(toFormats))
    }

    w.Flush()
}
