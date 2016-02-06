package drive

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/gyuho/goraph/graph"
    "github.com/sabhiram/go-git-ignore"
    "golang.org/x/net/context"
    "google.golang.org/api/drive/v3"
    "google.golang.org/api/googleapi"
)

const DefaultIgnoreFile = ".gdriveignore"

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

func prepareLocalFiles(root string) ([]*localFile, error) {
    var files []*localFile

    // Get absolute root path
    absRootPath, err := filepath.Abs(root)
    if err != nil {
        return nil, err
    }

    // Prepare ignorer
    shouldIgnore, err := prepareIgnorer(filepath.Join(absRootPath, DefaultIgnoreFile))
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

        // Get relative path from root
        relPath, err := filepath.Rel(absRootPath, absPath)
        if err != nil {
            return err
        }

        // Skip file if it is ignored by ignore file
        if shouldIgnore(relPath) {
            return nil
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

func (self *Drive) listAllFiles(q string, fields []googleapi.Field) ([]*drive.File, error) {
    var files []*drive.File

    err := self.service.Files.List().Q(q).Fields(fields...).PageSize(1000).Pages(context.TODO(), func(fl *drive.FileList) error {
        files = append(files, fl.Files...)
        return nil
    })

    return files, err
}

func (self *Drive) prepareRemoteFiles(rootDir *drive.File) ([]*remoteFile, error) {
    // Find all files which has rootDir as root
    query := fmt.Sprintf("appProperties has {key='syncRootId' and value='%s'}", rootDir.Id)
    fields := []googleapi.Field{"nextPageToken", "files(id,name,parents,md5Checksum,mimeType)"}
    files, err := self.listAllFiles(query, fields)
    if err != nil {
        return nil, fmt.Errorf("Failed listing files: %s", err)
    }

    if err := checkFiles(files); err != nil {
        return nil, err
    }

    relPaths, err := prepareRemoteRelPaths(rootDir, files)
    if err != nil {
        return nil, err
    }

    var remoteFiles []*remoteFile
    for _, f := range files {
        relPath, ok := relPaths[f.Id]
        if !ok {
            return nil, fmt.Errorf("File %s does not have a valid parent", f.Id)
        }
        remoteFiles = append(remoteFiles, &remoteFile{
            relPath: relPath,
            file: f,
        })
    }

    return remoteFiles, nil
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

func checkFiles(files []*drive.File) error {
    uniq := map[string]string{}

    for _, f := range files {
        // Ensure all files have exactly one parent
        if len(f.Parents) != 1 {
            return fmt.Errorf("File %s does not have exacly one parent", f.Id)
        }

        // Ensure that there are no duplicate files
        uniqKey := f.Name + f.Parents[0]
        if dupeId, isDupe := uniq[uniqKey]; isDupe {
            return fmt.Errorf("Found name collision between %s and %s", f.Id, dupeId)
        }
        uniq[uniqKey] = f.Id
    }

    return nil
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

type syncFiles struct {
    root *remoteFile
    local []*localFile
    remote []*remoteFile
}

func (self *syncFiles) filterMissingRemoteDirs() []*localFile {
    var files []*localFile

    for _, lf := range self.local {
        if lf.info.IsDir() && !self.existsRemote(lf) {
            files = append(files, lf)
        }
    }

    return files
}

func (self *syncFiles) filterMissingLocalDirs() []*remoteFile {
    var files []*remoteFile

    for _, rf := range self.remote {
        if isDir(rf.file) && !self.existsLocal(rf) {
            files = append(files, rf)
        }
    }

    return files
}

func (self *syncFiles) filterMissingRemoteFiles() []*localFile {
    var files []*localFile

    for _, lf := range self.local {
        if !lf.info.IsDir() && !self.existsRemote(lf) {
            files = append(files, lf)
        }
    }

    return files
}

func (self *syncFiles) filterMissingLocalFiles() []*remoteFile {
    var files []*remoteFile

    for _, rf := range self.remote {
        if !isDir(rf.file) && !self.existsLocal(rf) {
            files = append(files, rf)
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

func (self *syncFiles) filterChangedRemoteFiles() []*changedFile {
    var files []*changedFile

    for _, rf := range self.remote {
        // Skip directories
        if isDir(rf.file) {
            continue
        }

        // Skip local files that don't exist
        lf, found := self.findLocalByPath(rf.relPath)
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

func (self *syncFiles) filterExtraneousRemoteFiles() []*remoteFile {
    var files []*remoteFile

    for _, rf := range self.remote {
        if !self.existsLocal(rf) {
            files = append(files, rf)
        }
    }

    return files
}

func (self *syncFiles) filterExtraneousLocalFiles() []*localFile {
    var files []*localFile

    for _, lf := range self.local {
        if !self.existsRemote(lf) {
            files = append(files, lf)
        }
    }

    return files
}

func (self *syncFiles) existsRemote(lf *localFile) bool {
    _, found := self.findRemoteByPath(lf.relPath)
    return found
}

func (self *syncFiles) existsLocal(rf *remoteFile) bool {
    _, found := self.findLocalByPath(rf.relPath)
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

func (self *syncFiles) findLocalByPath(relPath string) (*localFile, bool) {
    for _, lf := range self.local {
        if relPath == lf.relPath {
            return lf, true
        }
    }

    return nil, false
}

type byLocalPathLength []*localFile

func (self byLocalPathLength) Len() int {
    return len(self)
}

func (self byLocalPathLength) Swap(i, j int) {
    self[i], self[j] = self[j], self[i]
}

func (self byLocalPathLength) Less(i, j int) bool {
    return pathLength(self[i].relPath) < pathLength(self[j].relPath)
}

type byRemotePathLength []*remoteFile

func (self byRemotePathLength) Len() int {
    return len(self)
}

func (self byRemotePathLength) Swap(i, j int) {
    self[i], self[j] = self[j], self[i]
}

func (self byRemotePathLength) Less(i, j int) bool {
    return pathLength(self[i].relPath) < pathLength(self[j].relPath)
}

type ignoreFunc func(string) bool

func prepareIgnorer(path string) (ignoreFunc, error) {
    acceptAll := func(string) bool {
        return false
    }

    if !fileExists(path) {
        return acceptAll, nil
    }

    ignorer, err := ignore.CompileIgnoreFile(path)
    if err != nil {
        return acceptAll, fmt.Errorf("Failed to prepare ignorer: %s", err)
    }

    return ignorer.MatchesPath, nil
}
