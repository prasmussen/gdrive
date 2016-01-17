package main

import (
	"fmt"
	"os"
    "./cli"
)

const Name = "gdrive"
const Version = "2.0.0"

const ClientId     = "367116221053-7n0vf5akeru7on6o2fjinrecpdoe99eg.apps.googleusercontent.com"
const ClientSecret = "1qsNodXNaWq1mQuBjUjmvhoO"

const DefaultMaxFiles = 100
const DefaultQuery = "trashed = false and 'me' in owners"

var DefaultConfigDir = GetDefaultConfigDir()
var DefaultTokenFilePath = GetDefaultTokenFilePath()


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
            Pattern: "[global options] list [options]",
            Description: "List files",
            Callback: listHandler,
            Flags: cli.Flags{
                "global options": globalFlags,
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
            Pattern: "[global options] download [options] <id>",
            Description: "Download file or directory",
            Callback: downloadHandler,
            Flags: cli.Flags{
                "global options": globalFlags,
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
            Pattern: "[global options] upload [options] <path>",
            Description: "Upload file or directory",
            Callback: uploadHandler,
            Flags: cli.Flags{
                "global options": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "recursive",
                        Patterns: []string{"-r", "--recursive"},
                        Description: "Upload directory recursively",
                        OmitValue: true,
                    },
                    cli.StringFlag{
                        Name: "parent",
                        Patterns: []string{"-p", "--parent"},
                        Description: "Parent id, used to upload file to a specific directory",
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
                },
            },
        },
        &cli.Handler{
            Pattern: "[global options] info [options] <id>",
            Description: "Show file info",
            Callback: infoHandler,
            Flags: cli.Flags{
                "global options": globalFlags,
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
            Pattern: "[global options] mkdir [options] <name>",
            Description: "Create directory",
            Callback: handler,
            Flags: cli.Flags{
                "global options": globalFlags,
                "options": []cli.Flag{
                    cli.StringFlag{
                        Name: "parent",
                        Patterns: []string{"-p", "--parent"},
                        Description: "Parent id of created directory",
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
            Pattern: "[global options] share <id>",
            Description: "Share file or directory",
            Callback: handler,
            Flags: cli.Flags{
                "global options": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "revoke",
                        Patterns: []string{"--revoke"},
                        Description: "Unshare file or directory",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global options] url [options] <id>",
            Description: "Get url to file or directory",
            Callback: handler,
            Flags: cli.Flags{
                "global options": globalFlags,
                "options": []cli.Flag{
                    cli.BoolFlag{
                        Name: "download",
                        Patterns: []string{"--download"},
                        Description: "Download url",
                        OmitValue: true,
                    },
                },
            },
        },
        &cli.Handler{
            Pattern: "[global options] delete <id>",
            Description: "Delete file or directory",
            Callback: deleteHandler,
            Flags: cli.Flags{
                "global options": globalFlags,
            },
        },
        &cli.Handler{
            Pattern: "[global options] quota [options]",
            Description: "Show free space",
            Callback: handler,
            Flags: cli.Flags{
                "global options": globalFlags,
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
            Pattern: "help <subcommand>",
            Description: "Print subcommand help",
            Callback: printCommandHelp,
        },
    }

    cli.SetHandlers(handlers)

    if ok := cli.Handle(os.Args[1:]); !ok {
        ExitF("No valid arguments given, use '%s help' to see available commands", Name)
    }
}
