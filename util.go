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

func ConfigFilePath(basePath, name string) string {
    return filepath.Join(basePath, name)
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

func checkErr(err error) {
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
