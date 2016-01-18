package drive

import (
    "os"
    "fmt"
    "strings"
    "strconv"
    "unicode/utf8"
    "time"
)

type kv struct {
    key string
    value string
}

func formatList(a []string) string {
    return strings.Join(a, ", ")
}

func formatSize(bytes int64, forceBytes bool) string {
    if bytes == 0 {
        return ""
    }

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

// Truncates string to given max length, and inserts ellipsis into
// the middle of the string to signify that the string has been truncated
func truncateString(str string, maxRunes int) string {
    indicator := "..."

    // Number of runes in string
    runeCount := utf8.RuneCountInString(str)

    // Return input string if length of input string is less than max length
    // Input string is also returned if max length is less than 9 which is the minmal supported length
    if runeCount <= maxRunes || maxRunes < 9 {
        return str
    }

    // Number of remaining runes to be removed
    remaining := (runeCount - maxRunes) + utf8.RuneCountInString(indicator)

    var truncated string
    var skip bool

    for leftOffset, char := range str {
        rightOffset := runeCount - (leftOffset + remaining)

        // Start skipping chars when the left and right offsets are equal
        // Or in the case where we wont be able to do an even split: when the left offset is larger than the right offset
        if leftOffset == rightOffset || (leftOffset > rightOffset && !skip) {
            skip = true
            truncated += indicator
        }

        if skip && remaining > 0 {
            // Skip char and decrement the remaining skip counter
            remaining--
            continue
        }

        // Add char to result string
        truncated += string(char)
    }

    // Return truncated string
    return truncated
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    if err == nil {
        return true
    }
    return false
}
