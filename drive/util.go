package drive

import (
    "io"
    "os"
    "fmt"
    "strings"
    "strconv"
    "unicode/utf8"
    "math"
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

func calcRate(bytes int64, start, end time.Time) int64 {
    seconds := float64(end.Sub(start).Seconds())
    if seconds < 1.0 {
        return bytes
    }
    return round(float64(bytes) / seconds)
}

func round(n float64) int64 {
    if n < 0 {
        return int64(math.Ceil(n - 0.5))
    }
    return int64(math.Floor(n + 0.5))
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

func intMax() int64 {
    return 1 << (strconv.IntSize - 1) - 1
}

type Progress struct {
    Writer io.Writer
    Reader io.Reader
    Size int64
    progress int64
    rate int64
    rateProgress int64
    rateUpdated time.Time
    updated time.Time
    done bool
}

func (self *Progress) Read(p []byte) (int, error) {
    // Read
    n, err := self.Reader.Read(p)

    now := time.Now()
    isLast := err != nil

    // Increment progress
    newProgress := self.progress + int64(n)
    self.progress = newProgress

    if self.rateUpdated.IsZero() {
        self.rateUpdated = now
    }

    // Update rate every 3 seconds
    if self.rateUpdated.Add(time.Second * 3).Before(now) {
        self.rate = calcRate(newProgress - self.rateProgress, self.rateUpdated, now)
        self.rateUpdated = now
        self.rateProgress = newProgress
    }

    // Draw progress every second
    if self.updated.Add(time.Second).Before(now) || isLast {
        self.Draw(isLast)
    }

    // Update last draw time
    self.updated = now

    // Mark as done if error occurs
    self.done = isLast

    return n, err
}

func (self *Progress) Draw(isLast bool) {
    if self.done {
        return
    }

    // Clear line
    fmt.Fprintf(self.Writer, "\r%50s", "")

    // Print progress
    fmt.Fprintf(self.Writer, "\r%s/%s", formatSize(self.progress, false), formatSize(self.Size, false))

    // Print rate
    if self.rate > 0 {
        fmt.Fprintf(self.Writer, ", Rate: %s/s", formatSize(self.rate, false))
    }

    if isLast {
        fmt.Fprintf(self.Writer, "\n")
    }
}
