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

func ExportFormat(info *drive.File, format string) (downloadUrl string, extension string, err error) {
	// See https://developers.google.com/drive/web/manage-downloads#downloading_google_documents
	switch format {
	case "docx":
		extension = ".docx"
		downloadUrl = info.ExportLinks["application/vnd.openxmlformats-officedocument.wordprocessingml.document"]
	case "xlsx":
		extension = ".xlsx"
		downloadUrl = info.ExportLinks["application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"]
	case "pptx":
		extension = ".pptx"
		downloadUrl = info.ExportLinks["application/application/vnd.openxmlformats-officedocument.presentationml.presentation"]
	case "odf":
		extension = ".odf"
		downloadUrl = info.ExportLinks["application/vnd.oasis.opendocument.text"]
	case "ods":
		extension = ".ods"
		downloadUrl = info.ExportLinks["application/x-vnd.oasis.opendocument.spreadsheet"]
	case "pdf":
		extension = ".pdf"
		downloadUrl = info.ExportLinks["application/pdf"]
	case "rtf":
		extension = ".rtf"
		downloadUrl = info.ExportLinks["application/rtf"]
	case "csv":
		extension = ".csv"
		downloadUrl = info.ExportLinks["text/csv"]
	case "html":
		extension = ".html"
		downloadUrl = info.ExportLinks["text/html"]
	case "txt":
		extension = ".txt"
		downloadUrl = info.ExportLinks["text/plain"]
	case "json":
		extension = ".json"
		downloadUrl = info.ExportLinks["application/vnd.google-apps.script+json"]
	default:
		err = fmt.Errorf("Unknown export format: %s. Known formats: docx, xlsx, pptx, odf, ods, pdf, rtf, csv, txt, html, json", format)
	}
	return
}
