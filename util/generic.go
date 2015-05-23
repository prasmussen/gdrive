package util

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Prompt user to input data
func Prompt(msg string) string {
	fmt.Printf(msg)
	var str string
	fmt.Scanln(&str)
	return str
}

// Returns true if file/directory exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func Mkdir(path string) error {
	dir := filepath.Dir(path)
	if FileExists(dir) {
		return nil
	}
	return os.Mkdir(dir, 0700)
}

// Returns the users home dir
func Homedir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA")
	}
	return os.Getenv("HOME")
}

func FormatBool(b bool) string {
	return strings.Title(strconv.FormatBool(b))
}

func FileSizeFormat(bytes int64, forceBytes bool) string {
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

// Truncates string to given max length, and inserts ellipsis into
// the middle of the string to signify that the string has been truncated
func TruncateString(str string, maxRunes int) string {
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

func ISODateToLocal(iso string) string {
	t, err := time.Parse(time.RFC3339, iso)
	if err != nil {
		return iso
	}
	local := t.Local()
	year, month, day := local.Date()
	hour, min, sec := local.Clock()
	return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, min, sec)
}

func MeasureTransferRate() func(int64) string {
	start := time.Now()

	return func(bytes int64) string {
		seconds := int64(time.Now().Sub(start).Seconds())
		if seconds < 1 {
			return fmt.Sprintf("%s/s", FileSizeFormat(bytes, false))
		}
		bps := bytes / seconds
		return fmt.Sprintf("%s/s", FileSizeFormat(bps, false))
	}
}

// Prints a map in the provided order with one key-value-pair per line
func Print(m map[string]string, keyOrder []string) {
	for _, key := range keyOrder {
		value, ok := m[key]
		if ok && value != "" {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
}

// Prints items in columns with header and correct padding
func PrintColumns(items []map[string]string, keyOrder []string, columnSpacing int, noHeader bool) {

	if !noHeader {
		// Create header
		header := make(map[string]string)
		for _, key := range keyOrder {
			header[key] = key
		}

		// Add header as the first element of items
		items = append([]map[string]string{header}, items...)
	}

	// Get a padding function for each column
	padFns := make(map[string]func(string) string)
	for _, key := range keyOrder {
		padFns[key] = columnPadder(items, key, columnSpacing)
	}

	// Loop, pad and print items
	for _, item := range items {
		var line string

		// Add each column to line with correct padding
		for _, key := range keyOrder {
			value, _ := item[key]
			line += padFns[key](value)
		}

		// Print line
		fmt.Println(line)
	}
}

// Returns a padding function, that pads input to the longest string in items
func columnPadder(items []map[string]string, key string, spacing int) func(string) string {
	// Holds length of longest string
	var max int

	// Find the longest string of type key in the array
	for _, item := range items {
		str := item[key]
		length := utf8.RuneCountInString(str)
		if length > max {
			max = length
		}
	}

	// Return padding function
	return func(str string) string {
		column := str
		for utf8.RuneCountInString(column) < max+spacing {
			column += " "
		}
		return column
	}
}

func inArray(needle string, haystack []string) bool {
	for _, x := range haystack {
		if needle == x {
			return true
		}
	}

	return false
}
