package main

import (
    "runtime"
    "path/filepath"
    "fmt"
    "os"
)

func GetDefaultConfigDir() string {
    return filepath.Join(Homedir(), ".gdrive")
}

func GetDefaultTokenFilePath() string {
    return filepath.Join(GetDefaultConfigDir(), "token.json")
}

func Homedir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("APPDATA")
	}
	return os.Getenv("HOME")
}

func ExitF(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Println("")
	os.Exit(1)
}
