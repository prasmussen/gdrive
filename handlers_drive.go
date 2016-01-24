package main

import (
	"fmt"
	"os"
	"io"
	"io/ioutil"
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
    err := newDrive(args).Download(drive.DownloadArgs{
        Out: os.Stdout,
        Id: args.String("id"),
        Force: args.Bool("force"),
        Path: args.String("path"),
        Recursive: args.Bool("recursive"),
        Stdout: args.Bool("stdout"),
        Progress: progressWriter(args.Bool("noProgress")),
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
        Progress: progressWriter(args.Bool("noProgress")),
    })
    checkErr(err)
}

func uploadHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).Upload(drive.UploadArgs{
        Out: os.Stdout,
        Progress: progressWriter(args.Bool("noProgress")),
        Path: args.String("path"),
        Name: args.String("name"),
        Parents: args.StringSlice("parent"),
        Mime: args.String("mime"),
        Recursive: args.Bool("recursive"),
        Share: args.Bool("share"),
        ChunkSize: args.Int64("chunksize"),
    })
    checkErr(err)
}

func uploadStdinHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).UploadStream(drive.UploadStreamArgs{
        Out: os.Stdout,
        In: os.Stdin,
        Name: args.String("name"),
        Parents: args.StringSlice("parent"),
        Mime: args.String("mime"),
        Share: args.Bool("share"),
        ChunkSize: args.Int64("chunksize"),
    })
    checkErr(err)
}

func updateHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).Update(drive.UpdateArgs{
        Out: os.Stdout,
        Id: args.String("id"),
        Path: args.String("path"),
        Name: args.String("name"),
        Parents: args.StringSlice("parent"),
        Mime: args.String("mime"),
        Stdin: args.Bool("stdin"),
        Share: args.Bool("share"),
        Progress: progressWriter(args.Bool("noProgress")),
        ChunkSize: args.Int64("chunksize"),
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

func deleteHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).Delete(drive.DeleteArgs{
        Out: os.Stdout,
        Id: args.String("id"),
    })
    checkErr(err)
}

func deleteRevisionHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).DeleteRevision(drive.DeleteRevisionArgs{
        Out: os.Stdout,
        FileId: args.String("fileId"),
        RevisionId: args.String("revisionId"),
    })
    checkErr(err)
}

func aboutHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).About(drive.AboutArgs{
        Out: os.Stdout,
        SizeInBytes: args.Bool("sizeInBytes"),
    })
    checkErr(err)
}

func aboutImportHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).AboutImport(drive.AboutImportArgs{
        Out: os.Stdout,
    })
    checkErr(err)
}

func aboutExportHandler(ctx cli.Context) {
    args := ctx.Args()
    err := newDrive(args).AboutExport(drive.AboutExportArgs{
        Out: os.Stdout,
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

func progressWriter(discard bool) io.Writer {
    if discard {
        return ioutil.Discard
    }
    return os.Stderr
}
