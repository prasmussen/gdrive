package main

import (
	"fmt"
	"os"

	"github.com/prasmussen/gdrive/cli"
)

const Name = "gdrive"
const Version = "2.1.1"

const DefaultMaxFiles = 30
const DefaultMaxChanges = 100
const DefaultNameWidth = 40
const DefaultPathWidth = 60
const DefaultUploadChunkSize = 8 * 1024 * 1024
const DefaultTimeout = 5 * 60
const DefaultQuery = "trashed = false and 'me' in owners"
const DefaultShareRole = "reader"
const DefaultShareType = "anyone"

var DefaultConfigDir = GetDefaultConfigDir()

func main() {
	globalFlags := []cli.Flag{
		cli.StringFlag{
			Name:         "configDir",
			Patterns:     []string{"-c", "--config"},
			Description:  fmt.Sprintf("Application path, default: %s", DefaultConfigDir),
			DefaultValue: DefaultConfigDir,
		},
		cli.StringFlag{
			Name:        "refreshToken",
			Patterns:    []string{"--refresh-token"},
			Description: "Oauth refresh token used to get access token (for advanced users)",
		},
		cli.StringFlag{
			Name:        "accessToken",
			Patterns:    []string{"--access-token"},
			Description: "Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)",
		},
		cli.StringFlag{
			Name:        "serviceAccount",
			Patterns:    []string{"--service-account"},
			Description: "Oauth service account filename, used for server to server communication without user interaction (filename path is relative to config dir)",
		},
	}

	handlers := []*cli.Handler{
		&cli.Handler{
			Pattern:     "[global] list [options]",
			Description: "List files",
			Callback:    listHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.IntFlag{
						Name:         "maxFiles",
						Patterns:     []string{"-m", "--max"},
						Description:  fmt.Sprintf("Max files to list, default: %d", DefaultMaxFiles),
						DefaultValue: DefaultMaxFiles,
					},
					cli.StringFlag{
						Name:         "query",
						Patterns:     []string{"-q", "--query"},
						Description:  fmt.Sprintf(`Default query: "%s". See https://developers.google.com/drive/search-parameters`, DefaultQuery),
						DefaultValue: DefaultQuery,
					},
					cli.StringFlag{
						Name:        "sortOrder",
						Patterns:    []string{"--order"},
						Description: "Sort order. See https://godoc.org/google.golang.org/api/drive/v3#FilesListCall.OrderBy",
					},
					cli.IntFlag{
						Name:         "nameWidth",
						Patterns:     []string{"--name-width"},
						Description:  fmt.Sprintf("Width of name column, default: %d, minimum: 9, use 0 for full width", DefaultNameWidth),
						DefaultValue: DefaultNameWidth,
					},
					cli.BoolFlag{
						Name:        "absPath",
						Patterns:    []string{"--absolute"},
						Description: "Show absolute path to file (will only show path from first parent)",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "skipHeader",
						Patterns:    []string{"--no-header"},
						Description: "Dont print the header",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "sizeInBytes",
						Patterns:    []string{"--bytes"},
						Description: "Size in bytes",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] download [options] <fileId>",
			Description: "Download file or directory",
			Callback:    downloadHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "force",
						Patterns:    []string{"-f", "--force"},
						Description: "Overwrite existing file",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "skip",
						Patterns:    []string{"-s", "--skip"},
						Description: "Skip existing files",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "recursive",
						Patterns:    []string{"-r", "--recursive"},
						Description: "Download directory recursively, documents will be skipped",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "path",
						Patterns:    []string{"--path"},
						Description: "Download path",
					},
					cli.BoolFlag{
						Name:        "delete",
						Patterns:    []string{"--delete"},
						Description: "Delete remote file when download is successful",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "stdout",
						Patterns:    []string{"--stdout"},
						Description: "Write file content to stdout",
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", DefaultTimeout),
						DefaultValue: DefaultTimeout,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] download query [options] <query>",
			Description: "Download all files and directories matching query",
			Callback:    downloadQueryHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "force",
						Patterns:    []string{"-f", "--force"},
						Description: "Overwrite existing file",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "skip",
						Patterns:    []string{"-s", "--skip"},
						Description: "Skip existing files",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "recursive",
						Patterns:    []string{"-r", "--recursive"},
						Description: "Download directories recursively, documents will be skipped",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "path",
						Patterns:    []string{"--path"},
						Description: "Download path",
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] upload [options] <path>",
			Description: "Upload file or directory",
			Callback:    uploadHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "recursive",
						Patterns:    []string{"-r", "--recursive"},
						Description: "Upload directory recursively",
						OmitValue:   true,
					},
					cli.StringSliceFlag{
						Name:        "parent",
						Patterns:    []string{"-p", "--parent"},
						Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
					},
					cli.StringFlag{
						Name:        "name",
						Patterns:    []string{"--name"},
						Description: "Filename",
					},
					cli.StringFlag{
						Name:        "description",
						Patterns:    []string{"--description"},
						Description: "File description",
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "mime",
						Patterns:    []string{"--mime"},
						Description: "Force mime type",
					},
					cli.BoolFlag{
						Name:        "share",
						Patterns:    []string{"--share"},
						Description: "Share file",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "delete",
						Patterns:    []string{"--delete"},
						Description: "Delete local file when upload is successful",
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", DefaultTimeout),
						DefaultValue: DefaultTimeout,
					},
					cli.IntFlag{
						Name:         "chunksize",
						Patterns:     []string{"--chunksize"},
						Description:  fmt.Sprintf("Set chunk size in bytes, default: %d", DefaultUploadChunkSize),
						DefaultValue: DefaultUploadChunkSize,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] upload - [options] <name>",
			Description: "Upload file from stdin",
			Callback:    uploadStdinHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringSliceFlag{
						Name:        "parent",
						Patterns:    []string{"-p", "--parent"},
						Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
					},
					cli.IntFlag{
						Name:         "chunksize",
						Patterns:     []string{"--chunksize"},
						Description:  fmt.Sprintf("Set chunk size in bytes, default: %d", DefaultUploadChunkSize),
						DefaultValue: DefaultUploadChunkSize,
					},
					cli.StringFlag{
						Name:        "description",
						Patterns:    []string{"--description"},
						Description: "File description",
					},
					cli.StringFlag{
						Name:        "mime",
						Patterns:    []string{"--mime"},
						Description: "Force mime type",
					},
					cli.BoolFlag{
						Name:        "share",
						Patterns:    []string{"--share"},
						Description: "Share file",
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", DefaultTimeout),
						DefaultValue: DefaultTimeout,
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] update [options] <fileId> <path>",
			Description: "Update file, this creates a new revision of the file",
			Callback:    updateHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringSliceFlag{
						Name:        "parent",
						Patterns:    []string{"-p", "--parent"},
						Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
					},
					cli.StringFlag{
						Name:        "name",
						Patterns:    []string{"--name"},
						Description: "Filename",
					},
					cli.StringFlag{
						Name:        "description",
						Patterns:    []string{"--description"},
						Description: "File description",
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "mime",
						Patterns:    []string{"--mime"},
						Description: "Force mime type",
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", DefaultTimeout),
						DefaultValue: DefaultTimeout,
					},
					cli.IntFlag{
						Name:         "chunksize",
						Patterns:     []string{"--chunksize"},
						Description:  fmt.Sprintf("Set chunk size in bytes, default: %d", DefaultUploadChunkSize),
						DefaultValue: DefaultUploadChunkSize,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] info [options] <fileId>",
			Description: "Show file info",
			Callback:    infoHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "sizeInBytes",
						Patterns:    []string{"--bytes"},
						Description: "Show size in bytes",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] mkdir [options] <name>",
			Description: "Create directory",
			Callback:    mkdirHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringSliceFlag{
						Name:        "parent",
						Patterns:    []string{"-p", "--parent"},
						Description: "Parent id of created directory, can be specified multiple times to give many parents",
					},
					cli.StringFlag{
						Name:        "description",
						Patterns:    []string{"--description"},
						Description: "Directory description",
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] share [options] <fileId>",
			Description: "Share file or directory",
			Callback:    shareHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringFlag{
						Name:         "role",
						Patterns:     []string{"--role"},
						Description:  fmt.Sprintf("Share role: owner/writer/commenter/reader, default: %s", DefaultShareRole),
						DefaultValue: DefaultShareRole,
					},
					cli.StringFlag{
						Name:         "type",
						Patterns:     []string{"--type"},
						Description:  fmt.Sprintf("Share type: user/group/domain/anyone, default: %s", DefaultShareType),
						DefaultValue: DefaultShareType,
					},
					cli.StringFlag{
						Name:        "email",
						Patterns:    []string{"--email"},
						Description: "The email address of the user or group to share the file with. Requires 'user' or 'group' as type",
					},
					cli.StringFlag{
						Name:        "domain",
						Patterns:    []string{"--domain"},
						Description: "The name of Google Apps domain. Requires 'domain' as type",
					},
					cli.BoolFlag{
						Name:        "discoverable",
						Patterns:    []string{"--discoverable"},
						Description: "Make file discoverable by search engines",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "revoke",
						Patterns:    []string{"--revoke"},
						Description: "Delete all sharing permissions (owner roles will be skipped)",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] share list <fileId>",
			Description: "List files permissions",
			Callback:    shareListHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
			},
		},
		&cli.Handler{
			Pattern:     "[global] share revoke <fileId> <permissionId>",
			Description: "Revoke permission",
			Callback:    shareRevokeHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
			},
		},
		&cli.Handler{
			Pattern:     "[global] delete [options] <fileId>",
			Description: "Delete file or directory",
			Callback:    deleteHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "recursive",
						Patterns:    []string{"-r", "--recursive"},
						Description: "Delete directory and all it's content",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] sync list [options]",
			Description: "List all syncable directories on drive",
			Callback:    listSyncHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "skipHeader",
						Patterns:    []string{"--no-header"},
						Description: "Dont print the header",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] sync content [options] <fileId>",
			Description: "List content of syncable directory",
			Callback:    listRecursiveSyncHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringFlag{
						Name:        "sortOrder",
						Patterns:    []string{"--order"},
						Description: "Sort order. See https://godoc.org/google.golang.org/api/drive/v3#FilesListCall.OrderBy",
					},
					cli.IntFlag{
						Name:         "pathWidth",
						Patterns:     []string{"--path-width"},
						Description:  fmt.Sprintf("Width of path column, default: %d, minimum: 9, use 0 for full width", DefaultPathWidth),
						DefaultValue: DefaultPathWidth,
					},
					cli.BoolFlag{
						Name:        "skipHeader",
						Patterns:    []string{"--no-header"},
						Description: "Dont print the header",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "sizeInBytes",
						Patterns:    []string{"--bytes"},
						Description: "Size in bytes",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] sync download [options] <fileId> <path>",
			Description: "Sync drive directory to local directory",
			Callback:    downloadSyncHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "keepRemote",
						Patterns:    []string{"--keep-remote"},
						Description: "Keep remote file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "keepLocal",
						Patterns:    []string{"--keep-local"},
						Description: "Keep local file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "keepLargest",
						Patterns:    []string{"--keep-largest"},
						Description: "Keep largest file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "deleteExtraneous",
						Patterns:    []string{"--delete-extraneous"},
						Description: "Delete extraneous local files",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "dryRun",
						Patterns:    []string{"--dry-run"},
						Description: "Show what would have been transferred",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", DefaultTimeout),
						DefaultValue: DefaultTimeout,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] sync upload [options] <path> <fileId>",
			Description: "Sync local directory to drive",
			Callback:    uploadSyncHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "keepRemote",
						Patterns:    []string{"--keep-remote"},
						Description: "Keep remote file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "keepLocal",
						Patterns:    []string{"--keep-local"},
						Description: "Keep local file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "keepLargest",
						Patterns:    []string{"--keep-largest"},
						Description: "Keep largest file when a conflict is encountered",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "deleteExtraneous",
						Patterns:    []string{"--delete-extraneous"},
						Description: "Delete extraneous remote files",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "dryRun",
						Patterns:    []string{"--dry-run"},
						Description: "Show what would have been transferred",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", DefaultTimeout),
						DefaultValue: DefaultTimeout,
					},
					cli.IntFlag{
						Name:         "chunksize",
						Patterns:     []string{"--chunksize"},
						Description:  fmt.Sprintf("Set chunk size in bytes, default: %d", DefaultUploadChunkSize),
						DefaultValue: DefaultUploadChunkSize,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] changes [options]",
			Description: "List file changes",
			Callback:    listChangesHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.IntFlag{
						Name:         "maxChanges",
						Patterns:     []string{"-m", "--max"},
						Description:  fmt.Sprintf("Max changes to list, default: %d", DefaultMaxChanges),
						DefaultValue: DefaultMaxChanges,
					},
					cli.StringFlag{
						Name:         "pageToken",
						Patterns:     []string{"--since"},
						Description:  fmt.Sprintf("Page token to start listing changes from"),
						DefaultValue: "1",
					},
					cli.BoolFlag{
						Name:        "now",
						Patterns:    []string{"--now"},
						Description: fmt.Sprintf("Get latest page token"),
						OmitValue:   true,
					},
					cli.IntFlag{
						Name:         "nameWidth",
						Patterns:     []string{"--name-width"},
						Description:  fmt.Sprintf("Width of name column, default: %d, minimum: 9, use 0 for full width", DefaultNameWidth),
						DefaultValue: DefaultNameWidth,
					},
					cli.BoolFlag{
						Name:        "skipHeader",
						Patterns:    []string{"--no-header"},
						Description: "Dont print the header",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] revision list [options] <fileId>",
			Description: "List file revisions",
			Callback:    listRevisionsHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.IntFlag{
						Name:         "nameWidth",
						Patterns:     []string{"--name-width"},
						Description:  fmt.Sprintf("Width of name column, default: %d, minimum: 9, use 0 for full width", DefaultNameWidth),
						DefaultValue: DefaultNameWidth,
					},
					cli.BoolFlag{
						Name:        "skipHeader",
						Patterns:    []string{"--no-header"},
						Description: "Dont print the header",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "sizeInBytes",
						Patterns:    []string{"--bytes"},
						Description: "Size in bytes",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] revision download [options] <fileId> <revId>",
			Description: "Download revision",
			Callback:    downloadRevisionHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "force",
						Patterns:    []string{"-f", "--force"},
						Description: "Overwrite existing file",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.BoolFlag{
						Name:        "stdout",
						Patterns:    []string{"--stdout"},
						Description: "Write file content to stdout",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "path",
						Patterns:    []string{"--path"},
						Description: "Download path",
					},
					cli.IntFlag{
						Name:         "timeout",
						Patterns:     []string{"--timeout"},
						Description:  fmt.Sprintf("Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: %d", DefaultTimeout),
						DefaultValue: DefaultTimeout,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] revision delete <fileId> <revId>",
			Description: "Delete file revision",
			Callback:    deleteRevisionHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
			},
		},
		&cli.Handler{
			Pattern:     "[global] import [options] <path>",
			Description: "Upload and convert file to a google document, see 'about import' for available conversions",
			Callback:    importHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.StringSliceFlag{
						Name:        "parent",
						Patterns:    []string{"-p", "--parent"},
						Description: "Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents",
					},
					cli.BoolFlag{
						Name:        "noProgress",
						Patterns:    []string{"--no-progress"},
						Description: "Hide progress",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "mime",
						Patterns:    []string{"--mime"},
						Description: "Mime type of imported file",
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] export [options] <fileId>",
			Description: "Export a google document",
			Callback:    exportHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "force",
						Patterns:    []string{"-f", "--force"},
						Description: "Overwrite existing file",
						OmitValue:   true,
					},
					cli.StringFlag{
						Name:        "mime",
						Patterns:    []string{"--mime"},
						Description: "Mime type of exported file",
					},
					cli.BoolFlag{
						Name:        "printMimes",
						Patterns:    []string{"--print-mimes"},
						Description: "Print available mime types for given file",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] about [options]",
			Description: "Google drive metadata, quota usage",
			Callback:    aboutHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
				cli.NewFlagGroup("options",
					cli.BoolFlag{
						Name:        "sizeInBytes",
						Patterns:    []string{"--bytes"},
						Description: "Show size in bytes",
						OmitValue:   true,
					},
				),
			},
		},
		&cli.Handler{
			Pattern:     "[global] about import",
			Description: "Show supported import formats",
			Callback:    aboutImportHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
			},
		},
		&cli.Handler{
			Pattern:     "[global] about export",
			Description: "Show supported export formats",
			Callback:    aboutExportHandler,
			FlagGroups: cli.FlagGroups{
				cli.NewFlagGroup("global", globalFlags...),
			},
		},
		&cli.Handler{
			Pattern:     "version",
			Description: "Print application version",
			Callback:    printVersion,
		},
		&cli.Handler{
			Pattern:     "help",
			Description: "Print help",
			Callback:    printHelp,
		},
		&cli.Handler{
			Pattern:     "help <command>",
			Description: "Print command help",
			Callback:    printCommandHelp,
		},
		&cli.Handler{
			Pattern:     "help <command> <subcommand>",
			Description: "Print subcommand help",
			Callback:    printSubCommandHelp,
		},
	}

	cli.SetHandlers(handlers)

	if ok := cli.Handle(os.Args[1:]); !ok {
		ExitF("No valid arguments given, use '%s help' to see available commands", Name)
	}
}
