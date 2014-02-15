package util

import (
    "fmt"
    "strings"
    "code.google.com/p/google-api-go-client/drive/v2"
)

func PreviewUrl(id string) string {
    //return fmt.Sprintf("https://drive.google.com/uc?id=%s&export=preview", id)
    return fmt.Sprintf("https://drive.google.com/uc?id=%s", id)
}

// Note to self: file.WebContentLink = https://docs.google.com/uc?id=<id>&export=download
func DownloadUrl(id string) string {
    return fmt.Sprintf("https://drive.google.com/uc?id=%s&export=download", id)
}

func ParentList(parents []*drive.ParentReference) string {
    ids := make([]string, 0)
    for _, parent := range parents {
        ids = append(ids, parent.Id)
    }

    return strings.Join(ids, ", ")
}
