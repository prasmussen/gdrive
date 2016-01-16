package drive

import (
    "fmt"
    "strings"
    "strconv"
    "time"
)


func PrintFileInfo(args PrintFileInfoArgs) {
    f := args.File

    items := []kv{
        kv{"Id", f.Id},
        kv{"Name", f.Name},
        kv{"Description", f.Description},
        kv{"Mime", f.MimeType},
        kv{"Size", formatSize(f.Size, args.SizeInBytes)},
        kv{"Created", formatDatetime(f.CreatedTime)},
        kv{"Modified", formatDatetime(f.ModifiedTime)},
        kv{"Md5sum", f.Md5Checksum},
        kv{"Shared", formatBool(f.Shared)},
        kv{"Parents", formatList(f.Parents)},
    }

    for _, item := range items {
        if item.value() != "" {
            fmt.Printf("%s: %s\n", item.key(), item.value())
        }
    }
}


func formatList(a []string) string {
    return strings.Join(a, ", ")
}

func formatSize(bytes int64, forceBytes bool) string {
    if forceBytes {
        return fmt.Sprintf("%v B", bytes)
    }

    units := []string{"B", "KB", "MB", "GB", "TB", "PB"}

    var i int
    value := float64(bytes)

    for value > 1000 {
        value /= 1000
        i++
    }
    return fmt.Sprintf("%.1f %s", value, units[i])
}

func formatBool(b bool) string {
    return strings.Title(strconv.FormatBool(b))
}

func formatDatetime(iso string) string {
    t, err := time.Parse(time.RFC3339, iso)
    if err != nil {
        return iso
    }
    local := t.Local()
    year, month, day := local.Date()
    hour, min, sec := local.Clock()
    return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
}
