package drive

import (
    "fmt"
    "io"
    "os"
    "time"
    "sort"
    "path/filepath"
    "github.com/gyuho/goraph/graph"
    "google.golang.org/api/googleapi"
    "google.golang.org/api/drive/v3"
)

type UploadSyncArgs struct {
    Out io.Writer
    Progress io.Writer
    Path string
    Parent string
    DeleteRemote bool
    ChunkSize int64
}

func (self *Drive) UploadSync(args UploadSyncArgs) error {
    if args.ChunkSize > intMax() - 1 {
        return fmt.Errorf("Chunk size is to big, max chunk size for this computer is %d", intMax() - 1)
    }

    fmt.Fprintln(args.Out, "Starting sync...")
    started := time.Now()

    // Create root directory if it does not exist
    rootDir, createdRoot, err := self.getOrCreateSyncRootDir(args)
    if err != nil {
        return err
    }

    if createdRoot {
        fmt.Fprintln(args.Out, "No existing files found on drive, starting from scratch")
    }

    fmt.Fprintln(args.Out, "Collecting local and remote file information...")
    files, err := self.prepareSyncFiles(args.Path, rootDir)
    if err != nil {
        return err
    }

    fmt.Fprintf(args.Out, "Found %d local file(s) and %d remote file(s)\n", len(files.local), len(files.remote))

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

    fmt.Fprintf(args.Out, "Sync finished in %s\n", time.Since(started))

    return nil
}

func (self *Drive) prepareSyncFiles(localPath string, root *drive.File) (*syncFiles, error) {
    localCh := make(chan struct{files []*localFile; err error})
    remoteCh := make(chan struct{files []*remoteFile; err error})

    go func() {
        files, err := prepareLocalFiles(localPath)
        localCh <- struct{files []*localFile; err error}{files, err}
    }()

    go func() {
        files, err := self.prepareRemoteFiles(root)
        remoteCh <- struct{files []*remoteFile; err error}{files, err}
    }()

    local := <-localCh
    if local.err != nil {
        return nil, local.err
    }

    remote := <-remoteCh
    if remote.err != nil {
        return nil, remote.err
    }

    return &syncFiles{
        root: &remoteFile{file: root},
        local: local.files,
        remote: remote.files,
    }, nil
}

func (self *Drive) getOrCreateSyncRootDir(args UploadSyncArgs) (*drive.File, bool, error) {
    // Root dir name
    name := filepath.Base(args.Path)

    // Build root dir query
    query := fmt.Sprintf("name = '%s' and appProperties has {key='isSyncRoot' and value='true'}", name)
    if args.Parent != "" {
        query += fmt.Sprintf(" and '%s' in parents", args.Parent)
    }

    // Find root dir
    fileList, err := self.service.Files.List().Q(query).Fields("files(id,name,mimeType)").Do()
    if err != nil {
        return nil, false, fmt.Errorf("Failed listing files: %s", err)
    }

    // More than one root dir found
    if len(fileList.Files) > 1 {
        return nil, false, fmt.Errorf("More than one root directories found, aborting...")
    }

    // Root dir found, return
    if len(fileList.Files) == 1 {
        return fileList.Files[0], false, nil
    }

    // Root dir not found, create new
    dstFile := &drive.File{
        Name: name,
        MimeType: DirectoryMimeType,
        AppProperties: map[string]string{"isSyncRoot": "true"},
    }

    // Add parent if provided
    if args.Parent != "" {
        dstFile.Parents = []string{args.Parent}
    }

    // Create directory
    f, err := self.service.Files.Create(dstFile).Do()
    if err != nil {
        return nil, false, fmt.Errorf("Failed to create directory: %s", err)
    }

    return f, true, nil
}

func (self *Drive) createMissingRemoteDirs(files *syncFiles, args UploadSyncArgs) (*syncFiles, error) {
    missingDirs := files.filterMissingRemoteDirs()
    missingCount := len(missingDirs)

    if missingCount > 0 {
        fmt.Fprintf(args.Out, "\n%d directories missing on drive\n", missingCount)
    }

    // Sort directories so that the dirs with the shortest path comes first
    sort.Sort(byPathLength(missingDirs))

    for i, lf := range missingDirs {
        parentPath := parentFilePath(lf.relPath)
        parent, ok := files.findRemoteByPath(parentPath)
        if !ok {
            return nil, fmt.Errorf("Could not find remote directory with path '%s', aborting...", parentPath)
        }

        dstFile := &drive.File{
            Name: lf.info.Name(),
            MimeType: DirectoryMimeType,
            Parents: []string{parent.file.Id},
            AppProperties: map[string]string{"syncRootId": files.root.file.Id},
        }

        fmt.Fprintf(args.Out, "[%04d/%04d] Creating directory: %s\n", i + 1, missingCount, filepath.Join(files.root.file.Name, lf.relPath))

        f, err := self.service.Files.Create(dstFile).Do()
        if err != nil {
            return nil, fmt.Errorf("Failed to create directory: %s", err)
        }

        files.remote = append(files.remote, &remoteFile{
            relPath: lf.relPath,
            file: f,
        })
    }

    return files, nil
}

func (self *Drive) uploadMissingFiles(files *syncFiles, args UploadSyncArgs) error {
    missingFiles := files.filterMissingRemoteFiles()
    missingCount := len(missingFiles)

    if missingCount > 0 {
        fmt.Fprintf(args.Out, "\n%d file(s) missing on drive\n", missingCount)
    }

    for i, lf := range missingFiles {
        parentPath := parentFilePath(lf.relPath)
        parent, ok := files.findRemoteByPath(parentPath)
        if !ok {
            return fmt.Errorf("Could not find remote directory with path '%s', aborting...", parentPath)
        }

        newArgs := args
        newArgs.Path = lf.absPath
        newArgs.Parent = parent.file.Id

        fmt.Fprintf(args.Out, "[%04d/%04d] Uploading %s -> %s\n", i + 1, missingCount, lf.absPath, filepath.Join(files.root.file.Name, lf.relPath))
        err := self.uploadMissingFile(files.root.file.Id, lf, newArgs)
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
        fmt.Fprintf(args.Out, "\n%d local file(s) has changed\n", changedCount)
    }

    for i, cf := range changedFiles {
        fmt.Fprintf(args.Out, "[%04d/%04d] Updating %s -> %s\n", i + 1, changedCount, cf.local.absPath, filepath.Join(files.root.file.Name, cf.local.relPath))
        err := self.updateChangedFile(cf, args)
        if err != nil {
            return err
        }
    }

    return nil
}

func (self *Drive) uploadMissingFile(rootId string, lf *localFile, args UploadSyncArgs) error {
    srcFile, err := os.Open(lf.absPath)
    if err != nil {
        return fmt.Errorf("Failed to open file: %s", err)
    }

    // Instantiate drive file
    dstFile := &drive.File{
        Name: lf.info.Name(),
        Parents: []string{args.Parent},
        AppProperties: map[string]string{"syncRootId": rootId},
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

func (self *Drive) prepareRemoteFiles(rootDir *drive.File) ([]*remoteFile, error) {
    // Find all files which has rootDir as root
    query := fmt.Sprintf("appProperties has {key='syncRootId' and value='%s'}", rootDir.Id)
    fileList, err := self.service.Files.List().Q(query).Fields("files(id,name,parents,md5Checksum,mimeType)").Do()
    if err != nil {
        return nil, fmt.Errorf("Failed listing files: %s", err)
    }

    if err := checkFiles(fileList.Files); err != nil {
        return nil, err
    }

    relPaths, err := prepareRemoteRelPaths(rootDir.Id, fileList.Files)
    if err != nil {
        return nil, err
    }

    var remoteFiles []*remoteFile
    for _, f := range fileList.Files {
        relPath, ok := relPaths[f.Id]
        if !ok {
            return nil, fmt.Errorf("File %s does not have a valid parent, aborting...", f.Id)
        }
        remoteFiles = append(remoteFiles, &remoteFile{
            relPath: relPath,
            file: f,
        })
    }

    return remoteFiles, nil
}

func checkFiles(files []*drive.File) error {
    uniq := map[string]string{}

    for _, f := range files {
        // Ensure all files have exactly one parent
        if len(f.Parents) != 1 {
            return fmt.Errorf("File %s does not have exacly one parent, aborting...", f.Id)
        }

        // Ensure that there are no duplicate files
        uniqKey := f.Name + f.Parents[0]
        if dupeId, isDupe := uniq[uniqKey]; isDupe {
            return fmt.Errorf("Found name collision between %s and %s, aborting", f.Id, dupeId)
        }
        uniq[uniqKey] = f.Id
    }

    return nil
}

func prepareRemoteRelPaths(rootId string, files []*drive.File) (map[string]string, error) {
    names := map[string]string{}
    idGraph := graph.NewDefaultGraph()

    for _, f := range files {
        // Store directory name for quick lookup
        names[f.Id] = f.Name

        // Store path between parent and child folder
        idGraph.AddVertex(f.Id)
        idGraph.AddVertex(f.Parents[0])
        idGraph.AddEdge(f.Parents[0], f.Id, 0)
    }

    paths := map[string]string{}

    for _, f := range files {
        // Find path from root to directory
        pathIds, _, err := graph.Dijkstra(idGraph, rootId, f.Id)
        if err != nil {
            return nil, err
        }

        // Convert path ids to path names
        var pathNames []string
        for _, id := range pathIds {
            pathNames = append(pathNames, names[id])
        }

        // Store relative file path from root to directory
        paths[f.Id] = filepath.Join(pathNames...)
    }

    return paths, nil
}

type localFile struct {
    absPath string
    relPath string
    info os.FileInfo
}

type remoteFile struct {
    relPath string
    file *drive.File
}

type changedFile struct {
    local *localFile
    remote *remoteFile
}

func prepareLocalFiles(root string) ([]*localFile, error) {
    var files []*localFile

    // Get absolute root path
    absRootPath, err := filepath.Abs(root)
    if err != nil {
        return nil, err
    }

    err = filepath.Walk(absRootPath, func(absPath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Skip root directory
        if absPath == absRootPath {
            return nil
        }

        relPath, err := filepath.Rel(absRootPath, absPath)
        if err != nil {
            return err
        }

        files = append(files, &localFile{
            absPath: absPath,
            relPath: relPath,
            info: info,
        })

        return nil
    })

    if err != nil {
        return nil, fmt.Errorf("Failed to prepare local files: %s", err)
    }

    return files, err
}

type syncFiles struct {
    root *remoteFile
    local []*localFile
    remote []*remoteFile
}

func (self *syncFiles) filterMissingRemoteDirs() []*localFile {
    var files []*localFile

    for _, f := range self.local {
        if f.info.IsDir() && !self.existsRemote(f) {
            files = append(files, f)
        }
    }

    return files
}

func (self *syncFiles) filterMissingRemoteFiles() []*localFile {
    var files []*localFile

    for _, f := range self.local {
        if !f.info.IsDir() && !self.existsRemote(f) {
            files = append(files, f)
        }
    }

    return files
}

func (self *syncFiles) filterChangedLocalFiles() []*changedFile {
    var files []*changedFile

    for _, lf := range self.local {
        // Skip directories
        if lf.info.IsDir() {
            continue
        }

        // Skip files that don't exist on drive
        rf, found := self.findRemoteByPath(lf.relPath)
        if !found {
            continue
        }

        // Add files where remote md5 sum does not match local
        if rf.file.Md5Checksum != md5sum(lf.absPath) {
            files = append(files, &changedFile{
                local: lf,
                remote: rf,
            })
        }
    }

    return files
}

func (self *syncFiles) existsRemote(lf *localFile) bool {
    _, found := self.findRemoteByPath(lf.relPath)
    return found
}

func (self *syncFiles) findRemoteByPath(relPath string) (*remoteFile, bool) {
    if relPath == "." {
        return self.root, true
    }

    for _, rf := range self.remote {
        if relPath == rf.relPath {
            return rf, true
        }
    }

    return nil, false
}

type byPathLength []*localFile

func (self byPathLength) Len() int {
    return len(self)
}

func (self byPathLength) Swap(i, j int) {
    self[i], self[j] = self[j], self[i]
}

func (self byPathLength) Less(i, j int) bool {
    return pathLength(self[i].relPath) < pathLength(self[j].relPath)
}
