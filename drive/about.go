package drive

import (
    "fmt"
    "os"
    "text/tabwriter"
)

type AboutArgs struct {
    SizeInBytes bool
    ImportFormats bool
    ExportFormats bool
}

func (self *Drive) About(args AboutArgs) {
    about, err := self.service.About.Get().Fields("exportFormats", "importFormats", "maxImportSizes", "maxUploadSize", "storageQuota", "user").Do()
    errorF(err, "Failed to get about %s", err)

    if args.ExportFormats {
        printSupportedFormats(about.ExportFormats)
        return
    }

    if args.ImportFormats {
        printSupportedFormats(about.ImportFormats)
        return
    }

    user := about.User
    quota := about.StorageQuota

    fmt.Printf("User: %s, %s\n", user.DisplayName, user.EmailAddress)
    fmt.Printf("Used: %s\n", formatSize(quota.UsageInDrive, args.SizeInBytes))
    fmt.Printf("Free: %s\n", formatSize(quota.Limit - quota.UsageInDrive, args.SizeInBytes))
    fmt.Printf("Total: %s\n", formatSize(quota.Limit, args.SizeInBytes))
    fmt.Printf("Max upload size: %s\n", formatSize(about.MaxUploadSize, args.SizeInBytes))
}

func printSupportedFormats(formats map[string][]string) {
    w := new(tabwriter.Writer)
    w.Init(os.Stdout, 0, 0, 3, ' ', 0)

    fmt.Fprintln(w, "From\tTo")

    for from, toFormats := range formats {
        fmt.Fprintf(w, "%s\t%s\n", from, formatList(toFormats))
    }

    w.Flush()
}
