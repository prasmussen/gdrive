package drive

import (
    "fmt"
    "io"
    "os"
    "sort"
    "time"
    "path/filepath"
    "google.golang.org/api/googleapi"
    "google.golang.org/api/drive/v3"
)

type DownloadSyncArgs struct {
    Out io.Writer
    Progress io.Writer
    RootId string
    Path string
    DryRun bool
    DeleteExtraneous bool
}

func (self *Drive) DownloadSync(args DownloadSyncArgs) error {
    fmt.Fprintln(args.Out, "Starting sync...")
    started := time.Now()

    // Get remote root dir
    rootDir, err := self.getSyncRoot(args.RootId)
    if err != nil {
        return err
    }

    fmt.Fprintln(args.Out, "Collecting local and remote file information...")
    files, err := self.prepareSyncFiles(args.Path, rootDir)
    if err != nil {
        return err
    }

    fmt.Fprintf(args.Out, "Found %d local files and %d remote files\n", len(files.local), len(files.remote))

    // Create missing directories
    err = self.createMissingLocalDirs(files, args)
    if err != nil {
        return err
    }

    // Download missing files
    err = self.downloadMissingFiles(files, args)
    if err != nil {
        return err
    }

    // Download files that has changed
    err = self.downloadChangedFiles(files, args)
    if err != nil {
        return err
    }

    // Delete extraneous local files
    if args.DeleteExtraneous {
        err = self.deleteExtraneousLocalFiles(files, args)
        if err != nil {
            return err
        }
    }
    fmt.Fprintf(args.Out, "Sync finished in %s\n", time.Since(started))

    return nil
}

func (self *Drive) getSyncRoot(rootId string) (*drive.File, error) {
    fields := []googleapi.Field{"id", "name", "mimeType", "appProperties"}
    f, err := self.service.Files.Get(rootId).Fields(fields...).Do()
    if err != nil {
        return nil, fmt.Errorf("Failed to find root dir: %s", err)
    }

    // Ensure file is a directory
    if !isDir(f) {
        return nil, fmt.Errorf("Provided root id is not a directory")
    }

    // Ensure directory is a proper syncRoot
    if _, ok := f.AppProperties["isSyncRoot"]; !ok {
        return nil, fmt.Errorf("Provided id is not a sync root directory")
    }

    return f, nil
}

func (self *Drive) createMissingLocalDirs(files *syncFiles, args DownloadSyncArgs) error {
    missingDirs := files.filterMissingLocalDirs()
    missingCount := len(missingDirs)

    if missingCount > 0 {
        fmt.Fprintf(args.Out, "\n%d local directories are missing\n", missingCount)
    }

    // Sort directories so that the dirs with the shortest path comes first
    sort.Sort(byRemotePathLength(missingDirs))

    for i, rf := range missingDirs {
        path, err := filepath.Abs(filepath.Join(args.Path, rf.relPath))
        if err != nil {
            return fmt.Errorf("Failed to determine local absolute path: %s", err)
        }
        fmt.Fprintf(args.Out, "[%04d/%04d] Creating directory: %s\n", i + 1, missingCount, path)

        if args.DryRun {
            continue
        }

        mkdir(path)
    }

    return nil
}

func (self *Drive) downloadMissingFiles(files *syncFiles, args DownloadSyncArgs) error {
    missingFiles := files.filterMissingLocalFiles()
    missingCount := len(missingFiles)

    if missingCount > 0 {
        fmt.Fprintf(args.Out, "\n%d local files are missing\n", missingCount)
    }

    for i, rf := range missingFiles {
        remotePath := filepath.Join(files.root.file.Name, rf.relPath)
        localPath, err := filepath.Abs(filepath.Join(args.Path, rf.relPath))
        if err != nil {
            return fmt.Errorf("Failed to determine local absolute path: %s", err)
        }
        fmt.Fprintf(args.Out, "[%04d/%04d] Downloading %s -> %s\n", i + 1, missingCount, remotePath, localPath)

        if args.DryRun {
            continue
        }

        err = self.downloadRemoteFile(rf.file.Id, localPath, args)
        if err != nil {
            return err
        }
    }

    return nil
}

func (self *Drive) downloadChangedFiles(files *syncFiles, args DownloadSyncArgs) error {
    changedFiles := files.filterChangedRemoteFiles()
    changedCount := len(changedFiles)

    if changedCount > 0 {
        fmt.Fprintf(args.Out, "\n%d remote files has changed\n", changedCount)
    }

    for i, cf := range changedFiles {
        remotePath := filepath.Join(files.root.file.Name, cf.remote.relPath)
        localPath, err := filepath.Abs(filepath.Join(args.Path, cf.remote.relPath))
        if err != nil {
            return fmt.Errorf("Failed to determine local absolute path: %s", err)
        }
        fmt.Fprintf(args.Out, "[%04d/%04d] Downloading %s -> %s\n", i + 1, changedCount, remotePath, localPath)

        if args.DryRun {
            continue
        }

        err = self.downloadRemoteFile(cf.remote.file.Id, localPath, args)
        if err != nil {
            return err
        }
    }

    return nil
}

func (self *Drive) downloadRemoteFile(id, fpath string, args DownloadSyncArgs) error {
    res, err := self.service.Files.Get(id).Download()
    if err != nil {
        return fmt.Errorf("Failed to download file: %s", err)
    }

    // Close body on function exit
    defer res.Body.Close()

    // Wrap response body in progress reader
    srcReader := getProgressReader(res.Body, args.Progress, res.ContentLength)

    // Ensure any parent directories exists
    if err = mkdir(fpath); err != nil {
        return err
    }

    // Create new file
    outFile, err := os.Create(fpath)
    if err != nil {
        return fmt.Errorf("Unable to create local file: %s", err)
    }

    // Close file on function exit
    defer outFile.Close()

    // Save file to disk
    _, err = io.Copy(outFile, srcReader)
    if err != nil {
        return fmt.Errorf("Download was interrupted: %s", err)
    }

    return nil
}

func (self *Drive) deleteExtraneousLocalFiles(files *syncFiles, args DownloadSyncArgs) error {
    extraneousFiles := files.filterExtraneousLocalFiles()
    extraneousCount := len(extraneousFiles)

    if extraneousCount > 0 {
        fmt.Fprintf(args.Out, "\n%d local files are extraneous\n", extraneousCount)
    }

    // Sort files so that the files with the longest path comes first
    sort.Sort(sort.Reverse(byLocalPathLength(extraneousFiles)))

    for i, lf := range extraneousFiles {
        fmt.Fprintf(args.Out, "[%04d/%04d] Deleting %s\n", i + 1, extraneousCount, lf.absPath)

        if args.DryRun {
            continue
        }

        err := os.Remove(lf.absPath)
        if err != nil {
            return fmt.Errorf("Failed to delete local file: %s", err)
        }
    }

    return nil
}
