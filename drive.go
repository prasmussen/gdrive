package main

import (
	"fmt"
	"github.com/prasmussen/gdrive/cli"
	"github.com/prasmussen/gdrive/gdrive"
	"github.com/prasmussen/gdrive/util"
	"github.com/prasmussen/google-api-go-client/googleapi"
	"github.com/voxelbrain/goptions"
	"os"
)

const (
	VersionNumber = "1.9.0"
)

type Options struct {
	Advanced      bool   `goptions:"-a, --advanced, description='Advanced Mode -- lets you specify your own oauth client id and secret on setup'"`
	AppPath       string `goptions:"-c, --config, description='Set application path where config and token is stored. Defaults to ~/.gdrive'"`
	Version       bool   `goptions:"-v, --version, description='Print version'"`
	goptions.Help `goptions:"-h, --help, description='Show this help'"`

	goptions.Verbs

	List struct {
		MaxResults   int    `goptions:"-m, --max, description='Max results'"`
		IncludeDocs  bool   `goptions:"--include-docs, description='Include google docs in listing'"`
		TitleFilter  string `goptions:"-t, --title, mutexgroup='query', description='Title filter'"`
		Query        string `goptions:"-q, --query, mutexgroup='query', description='Query (see https://developers.google.com/drive/search-parameters)'"`
		SharedStatus bool   `goptions:"-s, --shared, description='Show shared status (Note: this will generate 1 http req per file)'"`
		NoHeader     bool   `goptions:"-n, --noheader, description='Do not show the header'"`
		SizeInBytes  bool   `goptions:"--bytes, description='Show size in bytes'"`
	} `goptions:"list"`

	Info struct {
		FileId      string `goptions:"-i, --id, obligatory, description='File Id'"`
		SizeInBytes bool   `goptions:"--bytes, description='Show size in bytes'"`
	} `goptions:"info"`

	Folder struct {
		Title    string `goptions:"-t, --title, obligatory, description='Folder to create'"`
		ParentId string `goptions:"-p, --parent, description='Parent Id of the folder'"`
		Share    bool   `goptions:"--share, description='Share created folder'"`
	} `goptions:"folder"`

	Upload struct {
		File      *os.File `goptions:"-f, --file, mutexgroup='input', obligatory, rdonly, description='File or directory to upload'"`
		Stdin     bool     `goptions:"-s, --stdin, mutexgroup='input', obligatory, description='Use stdin as file content'"`
		Title     string   `goptions:"-t, --title, description='Title to give uploaded file. Defaults to filename'"`
		ParentId  string   `goptions:"-p, --parent, description='Parent Id of the file'"`
		Share     bool     `goptions:"--share, description='Share uploaded file'"`
		MimeType  string   `goptions:"--mimetype, description='The MIME type (default will try to figure it out)'"`
		Convert   bool     `goptions:"--convert, description='File will be converted to Google Docs format'"`
		ChunkSize int64    `goptions:"-C, --chunksize, description='Set chunk size in bytes. Minimum is 262144, default is 4194304. Recommended to be a power of two.'"`
	} `goptions:"upload"`

	Download struct {
		FileId string `goptions:"-i, --id, mutexgroup='download', obligatory, description='File Id'"`
		Format string `goptions:"--format, description='Download file in a specified format (needed for google docs)'"`
		Stdout bool   `goptions:"-s, --stdout, description='Write file content to stdout'"`
		Force  bool   `goptions:"--force, description='Overwrite existing file'"`
		Pop    bool   `goptions:"--pop, mutexgroup='download', description='Download latest file, and remove it from google drive'"`
	} `goptions:"download"`

	Delete struct {
		FileId string `goptions:"-i, --id, obligatory, description='File Id'"`
	} `goptions:"delete"`

	Share struct {
		FileId string `goptions:"-i, --id, obligatory, description='File Id'"`
	} `goptions:"share"`

	Unshare struct {
		FileId string `goptions:"-i, --id, obligatory, description='File Id'"`
	} `goptions:"unshare"`

	Url struct {
		FileId   string `goptions:"-i, --id, obligatory, description='File Id'"`
		Preview  bool   `goptions:"-p, --preview, mutexgroup='urltype', description='Generate preview url (default)'"`
		Download bool   `goptions:"-d, --download, mutexgroup='urltype', description='Generate download url'"`
	} `goptions:"url"`

	Quota struct {
		SizeInBytes bool `goptions:"--bytes, description='Show size in bytes'"`
	} `goptions:"quota"`
}

func main() {
	opts := &Options{}
	goptions.ParseAndFail(opts)

	// Print version number and exit if the version flag is set
	if opts.Version {
		fmt.Printf("gdrive v%s\n", VersionNumber)
		return
	}

	// Get authorized drive client
	drive, err := gdrive.New(opts.AppPath, opts.Advanced, true)
	if err != nil {
		writeError("An error occurred creating Drive client: %v\n", err)
	}

	switch opts.Verbs {
	case "list":
		args := opts.List
		err = cli.List(drive, args.Query, args.TitleFilter, args.MaxResults, args.SharedStatus, args.NoHeader, args.IncludeDocs, args.SizeInBytes)

	case "info":
		err = cli.Info(drive, opts.Info.FileId, opts.Info.SizeInBytes)

	case "folder":
		args := opts.Folder
		err = cli.Folder(drive, args.Title, args.ParentId, args.Share)

	case "upload":
		args := opts.Upload

		// Set custom chunksize if given
		if args.ChunkSize >= (1 << 18) {
			googleapi.SetChunkSize(args.ChunkSize)
		}

		if args.Stdin {
			err = cli.UploadStdin(drive, os.Stdin, args.Title, args.ParentId, args.Share, args.MimeType, args.Convert)
		} else {
			err = cli.Upload(drive, args.File, args.Title, args.ParentId, args.Share, args.MimeType, args.Convert)
		}

	case "download":
		args := opts.Download
		if args.Pop {
			err = cli.DownloadLatest(drive, args.Stdout, args.Format, args.Force)
		} else {
			err = cli.Download(drive, args.FileId, args.Stdout, false, args.Format, args.Force)
		}

	case "delete":
		err = cli.Delete(drive, opts.Delete.FileId)

	case "share":
		err = cli.Share(drive, opts.Share.FileId)

	case "unshare":
		err = cli.Unshare(drive, opts.Unshare.FileId)

	case "url":
		if opts.Url.Download {
			fmt.Println(util.DownloadUrl(opts.Url.FileId))
		} else {
			fmt.Println(util.PreviewUrl(opts.Url.FileId))
		}

	case "quota":
		err = cli.Quota(drive, opts.Quota.SizeInBytes)

	default:
		goptions.PrintHelp()
	}

	if err != nil {
		writeError("%s", err)
	}
}

func writeError(format string, err error) {
	fmt.Fprintf(os.Stderr, format, err)
	fmt.Print("\n")
	os.Exit(1)
}
