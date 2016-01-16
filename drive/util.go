package drive

import (
    "fmt"
    "os"
)

func errorF(err error, format string, a ...interface{}) {
    if err == nil {
        return
    }

    fmt.Fprintf(os.Stderr, format, a...)
    fmt.Println("")
    os.Exit(1)
}

func exitF(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Println("")
	os.Exit(1)
}

func fileExists(path string) bool {
    _, err := os.Stat(path)
    if err == nil {
        return true
    }
    return false
}
