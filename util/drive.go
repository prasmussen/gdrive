package util

import (
    "fmt"
)

func PreviewUrl(id string) string {
    //return fmt.Sprintf("https://drive.google.com/uc?id=%s&export=preview", id)
    return fmt.Sprintf("https://drive.google.com/uc?id=%s", id)
}

// Note to self: file.WebContentLink = https://docs.google.com/uc?id=<id>&export=download
func DownloadUrl(id string) string {
    return fmt.Sprintf("https://drive.google.com/uc?id=%s&export=download", id)
}
