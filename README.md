gdrive
======


## Note
~~This tool is not being actively maintained at the moment, ymmv~~ **Active maintenance will be resumed soon.**
For incremental and encrypted backup of unix systems [borg](https://github.com/borgbackup/borg)
is a great alternative and [rsync.net](http://rsync.net/products/attic.html) provides
a cheap and reliable backup target.

## Overview
gdrive is a command line utility for interacting with Google Drive.

## Prerequisites
None, binaries are statically linked.
If you want to compile from source you need the [go toolchain](http://golang.org/doc/install).
Version 1.5 or higher.

## Installation
### With [Homebrew](http://brew.sh) on Mac
```
brew install gdrive
```
### Other
Download `gdrive` from one of the links below. On unix systems
run `chmod +x gdrive` after download to make the binary executable.
The first time gdrive is launched (i.e. run `gdrive about` in your
terminal not just `gdrive`), you will be prompted for a verification code.
The code is obtained by following the printed url and authenticating with the
google account for the drive you want access to. This will create a token file
inside the .gdrive folder in your home directory. Note that anyone with access
to this file will also have access to your google drive.
If you want to manage multiple drives you can use the global `--config` flag
or set the environment variable `GDRIVE_CONFIG_DIR`.
Example: `GDRIVE_CONFIG_DIR="/home/user/.gdrive-secondary" gdrive list`
You will be prompted for a new verification code if the folder does not exist.

### Downloads
Check our [releases](https://github.com/gdrive-org/gdrive/releases) section.

| Filename               | Version | Description        | Shasum                                   |
|:-----------------------|:--------|:-------------------|:-----------------------------------------|
| [gdrive-osx-x64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-osx-x64) | 2.1.0 | OS X 64-bit | 2d817cab0912e0efa5529ada7a4bb73fe90d2d74 |
| [gdrive-osx-386](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-osx-386) | 2.1.0 | OS X 32-bit | ed1a12998270f8176c96729d80f3f7307a2e2db6 |
| [gdrive-linux-x64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-linux-x64) | 2.1.0 | Linux 64-bit | 6ef7c740e980358a9a5de36c1aac7ea375319aa3 |
| [gdrive-linux-386](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-linux-386) | 2.1.0 | Linux 32-bit | 03f423fc7a444e95de3998d24dbfbdb0335fd055 |
| [gdrive-linux-rpi](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-linux-rpi) | 2.1.0 | Linux Raspberry Pi | 0247d36e2591ea009ddbd003300bf10cb58089da |
| [gdrive-linux-arm64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-linux-arm64) | 2.1.0 | Linux arm 64-bit | 3a41c27e10587eef0be20fe00415c65723fa3aee |
| [gdrive-linux-arm](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-linux-arm) | 2.1.0 | Linux arm 32-bit | 0247d36e2591ea009ddbd003300bf10cb58089da |
| [gdrive-linux-mips64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-linux-mips64) | 2.1.0 | Linux mips 64-bit | 69bdf4b9a4ae1f9ead9d35ba118a70005bea3829 |
| [gdrive-linux-mips64le](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-linux-mips64le) | 2.1.0 | Linux mips 64-bit le | 421ae3662ed9b486649560512e3c8a43dbf6cdc6 |
| [gdrive-linux-ppc64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-linux-ppc64) | 2.1.0 | Linux PPC 64-bit | 9b1621d1656fa8e4f2be21202708ab3bb8c0aa77 |
| [gdrive-linux-ppc64le](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-linux-ppc64le) | 2.1.0 | Linux PPC 64-bit le | 0bf501abb505f851a7592456db304f2c4430e9cd |
| [gdrive-windows-386.exe](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-windows-386.exe) | 2.1.0 | Window 32-bit | 4974290122f635d9c71875bc2c760e46b943c5ab |
| [gdrive-windows-x64.exe](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-windows-x64.exe) | 2.1.0 | Windows 64-bit | 53780f9a4168c71fd9f3d429932207069d1dcf03 |
| [gdrive-dragonfly-x64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-dragonfly-x64) | 2.1.0 | DragonFly BSD 64-bit | 5871ee54f992b71b3a1586ec224d93a79265a8b4 |
| [gdrive-freebsd-x64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-freebsd-x64) | 2.1.0 | FreeBSD 64-bit | ad4d16b0144b8fd05ef7f3174a54fe39b04a766d |
| [gdrive-freebsd-386](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-freebsd-386) | 2.1.0 | FreeBSD 32-bit | cc636c839434689b73274e3574ac8477c9542880 |
| [gdrive-freebsd-arm](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-freebsd-arm) | 2.1.0 | FreeBSD arm | 19abe4b23ad7dd968d0089313eca3108601dca02 |
| [gdrive-netbsd-x64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-netbsd-x64) | 2.1.0 | NetBSD 64-bit | 07af67d4a149a926ee9ccf82591a593eb2330c36 |
| [gdrive-netbsd-386](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-netbsd-386) | 2.1.0 | NetBSD 32-bit | 95a090ff3ecb8b6e2180afb02bec6ba7419a33e0 |
| [gdrive-netbsd-arm](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-netbsd-arm) | 2.1.0 | NetBSD arm | b953d575f2c306053e31718050698c47388dea26 |
| [gdrive-openbsd-x64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-openbsd-x64) | 2.1.0 | OpenBSD 64-bit | ac8dbed27ffd7ed0b13b2abb8be4b5ca72da9d7c |
| [gdrive-openbsd-386](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-openbsd-386) | 2.1.0 | OpenBSD 32-bit | d7a2d1441d3b41e36da080643c3d5f6306181eae |
| [gdrive-openbsd-arm](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-openbsd-arm) | 2.1.0 | OpenBSD arm | 9e508032367471515199f6787aa771e271d307f4 |
| [gdrive-solaris-x64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-solaris-x64) | 2.1.0 | Solaris 64-bit | c18612da275065064aa650b669de7140dcd94a5f |
| [gdrive-plan9-x64](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-plan9-x64) | 2.1.0 | Plan9 64-bit | 8907f5ffb8c5d6b2e15cd9d8f5a51ef762173298 |
| [gdrive-plan9-386](https://github.com/gdrive-org/gdrive/releases/download/2.1.0/gdrive-plan9-386) | 2.1.0 | Plan9 32-bit | f60b749ba57b8b2d824d06a7bd788a8c5808a607 |

## Compile from source
```bash
go get github.com/prasmussen/gdrive
```
The gdrive binary should now be available at `$GOPATH/bin/gdrive`


## Gdrive 2
Gdrive 2 is more or less a full rewrite and is not backwards compatible
with gdrive 1 as all the command line arguments has changed slightly.
Gdrive 2 uses version 3 of the google drive api and my google-api-go-client
fork is no longer needed.

### Syncing
Gdrive 2 supports basic syncing. It only syncs one way at the time and works
more like rsync than e.g. dropbox. Files that are synced to google drive
are tagged with an appProperty so that the files on drive can be traversed
faster. This means that you can't upload files with `gdrive upload` into
a sync directory as the files would be missing the sync tag, and would be
ignored by the sync commands.
The current implementation is slow and uses a lot of memory if you are
syncing many files. Currently only one file is uploaded at the time,
the speed can be improved in the future by uploading several files concurrently.
To learn more see usage and the examples below.

### Service Account
For server to server communication, where user interaction is not a viable option, 
is it possible to use a service account, as described in this [Google document](https://developers.google.com/identity/protocols/OAuth2ServiceAccount).
If you want to use a service account, instead of being interactively prompted for
authentication, you need to use the `--service-account <serviceAccountCredentials>` 
global option, where `serviceAccountCredentials` is a file in JSON format obtained
through the Google API Console, and its location is relative to the config dir. 

#### .gdriveignore
Placing a .gdriveignore in the root of your sync directory can be used to
skip certain files from being synced. .gdriveignore follows the same
rules as [.gitignore](https://git-scm.com/docs/gitignore), except that gdrive only reads the .gdriveignore file in the root of the sync directory, not ones in any subdirectories.


## Usage
```
gdrive [global] list [options]                                 List files
gdrive [global] download [options] <fileId>                    Download file or directory
gdrive [global] download query [options] <query>               Download all files and directories matching query
gdrive [global] upload [options] <path>                        Upload file or directory
gdrive [global] upload - [options] <name>                      Upload file from stdin
gdrive [global] update [options] <fileId> <path>               Update file, this creates a new revision of the file
gdrive [global] info [options] <fileId>                        Show file info
gdrive [global] mkdir [options] <name>                         Create directory
gdrive [global] share [options] <fileId>                       Share file or directory
gdrive [global] share list <fileId>                            List files permissions
gdrive [global] share revoke <fileId> <permissionId>           Revoke permission
gdrive [global] delete [options] <fileId>                      Delete file or directory
gdrive [global] sync list [options]                            List all syncable directories on drive
gdrive [global] sync content [options] <fileId>                List content of syncable directory
gdrive [global] sync download [options] <fileId> <path>        Sync drive directory to local directory
gdrive [global] sync upload [options] <path> <fileId>          Sync local directory to drive
gdrive [global] changes [options]                              List file changes
gdrive [global] revision list [options] <fileId>               List file revisions
gdrive [global] revision download [options] <fileId> <revId>   Download revision
gdrive [global] revision delete <fileId> <revId>               Delete file revision
gdrive [global] import [options] <path>                        Upload and convert file to a google document, see 'about import' for available conversions
gdrive [global] export [options] <fileId>                      Export a google document
gdrive [global] about [options]                                Google drive metadata, quota usage
gdrive [global] about import                                   Show supported import formats
gdrive [global] about export                                   Show supported export formats
gdrive version                                                 Print application version
gdrive help                                                    Print help
gdrive help <command>                                          Print command help
gdrive help <command> <subcommand>                             Print subcommand help
```

#### List files
```
gdrive [global] list [options]

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)

options:
  -m, --max <maxFiles>       Max files to list, default: 30
  -q, --query <query>        Default query: "trashed = false and 'me' in owners". See https://developers.google.com/drive/search-parameters
  --order <sortOrder>        Sort order. See https://godoc.org/google.golang.org/api/drive/v3#FilesListCall.OrderBy
  --name-width <nameWidth>   Width of name column, default: 40, minimum: 9, use 0 for full width
  --absolute                 Show absolute path to file (will only show path from first parent)
  --no-header                Dont print the header
  --bytes                    Size in bytes
```

List file in subdirectory


```
./gdrive list --query " 'IdOfTheParentFolder' in parents"
```

#### Download file or directory
```
gdrive [global] download [options] <fileId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -f, --force           Overwrite existing file
  -r, --recursive       Download directory recursively, documents will be skipped
  --path <path>         Download path
  --delete              Delete remote file when download is successful
  --no-progress         Hide progress
  --stdout              Write file content to stdout
  --timeout <timeout>   Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: 300
```

#### Download all files and directories matching query
```
gdrive [global] download query [options] <query>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -f, --force       Overwrite existing file
  -r, --recursive   Download directories recursively, documents will be skipped
  --path <path>     Download path
  --no-progress     Hide progress
```

#### Upload file or directory
```
gdrive [global] upload [options] <path>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -r, --recursive               Upload directory recursively
  -p, --parent <parent>         Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents
  --name <name>                 Filename
  --description <description>   File description
  --no-progress                 Hide progress
  --mime <mime>                 Force mime type
  --share                       Share file
  --delete                      Delete local file when upload is successful
  --timeout <timeout>           Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: 300
  --chunksize <chunksize>       Set chunk size in bytes, default: 8388608
```

#### Upload file from stdin
```
gdrive [global] upload - [options] <name>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -p, --parent <parent>         Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents
  --chunksize <chunksize>       Set chunk size in bytes, default: 8388608
  --description <description>   File description
  --mime <mime>                 Force mime type
  --share                       Share file
  --timeout <timeout>           Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: 300
  --no-progress                 Hide progress
```

#### Update file, this creates a new revision of the file
```
gdrive [global] update [options] <fileId> <path>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -p, --parent <parent>         Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents
  --name <name>                 Filename
  --description <description>   File description
  --no-progress                 Hide progress
  --mime <mime>                 Force mime type
  --timeout <timeout>           Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: 300
  --chunksize <chunksize>       Set chunk size in bytes, default: 8388608
```

#### Show file info
```
gdrive [global] info [options] <fileId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  --bytes   Show size in bytes
```

#### Create directory
```
gdrive [global] mkdir [options] <name>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -p, --parent <parent>         Parent id of created directory, can be specified multiple times to give many parents
  --description <description>   Directory description
```

#### Share file or directory
```
gdrive [global] share [options] <fileId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  --role <role>     Share role: owner/writer/commenter/reader, default: reader
  --type <type>     Share type: user/group/domain/anyone, default: anyone
  --email <email>   The email address of the user or group to share the file with. Requires 'user' or 'group' as type
  --discoverable    Make file discoverable by search engines
  --revoke          Delete all sharing permissions (owner roles will be skipped)
```

#### List files permissions
```
gdrive [global] share list <fileId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
```

#### Revoke permission
```
gdrive [global] share revoke <fileId> <permissionId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
```

#### Delete file or directory
```
gdrive [global] delete [options] <fileId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -r, --recursive   Delete directory and all it's content
```

#### List all syncable directories on drive
```
gdrive [global] sync list [options]

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  --no-header   Dont print the header
```

#### List content of syncable directory
```
gdrive [global] sync content [options] <fileId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  --order <sortOrder>        Sort order. See https://godoc.org/google.golang.org/api/drive/v3#FilesListCall.OrderBy
  --path-width <pathWidth>   Width of path column, default: 60, minimum: 9, use 0 for full width
  --no-header                Dont print the header
  --bytes                    Size in bytes
```

#### Sync drive directory to local directory
```
gdrive [global] sync download [options] <fileId> <path>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  --keep-remote         Keep remote file when a conflict is encountered
  --keep-local          Keep local file when a conflict is encountered
  --keep-largest        Keep largest file when a conflict is encountered
  --delete-extraneous   Delete extraneous local files
  --dry-run             Show what would have been transferred
  --no-progress         Hide progress
  --timeout <timeout>   Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: 300
```

#### Sync local directory to drive
```
gdrive [global] sync upload [options] <path> <fileId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  --keep-remote             Keep remote file when a conflict is encountered
  --keep-local              Keep local file when a conflict is encountered
  --keep-largest            Keep largest file when a conflict is encountered
  --delete-extraneous       Delete extraneous remote files
  --dry-run                 Show what would have been transferred
  --no-progress             Hide progress
  --timeout <timeout>       Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: 300
  --chunksize <chunksize>   Set chunk size in bytes, default: 8388608
```

#### List file changes
```
gdrive [global] changes [options]

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -m, --max <maxChanges>     Max changes to list, default: 100
  --since <pageToken>        Page token to start listing changes from
  --now                      Get latest page token
  --name-width <nameWidth>   Width of name column, default: 40, minimum: 9, use 0 for full width
  --no-header                Dont print the header
```

#### List file revisions
```
gdrive [global] revision list [options] <fileId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  --name-width <nameWidth>   Width of name column, default: 40, minimum: 9, use 0 for full width
  --no-header                Dont print the header
  --bytes                    Size in bytes
```

#### Download revision
```
gdrive [global] revision download [options] <fileId> <revId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -f, --force           Overwrite existing file
  --no-progress         Hide progress
  --stdout              Write file content to stdout
  --path <path>         Download path
  --timeout <timeout>   Set timeout in seconds, use 0 for no timeout. Timeout is reached when no data is transferred in set amount of seconds, default: 300
```

#### Delete file revision
```
gdrive [global] revision delete <fileId> <revId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
```

#### Upload and convert file to a google document, see 'about import' for available conversions
```
gdrive [global] import [options] <path>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -p, --parent <parent>   Parent id, used to upload file to a specific directory, can be specified multiple times to give many parents
  --no-progress           Hide progress
```

#### Export a google document
```
gdrive [global] export [options] <fileId>

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  -f, --force     Overwrite existing file
  --mime <mime>   Mime type of exported file
  --print-mimes   Print available mime types for given file
```

#### Google drive metadata, quota usage
```
gdrive [global] about [options]

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
  
options:
  --bytes   Show size in bytes
```

#### Show supported import formats
```
gdrive [global] about import

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
```

#### Show supported export formats
```
gdrive [global] about export

global:
  -c, --config <configDir>         Application path, default: /Users/<user>/.gdrive
  --refresh-token <refreshToken>   Oauth refresh token used to get access token (for advanced users)
  --access-token <accessToken>     Oauth access token, only recommended for short-lived requests because of short lifetime (for advanced users)
  --service-account <accountFile>  Oauth service account filename, used for server to server communication without user interaction (file is relative to config dir)
```


## Examples
#### List files
```
$ gdrive list
Id                             Name                    Type   Size     Created
0B3X9GlR6EmbnZ3gyeGw4d3ozbUk   drive-windows-x64.exe   bin    6.6 MB   2015-07-18 16:43:58
0B3X9GlR6EmbnTXlSc1FqV1dvSTQ   drive-windows-386.exe   bin    5.2 MB   2015-07-18 16:43:53
0B3X9GlR6EmbnVjIzMDRqck1aekE   drive-osx-x64           bin    6.5 MB   2015-07-18 16:43:50
0B3X9GlR6EmbnbEpXdlhza25zT1U   drive-osx-386           bin    5.2 MB   2015-07-18 16:43:41
0B3X9GlR6Embnb095MGxEYmJhY2c   drive-linux-x64         bin    6.5 MB   2015-07-18 16:43:38
```

#### List largest files
```
$ gdrive list --query "name contains 'gdrive'" --order "quotaBytesUsed desc" -m 3
Id                             Name                     Type   Size     Created
0B3X9GlR6EmbnZXpDRG1xblM2LTg   gdrive-linux-mips64      bin    8.5 MB   2016-02-22 21:07:04
0B3X9GlR6EmbnNW5CTV8xdFkxTjg   gdrive-linux-mips64le    bin    8.5 MB   2016-02-22 21:07:07
0B3X9GlR6EmbnZ1NGS25FdEVlWEk   gdrive-osx-x64           bin    8.3 MB   2016-02-21 20:22:13
```

#### Upload file
```
$ gdrive upload gdrive-osx-x64
Uploading gdrive-osx-x64
Uploaded 0B3X9GlR6EmbnZ1NGS25FdEVlWEk at 3.8 MB/s, total 8.3 MB
```

#### Make directory
```
$ gdrive mkdir gdrive-bin
Directory 0B3X9GlR6EmbnY1RLVTk5VUtOVkk created
```

#### Upload file to directory
```
$ gdrive upload --parent 0B3X9GlR6EmbnY1RLVTk5VUtOVkk gdrive-osx-x64
Uploading gdrive-osx-x64
Uploaded 0B3X9GlR6EmbnNTk0SkV0bm5Hd0E at 2.5 MB/s, total 8.3 MB
```

#### Download file
```
$ gdrive download 0B3X9GlR6EmbnZ1NGS25FdEVlWEk
Downloading gdrive-osx-x64 -> gdrive-osx-x64
Downloaded 0B3X9GlR6EmbnZ1NGS25FdEVlWEk at 8.3 MB/s, total 8.3 MB
```

#### Share a file
```
$ gdrive share 0B3X9GlR6EmbnNTk0SkV0bm5Hd0E
Granted reader permission to anyone
```

#### Pipe content directly to google drive
```
$ echo "Hello World" | gdrive upload - hello.txt
Uploading hello.txt
Uploaded 0B3X9GlR6EmbnaXVrOUpIcWlUS0E at 8.0 B/s, total 12.0 B
```

#### Print file to stdout
```
$ gdrive download --stdout 0B3X9GlR6EmbnaXVrOUpIcWlUS0E
Hello World
```

#### Get file info
```
$ gdrive info 0B3X9GlR6EmbnNTk0SkV0bm5Hd0E
Id: 0B3X9GlR6EmbnNTk0SkV0bm5Hd0E
Name: gdrive-osx-x64
Path: gdrive-bin/gdrive-osx-x64
Mime: application/octet-stream
Size: 8.3 MB
Created: 2016-02-21 20:47:04
Modified: 2016-02-21 20:47:04
Md5sum: b607f29231a3b2d16098c4212516470f
Shared: True
Parents: 0B3X9GlR6EmbnY1RLVTk5VUtOVkk
ViewUrl: https://drive.google.com/file/d/0B3X9GlR6EmbnNTk0SkV0bm5Hd0E/view?usp=drivesdk
DownloadUrl: https://docs.google.com/uc?id=0B3X9GlR6EmbnNTk0SkV0bm5Hd0E&export=download
```

#### Update file (create new revision)
```
$ gdrive update 0B3X9GlR6EmbnNTk0SkV0bm5Hd0E gdrive-osx-x64
Uploading gdrive-osx-x64
Updated 0B3X9GlR6EmbnNTk0SkV0bm5Hd0E at 2.0 MB/s, total 8.3 MB
```

#### List file revisions
```
$ gdrive revision list 0B3X9GlR6EmbnNTk0SkV0bm5Hd0E
Id                                                    Name             Size     Modified              KeepForever
0B3X9GlR6EmbnOFlHSTZQNWJWMGN2ckZucC9VaEUwczV1cUNrPQ   gdrive-osx-x64   8.3 MB   2016-02-21 20:47:04   False
0B3X9GlR6EmbndVEwMlZCUldGWUlPb2lTS25rOFo1L2t6c2ZVPQ   gdrive-osx-x64   8.3 MB   2016-02-21 21:12:09   False
```

#### Download revision
```
$ gdrive revision download 0B3X9GlR6EmbnNTk0SkV0bm5Hd0E 0B3X9GlR6EmbnOFlHSTZQNWJWMGN2ckZucC9VaEUwczV1cUNrPQ
Downloading gdrive-osx-x64 -> gdrive-osx-x64
Download complete, rate: 8.3 MB/s, total size: 8.3 MB
```

#### Export google doc as docx
```
$ gdrive export --mime application/vnd.openxmlformats-officedocument.wordprocessingml.document 1Kt5A8X7X2RQrEi5t6Y9W1LayRc4hyrFiG63y2dIJEvk
Exported 'foo.docx' with mime type: 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
```

#### Import csv as google spreadsheet
```
$ gdrive import foo.csv
Imported 1mTl3DjIvap4tpTX_oMkDcbDT8ShtiGJRlozTfkXpeko with mime type: 'application/vnd.google-apps.spreadsheet'
```

#### Syncing directory to drive
```
# Create directory on drive
$ gdrive mkdir drive-bin
Directory 0B3X9GlR6EmbnOEd6cEh6bU9XZWM created

# Sync to drive
$ gdrive sync upload _release/bin 0B3X9GlR6EmbnOEd6cEh6bU9XZWM
Starting sync...
Collecting local and remote file information...
Found 32 local files and 0 remote files

6 remote directories are missing
[0001/0006] Creating directory drive-bin/bsd
[0002/0006] Creating directory drive-bin/linux
[0003/0006] Creating directory drive-bin/osx
[0004/0006] Creating directory drive-bin/plan9
[0005/0006] Creating directory drive-bin/solaris
[0006/0006] Creating directory drive-bin/windows

26 remote files are missing
[0001/0026] Uploading bsd/gdrive-dragonfly-x64 -> drive-bin/bsd/gdrive-dragonfly-x64
[0002/0026] Uploading bsd/gdrive-freebsd-386 -> drive-bin/bsd/gdrive-freebsd-386
[0003/0026] Uploading bsd/gdrive-freebsd-arm -> drive-bin/bsd/gdrive-freebsd-arm
[0004/0026] Uploading bsd/gdrive-freebsd-x64 -> drive-bin/bsd/gdrive-freebsd-x64
[0005/0026] Uploading bsd/gdrive-netbsd-386 -> drive-bin/bsd/gdrive-netbsd-386
[0006/0026] Uploading bsd/gdrive-netbsd-arm -> drive-bin/bsd/gdrive-netbsd-arm
[0007/0026] Uploading bsd/gdrive-netbsd-x64 -> drive-bin/bsd/gdrive-netbsd-x64
[0008/0026] Uploading bsd/gdrive-openbsd-386 -> drive-bin/bsd/gdrive-openbsd-386
[0009/0026] Uploading bsd/gdrive-openbsd-arm -> drive-bin/bsd/gdrive-openbsd-arm
[0010/0026] Uploading bsd/gdrive-openbsd-x64 -> drive-bin/bsd/gdrive-openbsd-x64
[0011/0026] Uploading linux/gdrive-linux-386 -> drive-bin/linux/gdrive-linux-386
[0012/0026] Uploading linux/gdrive-linux-arm -> drive-bin/linux/gdrive-linux-arm
[0013/0026] Uploading linux/gdrive-linux-arm64 -> drive-bin/linux/gdrive-linux-arm64
[0014/0026] Uploading linux/gdrive-linux-mips64 -> drive-bin/linux/gdrive-linux-mips64
[0015/0026] Uploading linux/gdrive-linux-mips64le -> drive-bin/linux/gdrive-linux-mips64le
[0016/0026] Uploading linux/gdrive-linux-ppc64 -> drive-bin/linux/gdrive-linux-ppc64
[0017/0026] Uploading linux/gdrive-linux-ppc64le -> drive-bin/linux/gdrive-linux-ppc64le
[0018/0026] Uploading linux/gdrive-linux-x64 -> drive-bin/linux/gdrive-linux-x64
[0019/0026] Uploading osx/gdrive-osx-386 -> drive-bin/osx/gdrive-osx-386
[0020/0026] Uploading osx/gdrive-osx-arm -> drive-bin/osx/gdrive-osx-arm
[0021/0026] Uploading osx/gdrive-osx-x64 -> drive-bin/osx/gdrive-osx-x64
[0022/0026] Uploading plan9/gdrive-plan9-386 -> drive-bin/plan9/gdrive-plan9-386
[0023/0026] Uploading plan9/gdrive-plan9-x64 -> drive-bin/plan9/gdrive-plan9-x64
[0024/0026] Uploading solaris/gdrive-solaris-x64 -> drive-bin/solaris/gdrive-solaris-x64
[0025/0026] Uploading windows/gdrive-windows-386.exe -> drive-bin/windows/gdrive-windows-386.exe
[0026/0026] Uploading windows/gdrive-windows-x64.exe -> drive-bin/windows/gdrive-windows-x64.exe
Sync finished in 1m18.891946279s

# Add new local file
$ echo "google drive binaries" > _release/bin/readme.txt

# Sync again
$ gdrive sync upload _release/bin 0B3X9GlR6EmbnOEd6cEh6bU9XZWM
Starting sync...
Collecting local and remote file information...
Found 33 local files and 32 remote files

1 remote files are missing
[0001/0001] Uploading readme.txt -> drive-bin/readme.txt
Sync finished in 2.201339535s

# Modify local file
$ echo "for all platforms" >> _release/bin/readme.txt

# Sync again
$ gdrive sync upload _release/bin 0B3X9GlR6EmbnOEd6cEh6bU9XZWM
Starting sync...
Collecting local and remote file information...
Found 33 local files and 33 remote files

1 local files has changed
[0001/0001] Updating readme.txt -> drive-bin/readme.txt
Sync finished in 1.890244258s
```

#### List content of sync directory
```
$ gdrive sync content 0B3X9GlR6EmbnOEd6cEh6bU9XZWM
Id                             Path                             Type   Size     Modified
0B3X9GlR6EmbnMldxMFV1UGVMTlE   bsd                              dir             2016-02-21 22:54:00
0B3X9GlR6EmbnM05sQ3hVUnJnOXc   bsd/gdrive-dragonfly-x64         bin    7.8 MB   2016-02-21 22:54:14
0B3X9GlR6EmbnVy1KXzA4dlU5RVE   bsd/gdrive-freebsd-386           bin    6.1 MB   2016-02-21 22:54:18
0B3X9GlR6Embnb29QQkFtSlRiZnc   bsd/gdrive-freebsd-arm           bin    6.1 MB   2016-02-21 22:54:20
0B3X9GlR6EmbnMkFQYVpSaHhHTXM   bsd/gdrive-freebsd-x64           bin    7.8 MB   2016-02-21 22:54:23
0B3X9GlR6EmbnVmJRMl9hUDloVU0   bsd/gdrive-netbsd-386            bin    6.1 MB   2016-02-21 22:54:25
0B3X9GlR6EmbnLVlTZWpxOEF4Q2s   bsd/gdrive-netbsd-arm            bin    6.1 MB   2016-02-21 22:54:28
0B3X9GlR6EmbnOENUZmh3anJmNG8   bsd/gdrive-netbsd-x64            bin    7.8 MB   2016-02-21 22:54:30
0B3X9GlR6EmbnWTRoQ2ZVQXRfQlU   bsd/gdrive-openbsd-386           bin    6.1 MB   2016-02-21 22:54:32
0B3X9GlR6EmbncEtlN3ZuQ0VUWms   bsd/gdrive-openbsd-arm           bin    6.1 MB   2016-02-21 22:54:35
0B3X9GlR6EmbnMlFLY1ptNEFyZWc   bsd/gdrive-openbsd-x64           bin    7.8 MB   2016-02-21 22:54:38
0B3X9GlR6EmbncGtSajQyNzloVEE   linux                            dir             2016-02-21 22:54:01
0B3X9GlR6EmbnMWVudkJmb1NZdmM   linux/gdrive-linux-386           bin    6.1 MB   2016-02-21 22:54:40
0B3X9GlR6Embnbnpla1R2VHV5T2M   linux/gdrive-linux-arm           bin    6.1 MB   2016-02-21 22:54:42
0B3X9GlR6EmbnM0s2cU1YWkNJSjA   linux/gdrive-linux-arm64         bin    7.7 MB   2016-02-21 22:54:45
0B3X9GlR6EmbnNU9NNi1TdDc4S2c   linux/gdrive-linux-mips64        bin    8.5 MB   2016-02-21 22:54:47
0B3X9GlR6EmbnSmdQNjRKZ2dWV1U   linux/gdrive-linux-mips64le      bin    8.5 MB   2016-02-21 22:54:50
0B3X9GlR6EmbnS0g0OVgxMHY5Z3c   linux/gdrive-linux-ppc64         bin    7.8 MB   2016-02-21 22:54:52
0B3X9GlR6EmbneVp6ZXRpR3FhWlU   linux/gdrive-linux-ppc64le       bin    7.8 MB   2016-02-21 22:54:54
0B3X9GlR6EmbnczdJT195dFVxdU0   linux/gdrive-linux-x64           bin    7.8 MB   2016-02-21 22:54:57
0B3X9GlR6EmbnTXZXeDRnSDdVS1E   osx                              dir             2016-02-21 22:54:02
0B3X9GlR6EmbnWnRheXJNR0pUMU0   osx/gdrive-osx-386               bin    6.6 MB   2016-02-21 22:54:59
0B3X9GlR6EmbnRzNqMWFXdDR1Rms   osx/gdrive-osx-arm               bin    6.6 MB   2016-02-21 22:55:01
0B3X9GlR6EmbnaDlVWTZDd0JIeEU   osx/gdrive-osx-x64               bin    8.3 MB   2016-02-21 22:55:04
0B3X9GlR6EmbnWW84UFBvbHlURXM   plan9                            dir             2016-02-21 22:54:02
0B3X9GlR6EmbnTmc0a2RNdDZDRUU   plan9/gdrive-plan9-386           bin    5.8 MB   2016-02-21 22:55:07
0B3X9GlR6EmbnT1pYZ2p4Sk9FVFk   plan9/gdrive-plan9-x64           bin    7.4 MB   2016-02-21 22:55:10
0B3X9GlR6EmbnbnZnXzlYVHoxdk0   readme.txt                       bin    40.0 B   2016-02-21 22:59:56
0B3X9GlR6EmbnSWF1QUlta3RnaGc   solaris                          dir             2016-02-21 22:54:03
0B3X9GlR6EmbnaWFOV0YxSGs5Znc   solaris/gdrive-solaris-x64       bin    7.7 MB   2016-02-21 22:55:13
0B3X9GlR6EmbnNE5ySkEzbWQ4Qms   windows                          dir             2016-02-21 22:54:03
0B3X9GlR6EmbnX1RIT2w1TWZYWFU   windows/gdrive-windows-386.exe   bin    6.1 MB   2016-02-21 22:55:15
0B3X9GlR6EmbndmVMU05POGRPS3c   windows/gdrive-windows-x64.exe   bin    7.8 MB   2016-02-21 22:55:18
```
