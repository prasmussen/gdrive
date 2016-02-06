package main

import (
	"fmt"
    "os"
    "io"
    "crypto/md5"
	"./drive"
)

type Md5Comparer struct {}

func (self Md5Comparer) Changed(local *drive.LocalFile, remote *drive.RemoteFile) bool {
    return remote.Md5() != md5sum(local.AbsPath())
}

func md5sum(path string) string {
    h := md5.New()
    f, err := os.Open(path)
    if err != nil {
        return ""
    }
    defer f.Close()

    io.Copy(h, f)
    return fmt.Sprintf("%x", h.Sum(nil))
}
