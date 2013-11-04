package main

import (
    "fmt"
    "os"
    "github.com/voxelbrain/goptions"
    "./gdrive"
    "./util"
    "./cli"
)

const (
    VersionNumber = "1.0.1"
)

type Options struct {
    Advanced bool `goptions:"-a, --advanced, description='Advanced Mode -- lets you specify your own oauth client id and secret on setup'"`
    AppPath string `goptions:"-c, --config, description='Set application path where config and token is stored. Defaults to ~/.gdrive'"`
    Version bool `goptions:"-v, --version, description='Print version'"`
    goptions.Help `goptions:"-h, --help, description='Show this help'"`

    goptions.Verbs

    List struct {
        MaxResults int `goptions:"-m, --max, description='Max results'"`
        TitleFilter string `goptions:"-t, --title, mutexgroup='query', description='Title filter'"`
        Query string `goptions:"-q, --query, mutexgroup='query', description='Query (see https://developers.google.com/drive/search-parameters)'"`
        SharedStatus bool `goptions:"-s, --shared, description='Show shared status (Note: this will generate 1 http req per file)'"`
    } `goptions:"list"`

    Info struct {
        FileId string `goptions:"-i, --id, obligatory, description='File Id'"`
    } `goptions:"info"`

    Upload struct {
        File *os.File `goptions:"-f, --file, mutexgroup='input', obligatory, rdonly, description='File to upload'"`
        Stdin bool `goptions:"-s, --stdin, mutexgroup='input', obligatory, description='Use stdin as file content'"`
        Title string `goptions:"-t, --title, description='Title to give uploaded file. Defaults to filename'"`
        Share bool `goptions:"--share, description='Share uploaded file'"`
    } `goptions:"upload"`

    Download struct {
        FileId string `goptions:"-i, --id, mutexgroup='download', obligatory, description='File Id'"`
        Stdout bool `goptions:"-s, --stdout, description='Write file content to stdout'"`
        Pop bool `goptions:"--pop, mutexgroup='download', description='Download latest file, and remove it from google drive'"`
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
        FileId string `goptions:"-i, --id, obligatory, description='File Id'"`
        Preview bool `goptions:"-p, --preview, mutexgroup='urltype', description='Generate preview url (default)'"`
        Download bool `goptions:"-d, --download, mutexgroup='urltype', description='Generate download url'"`
    } `goptions:"url"`
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
    drive, err := gdrive.New(opts.AppPath, opts.Advanced)
    if err != nil {
        fmt.Printf("An error occurred creating Drive client: %v\n", err)
        os.Exit(1)
    }

    switch opts.Verbs {
        case "list":
            args := opts.List
            cli.List(drive, args.Query, args.TitleFilter, args.MaxResults, args.SharedStatus)

        case "info":
            cli.Info(drive, opts.Info.FileId)

        case "upload":
            args := opts.Upload
            if args.Stdin {
                cli.Upload(drive, os.Stdin, args.Title, args.Share)
            } else {
                cli.Upload(drive, args.File, args.Title, args.Share)
            }

        case "download":
            args := opts.Download
            if args.Pop {
                cli.DownloadLatest(drive, args.Stdout)
            } else {
                cli.Download(drive, args.FileId, args.Stdout, false)
            }

        case "delete":
            cli.Delete(drive, opts.Delete.FileId)

        case "share":
            cli.Share(drive, opts.Share.FileId)

        case "unshare":
            cli.Unshare(drive, opts.Unshare.FileId)

        case "url":
            if opts.Url.Download {
                fmt.Println(util.DownloadUrl(opts.Url.FileId))
            } else {
                fmt.Println(util.PreviewUrl(opts.Url.FileId))
            }

        default:
            goptions.PrintHelp()
    }
}


