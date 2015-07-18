package util

import (
	"fmt"
	"github.com/prasmussen/google-api-go-client/drive/v2"
	"strings"
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

func InternalDownloadUrlAndExtension(info *drive.File, format string) (downloadUrl string, extension string, err error) {
	// Make a list of available mime types for this file
	availableMimeTypes := make([]string, 0)
	for mime, _ := range info.ExportLinks {
		availableMimeTypes = append(availableMimeTypes, mime)
	}

	mimeExtensions := map[string]string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":               "docx",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":                     "xlsx",
		"application/application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
		"application/vnd.oasis.opendocument.text":                                               "odf",
		"application/x-vnd.oasis.opendocument.spreadsheet":                                      "ods",
		"application/pdf":                                                                       "pdf",
		"application/rtf":                                                                       "rtf",
		"text/csv":                                                                              "csv",
		"text/html":                                                                             "html",
		"text/plain":                                                                            "txt",
		"application/vnd.google-apps.script+json":                                               "json",
	}

	// Make a list of available formats for this file
	availableFormats := make([]string, 0)
	for _, mime := range availableMimeTypes {
		if ext, ok := mimeExtensions[mime]; ok {
			availableFormats = append(availableFormats, ext)
		}
	}

	// Return DownloadUrl if no format is specified
	if format == "" {
		if info.DownloadUrl == "" {
			if len(availableFormats) > 0 {
				return "", "", fmt.Errorf("A format needs to be specified to download this file (--format). Available formats: %s", strings.Join(availableFormats, ", "))
			} else {
				return "", "", fmt.Errorf("Download is not supported for this filetype")
			}
		}
		return info.DownloadUrl, "", nil
	}

	// Ensure that the specified format is available
	if !inArray(format, availableFormats) {
		if len(availableFormats) > 0 {
			return "", "", fmt.Errorf("Invalid format. Available formats: %s", strings.Join(availableFormats, ", "))
		} else {
			return "", "", fmt.Errorf("No export formats are available for this file")
		}
	}

	// Grab download url
	for mime, f := range mimeExtensions {
		if f == format {
			downloadUrl = info.ExportLinks[mime]
			break
		}
	}

	extension = "." + format
	return
}
