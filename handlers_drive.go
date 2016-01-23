package main

import (
	"fmt"
	"os"
    "./cli"
	"./auth"
	"./drive"
)

const ClientId     = "367116221053-7n0vf5akeru7on6o2fjinrecpdoe99eg.apps.googleusercontent.com"
const ClientSecret = "1qsNodXNaWq1mQuBjUjmvhoO"
const TokenFilename = "token_v2.json"


func listHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).List(drive.ListFilesArgs{
        Out: os.Stdout,
        MaxFiles: args.Int64("maxFiles"),
        NameWidth: args.Int64("nameWidth"),
        Query: args.String("query"),
        SkipHeader: args.Bool("skipHeader"),
        SizeInBytes: args.Bool("sizeInBytes"),
    })
    checkErr(err)
}

func downloadHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).Download(drive.DownloadFileArgs{
        Out: os.Stdout,
        Id: args.String("id"),
        Force: args.Bool("force"),
        Stdout: args.Bool("stdout"),
        NoProgress: args.Bool("noProgress"),
    })
    checkErr(err)
}

func downloadRevisionHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).DownloadRevision(drive.DownloadRevisionArgs{
        Out: os.Stdout,
        FileId: args.String("fileId"),
        RevisionId: args.String("revisionId"),
        Force: args.Bool("force"),
        Stdout: args.Bool("stdout"),
        NoProgress: args.Bool("noProgress"),
    })
    checkErr(err)
}

func uploadHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).Upload(drive.UploadFileArgs{
        Out: os.Stdout,
        Path: args.String("path"),
        Name: args.String("name"),
        Parents: args.StringSlice("parent"),
        Mime: args.String("mime"),
        Recursive: args.Bool("recursive"),
        Stdin: args.Bool("stdin"),
        Share: args.Bool("share"),
        NoProgress: args.Bool("noProgress"),
    })
    checkErr(err)
}

func infoHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).Info(drive.FileInfoArgs{
        Out: os.Stdout,
        Id: args.String("id"),
        SizeInBytes: args.Bool("sizeInBytes"),
    })
    checkErr(err)
}

func exportHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).Export(drive.ExportArgs{
        Out: os.Stdout,
        Id: args.String("id"),
        Mime: args.String("mime"),
        PrintMimes: args.Bool("printMimes"),
        Force: args.Bool("force"),
    })
    checkErr(err)
}

func listRevisionsHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).ListRevisions(drive.ListRevisionsArgs{
        Out: os.Stdout,
        Id: args.String("id"),
        NameWidth: args.Int64("nameWidth"),
        SizeInBytes: args.Bool("sizeInBytes"),
        SkipHeader: args.Bool("skipHeader"),
    })
    checkErr(err)
}

func mkdirHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).Mkdir(drive.MkdirArgs{
        Out: os.Stdout,
        Name: args.String("name"),
        Parents: args.StringSlice("parent"),
        Share: args.Bool("share"),
    })
    checkErr(err)
}

func shareHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).Share(drive.ShareArgs{
        Out: os.Stdout,
        FileId: args.String("id"),
        Role: args.String("role"),
        Type: args.String("type"),
        Email: args.String("email"),
        Discoverable: args.Bool("discoverable"),
        Revoke: args.Bool("revoke"),
    })
    checkErr(err)
}

func urlHandler(ctx cli.Context) {
    args := ctx.Args()
    newDrive(args).Url(drive.UrlArgs{
        Out: os.Stdout,
        FileId: args.String("id"),
        DownloadUrl: args.Bool("download"),
    })
}

func deleteHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).Delete(drive.DeleteArgs{
        Out: os.Stdout,
        Id: args.String("id"),
    })
    checkErr(err)
}

func aboutHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).About(drive.AboutArgs{
        Out: os.Stdout,
        SizeInBytes: args.Bool("sizeInBytes"),
        ImportFormats: args.Bool("importFormats"),
        ExportFormats: args.Bool("exportFormats"),
    })
    checkErr(err)
}

func newDrive(args cli.Arguments) *drive.Drive {
    configDir := args.String("configDir")
    tokenPath := ConfigFilePath(configDir, TokenFilename)
    oauth, err := auth.NewOauthClient(ClientId, ClientSecret, tokenPath, authCodePrompt)
    if err != nil {
        ExitF("Failed getting oauth client: %s", err.Error())
    }

    client, err := drive.New(oauth)
    if err != nil {
        ExitF("Failed getting drive: %s", err.Error())
    }

    return client
}

func authCodePrompt(url string) func() string {
    return func() string {
        fmt.Println("Authentication needed")
        fmt.Println("Go to the following url in your browser:")
        fmt.Printf("%s\n\n", url)
        fmt.Print("Enter verification code: ")

        var code string
        if _, err := fmt.Scan(&code); err != nil {
            fmt.Printf("Failed reading code: %s", err.Error())
        }
        return code
    }
}
