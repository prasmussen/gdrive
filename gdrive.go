package main

import (
	"fmt"
	"os"
    "./cli"
)

const Name = "gdrive"
const Version = "2.0.0"

const DefaultMaxFiles = 30
const DefaultMaxChanges = 100
const DefaultNameWidth = 40
const DefaultUploadChunkSize = 8 * 1024 * 1024
const DefaultQuery = "trashed = false and 'me' in owners"
const DefaultShareRole = "reader"
const DefaultShareType = "anyone"
var DefaultConfigDir = GetDefaultConfigDir()


func main() {
    globalFlags := []cli.Flag{
        cli.StringFlag{
            Name: "configDir",
            Patterns: []string{"-c", "--config"},
            Description: fmt.Sprintf("Application path, default: %s", DefaultConfigDir),
            DefaultValue: DefaultConfigDir,
        },
    }

    handlers := []*cli.Handler{
        &cli.Handler{
            Pattern: "[global] list [options]",
            Description: "List files",
            Callback: listHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.IntFlag{
                        Name: "maxFiles",
                        Patterns: []string{"-m", "--max"},
                        Description: fmt.Sprintf("Max files to list, default: %d", DefaultMaxFiles),
                        DefaultValue: DefaultMaxFiles,
                    },
                    cli.StringFlag{
                        Name: "query",
                        Patterns: []string{"-q", "--query"},
                        Description: fmt.Sprintf(`Default query: "%s". See https://developers.google.com/drive/search-parameters`, DefaultQuery),
                        DefaultValue: DefaultQuery,
                    },
                    cli.StringFlag{
                        Name: "sortOrder",
                        Patterns: []string{"--order"},
                        Description: "Sort order. See https://godoc.org/google.golang.org/api/drive/v3#FilesListCall.OrderBy",
                    },
                    cli.IntFlag{
                        Name: "nameWidth",
                        Patterns: []string{"--name-width"},
                        Description: fmt.Sprintf("Width of name column, default: %d, minimum: 9, use 0 for full width", DefaultNameWidth),
                        DefaultValue: DefaultNameWidth,
                    },
                    cli.BoolFlag{
                        Name: "skipHeader",
                        Patterns: []string{"--no-header"},
                        Description: "Dont print the header",
                        OmitValue: true,
                    },
                    cli.BoolFlag{
                        Name: "sizeInBytes",
                        Patterns: []string{"--bytes"},
                        Description: "Size in bytes",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] list changes [options]",
            Description: "List file changes",
            Callback: listChangesHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.IntFlag{
                        Name: "maxChanges",
                        Patterns: []string{"-m", "--max"},
                        Description: fmt.Sprintf("Max changes to list, default: %d", DefaultMaxChanges),
                        DefaultValue: DefaultMaxChanges,
                    },
                    cli.StringFlag{
                        Name: "pageToken",
                        Patterns: []string{"--since"},
                        Description: fmt.Sprintf("Page token to start listing changes from"),
                        DefaultValue: "1",
                    },
                    cli.BoolFlag{
                        Name: "now",
                        Patterns: []string{"--now"},
                        Description: fmt.Sprintf("Get latest page token"),
                        OmitValue: true,
                    },
                    cli.IntFlag{
                        Name: "nameWidth",
                        Patterns: []string{"--name-width"},
                        Description: fmt.Sprintf("Width of name column, default: %d, minimum: 9, use 0 for full width", DefaultNameWidth),
                        DefaultValue: DefaultNameWidth,
                    },
                    cli.BoolFlag{
                        Name: "skipHeader",
                        Patterns: []string{"--no-header"},
                        Description: "Dont print the header",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] list revisions [options] <id>",
            Description: "List file revisions",
            Callback: listRevisionsHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.IntFlag{
                        Name: "nameWidth",
                        Patterns: []string{"--name-width"},
                        Description: fmt.Sprintf("Width of name column, default: %d, minimum: 9, use 0 for full width", DefaultNameWidth),
                        DefaultValue: DefaultNameWidth,
                    },
                    cli.BoolFlag{
                        Name: "skipHeader",
                        Patterns: []string{"--no-header"},
                        Description: "Dont print the header",
                        OmitValue: true,
                    },
                    cli.BoolFlag{
                        Name: "sizeInBytes",
                        Patterns: []string{"--bytes"},
                        Description: "Size in bytes",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] download [options] <id>",
            Description: "Download file or directory",
            Callback: downloadHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "force",
                        Patterns: []string{"-f", "--force"},
                        Description: "Overwrite existing file",
                        OmitValue: true,
                    },
                    cli.BoolFlag{
                        Name: "recursive",
                        Patterns: []string{"-r", "--recursive"},
                        Description: "Download directory recursively, documents will be skipped",
                        OmitValue: true,
                    },
                    cli.StringFlag{
                        Name: "path",
                        Patterns: []string{"--path"},
                        Description: "Download path",
                    },
                    cli.BoolFlag{
                        Name: "noProgress",
                        Patterns: []string{"--no-progress"},
                        Description: "Hide progress",
                        OmitValue: true,
                    },
                    cli.BoolFlag{
                        Name: "stdout",
                        Patterns: []string{"--stdout"},
                        Description: "Write file content to stdout",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] download revision [options] <fileId> <revisionId>",
            Description: "Download revision",
            Callback: downloadRevisionHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "force",
                        Patterns: []string{"-f", "--force"},
                        Description: "Overwrite existing file",
                        OmitValue: true,
                    },
                    cli.BoolFlag{
                        Name: "noProgress",
                        Patterns: []string{"--no-progress"},
                        Description: "Hide progress",
                        OmitValue: true,
                    },
                    cli.BoolFlag{
                        Name: "stdout",
                        Patterns: []string{"--stdout"},
                        Description: "Write file content to stdout",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] upload [options] <path>",
            Description: "Upload file or directory",
            Callback: uploadHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "recursive",
                        Patterns: []string{"-r", "--recursive"},
                        Description: "Upload directory recursively",
                        OmitValue: true,
                    },
                    cli.StringSliceFlag{
                        Name: "parent",
                        Patterns: []string{"-p", "--parent"},
                        Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
                    },
                    cli.StringFlag{
                        Name: "name",
                        Patterns: []string{"--name"},
                        Description: "Filename",
                    },
                    cli.BoolFlag{
                        Name: "noProgress",
                        Patterns: []string{"--no-progress"},
                        Description: "Hide progress",
                        OmitValue: true,
                    },
                    cli.StringFlag{
                        Name: "mime",
                        Patterns: []string{"--mime"},
                        Description: "Force mime type",
                    },
                    cli.BoolFlag{
                        Name: "share",
                        Patterns: []string{"--share"},
                        Description: "Share file",
                        OmitValue: true,
                    },
                    cli.IntFlag{
                        Name: "chunksize",
                        Patterns: []string{"--chunksize"},
                        Description: fmt.Sprintf("Set chunk size in bytes, default: %d", DefaultUploadChunkSize),
                        DefaultValue: DefaultUploadChunkSize,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] upload stdin [options] <name>",
            Description: "Upload file from stdin",
            Callback: uploadStdinHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.StringSliceFlag{
                        Name: "parent",
                        Patterns: []string{"-p", "--parent"},
                        Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
                    },
                    cli.IntFlag{
                        Name: "chunksize",
                        Patterns: []string{"--chunksize"},
                        Description: fmt.Sprintf("Set chunk size in bytes, default: %d", DefaultUploadChunkSize),
                        DefaultValue: DefaultUploadChunkSize,
                    },
                    cli.StringFlag{
                        Name: "mime",
                        Patterns: []string{"--mime"},
                        Description: "Force mime type",
                    },
                    cli.BoolFlag{
                        Name: "share",
                        Patterns: []string{"--share"},
                        Description: "Share file",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] upload sync [options] <path> <id>",
            Description: "Sync local directory to drive",
            Callback: uploadSyncHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "noProgress",
                        Patterns: []string{"--no-progress"},
                        Description: "Hide progress",
                        OmitValue: true,
                    },
                    cli.BoolFlag{
                        Name: "deleteExtraneous",
                        Patterns: []string{"--delete-extraneous"},
                        Description: "Delete extraneous files from drive",
                        OmitValue: true,
                    },
                    cli.IntFlag{
                        Name: "chunksize",
                        Patterns: []string{"--chunksize"},
                        Description: fmt.Sprintf("Set chunk size in bytes, default: %d", DefaultUploadChunkSize),
                        DefaultValue: DefaultUploadChunkSize,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] update [options] <id> <path>",
            Description: "Update file, this creates a new revision of the file",
            Callback: updateHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.StringSliceFlag{
                        Name: "parent",
                        Patterns: []string{"-p", "--parent"},
                        Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
                    },
                    cli.StringFlag{
                        Name: "name",
                        Patterns: []string{"--name"},
                        Description: "Filename",
                    },
                    cli.BoolFlag{
                        Name: "noProgress",
                        Patterns: []string{"--no-progress"},
                        Description: "Hide progress",
                        OmitValue: true,
                    },
                    cli.BoolFlag{
                        Name: "stdin",
                        Patterns: []string{"--stdin"},
                        Description: "Use stdin as file content",
                        OmitValue: true,
                    },
                    cli.StringFlag{
                        Name: "mime",
                        Patterns: []string{"--mime"},
                        Description: "Force mime type",
                    },
                    cli.BoolFlag{
                        Name: "share",
                        Patterns: []string{"--share"},
                        Description: "Share file",
                        OmitValue: true,
                    },
                    cli.IntFlag{
                        Name: "chunksize",
                        Patterns: []string{"--chunksize"},
                        Description: fmt.Sprintf("Set chunk size in bytes, default: %d", DefaultUploadChunkSize),
                        DefaultValue: DefaultUploadChunkSize,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] info [options] <id>",
            Description: "Show file info",
            Callback: infoHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "sizeInBytes",
                        Patterns: []string{"--bytes"},
                        Description: "Show size in bytes",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] import [options] <path>",
            Description: "Upload and convert file to a google document, see 'about import' for available conversions",
            Callback: importHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.StringSliceFlag{
                        Name: "parent",
                        Patterns: []string{"-p", "--parent"},
                        Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
                    },
                    cli.BoolFlag{
                        Name: "noProgress",
                        Patterns: []string{"--no-progress"},
                        Description: "Hide progress",
                        OmitValue: true,
                    },
                    cli.BoolFlag{
                        Name: "share",
                        Patterns: []string{"--share"},
                        Description: "Share file",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] export [options] <id>",
            Description: "Export a google document",
            Callback: exportHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "force",
                        Patterns: []string{"-f", "--force"},
                        Description: "Overwrite existing file",
                        OmitValue: true,
                    },
                    cli.StringFlag{
                        Name: "mime",
                        Patterns: []string{"--mime"},
                        Description: "Mime type of exported file",
                    },
                    cli.BoolFlag{
                        Name: "printMimes",
                        Patterns: []string{"--print-mimes"},
                        Description: "Print available mime types for given file",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] mkdir [options] <name>",
            Description: "Create directory",
            Callback: mkdirHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.StringSliceFlag{
                        Name: "parent",
                        Patterns: []string{"-p", "--parent"},
                        Description: "Parent id of created directory, can be specified multiple times to give many parents",
                    },
                    cli.BoolFlag{
                        Name: "share",
                        Patterns: []string{"--share"},
                        Description: "Share created directory",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] share [options] <id>",
            Description: "Share file or directory",
            Callback: shareHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "discoverable",
                        Patterns: []string{"--discoverable"},
                        Description: "Make file discoverable by search engines",
                        OmitValue: true,
                    },
                    cli.StringFlag{
                        Name: "role",
                        Patterns: []string{"--role"},
                        Description: fmt.Sprintf("Share role. Default: %s", DefaultShareRole),
                        DefaultValue: DefaultShareRole,
                    },
                    cli.StringFlag{
                        Name: "type",
                        Patterns: []string{"--type"},
                        Description: fmt.Sprintf("Share type. Default: %s", DefaultShareType),
                        DefaultValue: DefaultShareType,
                    },
                    cli.StringFlag{
                        Name: "email",
                        Patterns: []string{"--email"},
                        Description: "The email address of the user or group to share the file with. Requires 'user' or 'group' as type",
                    },
                    cli.BoolFlag{
                        Name: "revoke",
                        Patterns: []string{"--revoke"},
                        Description: "Delete all sharing permissions",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] delete <id>",
            Description: "Delete file or directory",
            Callback: deleteHandler,
            Flags: cli.Flags{
                "global": globalFlags,
            },
        },
        &cli.Handler{
            Pattern: "[global] delete revision <fileId> <revisionId>",
            Description: "Delete file revision",
            Callback: deleteRevisionHandler,
            Flags: cli.Flags{
                "global": globalFlags,
            },
        },
        &cli.Handler{
            Pattern: "[global] about [options]",
            Description: "Google drive metadata, quota usage",
            Callback: aboutHandler,
            Flags: cli.Flags{
                "global": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "sizeInBytes",
                        Patterns: []string{"--bytes"},
                        Description: "Show size in bytes",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global] about import",
            Description: "Show supported import formats",
            Callback: aboutImportHandler,
            Flags: cli.Flags{
                "global": globalFlags,
            },
        },
        &cli.Handler{
            Pattern: "[global] about export",
            Description: "Show supported export formats",
            Callback: aboutExportHandler,
            Flags: cli.Flags{
                "global": globalFlags,
            },
        },
        &cli.Handler{
            Pattern: "version",
            Description: "Print application version",
            Callback: printVersion,
        },
        &cli.Handler{
            Pattern: "help",
            Description: "Print help",
            Callback: printHelp,
        },
        &cli.Handler{
            Pattern: "help <command>",
            Description: "Print command help",
            Callback: printCommandHelp,
        },
        &cli.Handler{
            Pattern: "help <command> <subcommand>",
            Description: "Print subcommand help",
            Callback: printSubCommandHelp,
        },
    }

    cli.SetHandlers(handlers)

    if ok := cli.Handle(os.Args[1:]); !ok {
        ExitF("No valid arguments given, use '%s help' to see available commands", Name)
    }
}
