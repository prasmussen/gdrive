package drive

import (
    "fmt"
    "io"
    "os"
    "time"
    "sort"
    "path/filepath"
    "google.golang.org/api/googleapi"
    "google.golang.org/api/drive/v3"
)

type UploadSyncArgs struct {
    Out io.Writer
    Progress io.Writer
    Path string
    RootId string
    DryRun bool
    DeleteExtraneous bool
    ChunkSize int64
    Comparer FileComparer
}

func (self *Drive) UploadSync(args UploadSyncArgs) error {
    if args.ChunkSize > intMax() - 1 {
        return fmt.Errorf("Chunk size is to big, max chunk size for this computer is %d", intMax() - 1)
    }

    fmt.Fprintln(args.Out, "Starting sync...")
    started := time.Now()

    // Create root directory if it does not exist
    rootDir, err := self.prepareSyncRoot(args)
    if err != nil {
        return err
    }

    fmt.Fprintln(args.Out, "Collecting local and remote file information...")
    files, err := self.prepareSyncFiles(args.Path, rootDir, args.Comparer)
    if err != nil {
        return err
    }

    fmt.Fprintf(args.Out, "Found %d local files and %d remote files\n", len(files.local), len(files.remote))

    // Create missing directories
    files, err = self.createMissingRemoteDirs(files, args)
    if err != nil {
        return err
    }

    // Upload missing files
    err = self.uploadMissingFiles(files, args)
    if err != nil {
        return err
    }

    // Update modified files
    err = self.updateChangedFiles(files, args)
    if err != nil {
        return err
    }

    // Delete extraneous files on drive
    if args.DeleteExtraneous {
        err = self.deleteExtraneousRemoteFiles(files, args)
        if err != nil {
            return err
        }
    }
    fmt.Fprintf(args.Out, "Sync finished in %s\n", time.Since(started))

    return nil
}

func (self *Drive) prepareSyncRoot(args UploadSyncArgs) (*drive.File, error) {
    fields := []googleapi.Field{"id", "name", "mimeType", "appProperties"}
    f, err := self.service.Files.Get(args.RootId).Fields(fields...).Do()
    if err != nil {
        return nil, fmt.Errorf("Failed to find root dir: %s", err)
    }

    // Ensure file is a directory
    if !isDir(f) {
        return nil, fmt.Errorf("Provided root id is not a directory")
    }

    // Return directory if syncRoot property is already set
    if _, ok := f.AppProperties["isSyncRoot"]; ok {
        return f, nil
    }

    // This is the first time this directory have been used for sync
    // Check if the directory is empty
    isEmpty, err := self.dirIsEmpty(f.Id)
    if err != nil {
        return nil, fmt.Errorf("Failed to check if root dir is empty: %s", err)
    }

    // Ensure that the directory is empty
    if !isEmpty {
        return nil, fmt.Errorf("Root directoy is not empty, the initial sync requires an empty directory")
    }

    // Update directory with syncRoot property
    dstFile := &drive.File{
        AppProperties: map[string]string{"isSyncRoot": "true"},
    }

    f, err = self.service.Files.Update(f.Id, dstFile).Fields(fields...).Do()
    if err != nil {
        return nil, fmt.Errorf("Failed to update root directory: %s", err)
    }

    return f, nil
}

func (self *Drive) createMissingRemoteDirs(files *syncFiles, args UploadSyncArgs) (*syncFiles, error) {
    missingDirs := files.filterMissingRemoteDirs()
    missingCount := len(missingDirs)

    if missingCount > 0 {
        fmt.Fprintf(args.Out, "\n%d remote directories are missing\n", missingCount)
    }

    // Sort directories so that the dirs with the shortest path comes first
    sort.Sort(byLocalPathLength(missingDirs))

    for i, lf := range missingDirs {
        parentPath := parentFilePath(lf.relPath)
        parent, ok := files.findRemoteByPath(parentPath)
        if !ok {
            return nil, fmt.Errorf("Could not find remote directory with path '%s'", parentPath)
        }

        dstFile := &drive.File{
            Name: lf.info.Name(),
            MimeType: DirectoryMimeType,
            Parents: []string{parent.file.Id},
            AppProperties: map[string]string{"syncRootId": args.RootId},
        }

        fmt.Fprintf(args.Out, "[%04d/%04d] Creating directory: %s\n", i + 1, missingCount, filepath.Join(files.root.file.Name, lf.relPath))

        if args.DryRun {
            files.remote = append(files.remote, &RemoteFile{
                relPath: lf.relPath,
                file: dstFile,
            })
        } else {
            f, err := self.service.Files.Create(dstFile).Do()
            if err != nil {
                return nil, fmt.Errorf("Failed to create directory: %s", err)
            }

            files.remote = append(files.remote, &RemoteFile{
                relPath: lf.relPath,
                file: f,
            })
        }
    }

    return files, nil
}

func (self *Drive) uploadMissingFiles(files *syncFiles, args UploadSyncArgs) error {
    missingFiles := files.filterMissingRemoteFiles()
    missingCount := len(missingFiles)

    if missingCount > 0 {
        fmt.Fprintf(args.Out, "\n%d remote files are missing\n", missingCount)
    }

    for i, lf := range missingFiles {
        parentPath := parentFilePath(lf.relPath)
        parent, ok := files.findRemoteByPath(parentPath)
        if !ok {
            return fmt.Errorf("Could not find remote directory with path '%s'", parentPath)
        }

        fmt.Fprintf(args.Out, "[%04d/%04d] Uploading %s -> %s\n", i + 1, missingCount, lf.absPath, filepath.Join(files.root.file.Name, lf.relPath))

        if args.DryRun {
            continue
        }

        err := self.uploadMissingFile(parent.file.Id, lf, args)
        if err != nil {
            return err
        }
    }

    return nil
}

func (self *Drive) updateChangedFiles(files *syncFiles, args UploadSyncArgs) error {
    changedFiles := files.filterChangedLocalFiles()
    changedCount := len(changedFiles)

    if changedCount > 0 {
        fmt.Fprintf(args.Out, "\n%d local files has changed\n", changedCount)
    }

    for i, cf := range changedFiles {
        fmt.Fprintf(args.Out, "[%04d/%04d] Updating %s -> %s\n", i + 1, changedCount, cf.local.absPath, filepath.Join(files.root.file.Name, cf.local.relPath))

        if args.DryRun {
            continue
        }

        err := self.updateChangedFile(cf, args)
        if err != nil {
            return err
        }
    }

    return nil
}

func (self *Drive) deleteExtraneousRemoteFiles(files *syncFiles, args UploadSyncArgs) error {
    extraneousFiles := files.filterExtraneousRemoteFiles()
    extraneousCount := len(extraneousFiles)

    if extraneousCount > 0 {
        fmt.Fprintf(args.Out, "\n%d remote files are extraneous\n", extraneousCount)
    }

    // Sort files so that the files with the longest path comes first
    sort.Sort(sort.Reverse(byRemotePathLength(extraneousFiles)))

    for i, rf := range extraneousFiles {
        fmt.Fprintf(args.Out, "[%04d/%04d] Deleting %s\n", i + 1, extraneousCount, filepath.Join(files.root.file.Name, rf.relPath))

        if args.DryRun {
            continue
        }

        err := self.deleteRemoteFile(rf, args)
        if err != nil {
            return err
        }
    }

    return nil
}

func (self *Drive) uploadMissingFile(parentId string, lf *LocalFile, args UploadSyncArgs) error {
    srcFile, err := os.Open(lf.absPath)
    if err != nil {
        return fmt.Errorf("Failed to open file: %s", err)
    }

    // Close file on function exit
    defer srcFile.Close()

    // Instantiate drive file
    dstFile := &drive.File{
        Name: lf.info.Name(),
        Parents: []string{parentId},
        AppProperties: map[string]string{"syncRootId": args.RootId},
    }

    // Chunk size option
    chunkSize := googleapi.ChunkSize(int(args.ChunkSize))

    // Wrap file in progress reader
    srcReader := getProgressReader(srcFile, args.Progress, lf.info.Size())

    _, err = self.service.Files.Create(dstFile).Fields("id", "name", "size", "md5Checksum").Media(srcReader, chunkSize).Do()
    if err != nil {
        return fmt.Errorf("Failed to upload file: %s", err)
    }

    return nil
}

func (self *Drive) updateChangedFile(cf *changedFile, args UploadSyncArgs) error {
    srcFile, err := os.Open(cf.local.absPath)
    if err != nil {
        return fmt.Errorf("Failed to open file: %s", err)
    }

    // Close file on function exit
    defer srcFile.Close()

    // Instantiate drive file
    dstFile := &drive.File{}

    // Chunk size option
    chunkSize := googleapi.ChunkSize(int(args.ChunkSize))

    // Wrap file in progress reader
    srcReader := getProgressReader(srcFile, args.Progress, cf.local.info.Size())

    _, err = self.service.Files.Update(cf.remote.file.Id, dstFile).Media(srcReader, chunkSize).Do()
    if err != nil {
        return fmt.Errorf("Failed to update file: %s", err)
    }

    return nil
}

func (self *Drive) deleteRemoteFile(rf *RemoteFile, args UploadSyncArgs) error {
    err := self.service.Files.Delete(rf.file.Id).Do()
    if err != nil {
        return fmt.Errorf("Failed to delete file: %s", err)
    }

    return nil
}

func (self *Drive) dirIsEmpty(id string) (bool, error) {
    query := fmt.Sprintf("'%s' in parents", id)
    fileList, err := self.service.Files.List().Q(query).Do()
    if err != nil {
        return false, fmt.Errorf("Empty dir check failed: ", err)
    }

    return len(fileList.Files) == 0, nil
}
