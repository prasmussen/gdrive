package drive

import (
	"bytes"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type DownloadSyncArgs struct {
	Out              io.Writer
	Progress         io.Writer
	RootId           string
	Path             string
	DryRun           bool
	DeleteExtraneous bool
	Timeout          time.Duration
	Resolution       ConflictResolution
	Comparer         FileComparer
}

func (self *Drive) DownloadSync(args DownloadSyncArgs) error {
	fmt.Fprintln(args.Out, "Starting sync...")
	started := time.Now()

	// Get remote root dir
	rootDir, err := self.getSyncRoot(args.RootId)
	if err != nil {
		return err
	}

	fmt.Fprintln(args.Out, "Collecting file information...")
	files, err := self.prepareSyncFiles(args.Path, rootDir, args.Comparer)
	if err != nil {
		return err
	}

	// Find changed files
	changedFiles := files.filterChangedRemoteFiles()

	fmt.Fprintf(args.Out, "Found %d local files and %d remote files\n", len(files.local), len(files.remote))

	// Ensure that we don't overwrite any local changes
	if args.Resolution == NoResolution {
		err = ensureNoLocalModifications(changedFiles)
		if err != nil {
			return fmt.Errorf("Conflict detected!\nThe following files have changed and the local file are newer than it's remote counterpart:\n\n%s\nNo conflict resolution was given, aborting...", err)
		}
	}

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
	err = self.downloadChangedFiles(changedFiles, args)
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
	f, err := self.service.Files.Get(rootId).SupportsTeamDrives(true).Fields(fields...).Do()
	if err != nil {
		return nil, fmt.Errorf("Failed to find root dir: %s", err)
	}

	// Ensure file is a directory
	if !isDir(f) {
		return nil, fmt.Errorf("Provided root id is not a directory")
	}

	// Ensure directory is a proper syncRoot
	if _, ok := f.AppProperties["syncRoot"]; !ok {
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
		absPath, err := filepath.Abs(filepath.Join(args.Path, rf.relPath))
		if err != nil {
			return fmt.Errorf("Failed to determine local absolute path: %s", err)
		}
		fmt.Fprintf(args.Out, "[%04d/%04d] Creating directory %s\n", i+1, missingCount, filepath.Join(filepath.Base(args.Path), rf.relPath))

		if args.DryRun {
			continue
		}

		os.MkdirAll(absPath, 0775)
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
		absPath, err := filepath.Abs(filepath.Join(args.Path, rf.relPath))
		if err != nil {
			return fmt.Errorf("Failed to determine local absolute path: %s", err)
		}
		fmt.Fprintf(args.Out, "[%04d/%04d] Downloading %s -> %s\n", i+1, missingCount, rf.relPath, filepath.Join(filepath.Base(args.Path), rf.relPath))

		err = self.downloadRemoteFile(rf.file.Id, absPath, args, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Drive) downloadChangedFiles(changedFiles []*changedFile, args DownloadSyncArgs) error {
	changedCount := len(changedFiles)

	if changedCount > 0 {
		fmt.Fprintf(args.Out, "\n%d remote files has changed\n", changedCount)
	}

	for i, cf := range changedFiles {
		if skip, reason := checkLocalConflict(cf, args.Resolution); skip {
			fmt.Fprintf(args.Out, "[%04d/%04d] Skipping %s (%s)\n", i+1, changedCount, cf.remote.relPath, reason)
			continue
		}

		absPath, err := filepath.Abs(filepath.Join(args.Path, cf.remote.relPath))
		if err != nil {
			return fmt.Errorf("Failed to determine local absolute path: %s", err)
		}
		fmt.Fprintf(args.Out, "[%04d/%04d] Downloading %s -> %s\n", i+1, changedCount, cf.remote.relPath, filepath.Join(filepath.Base(args.Path), cf.remote.relPath))

		err = self.downloadRemoteFile(cf.remote.file.Id, absPath, args, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Drive) downloadRemoteFile(id, fpath string, args DownloadSyncArgs, try int) error {
	if args.DryRun {
		return nil
	}

	// Get timeout reader wrapper and context
	timeoutReaderWrapper, ctx := getTimeoutReaderWrapperContext(args.Timeout)

	res, err := self.service.Files.Get(id).SupportsTeamDrives(true).Context(ctx).Download()
	if err != nil {
		if isBackendOrRateLimitError(err) && try < MaxErrorRetries {
			exponentialBackoffSleep(try)
			try++
			return self.downloadRemoteFile(id, fpath, args, try)
		} else if isTimeoutError(err) {
			return fmt.Errorf("Failed to download file: timeout, no data was transferred for %v", args.Timeout)
		} else {
			return fmt.Errorf("Failed to download file: %s", err)
		}
	}

	// Close body on function exit
	defer res.Body.Close()

	// Wrap response body in progress reader
	progressReader := getProgressReader(res.Body, args.Progress, res.ContentLength)

	// Wrap reader in timeout reader
	reader := timeoutReaderWrapper(progressReader)

	// Ensure any parent directories exists
	if err = mkdir(fpath); err != nil {
		return err
	}

	// Download to tmp file
	tmpPath := fpath + ".incomplete"

	// Create new file
	outFile, err := os.Create(tmpPath)
	if err != nil {
		return fmt.Errorf("Unable to create local file: %s", err)
	}

	// Save file to disk
	_, err = io.Copy(outFile, reader)
	if err != nil {
		outFile.Close()
		if try < MaxErrorRetries {
			exponentialBackoffSleep(try)
			try++
			return self.downloadRemoteFile(id, fpath, args, try)
		} else {
			os.Remove(tmpPath)
			return fmt.Errorf("Download was interrupted: %s", err)
		}
	}

	// Close file
	outFile.Close()

	// Rename tmp file to proper filename
	return os.Rename(tmpPath, fpath)
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
		fmt.Fprintf(args.Out, "[%04d/%04d] Deleting %s\n", i+1, extraneousCount, lf.absPath)

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

func checkLocalConflict(cf *changedFile, resolution ConflictResolution) (bool, string) {
	// No conflict unless local file was last modified
	if cf.compareModTime() != LocalLastModified {
		return false, ""
	}

	// Don't skip if want to keep the remote file
	if resolution == KeepRemote {
		return false, ""
	}

	// Skip if we want to keep the local file
	if resolution == KeepLocal {
		return true, "conflicting file, keeping local file"
	}

	if resolution == KeepLargest {
		largest := cf.compareSize()

		// Skip if the local file is largest
		if largest == LocalLargestSize {
			return true, "conflicting file, local file is largest, keeping local"
		}

		// Don't skip if the remote file is largest
		if largest == RemoteLargestSize {
			return false, ""
		}

		// Keep local if both files have the same size
		if largest == EqualSize {
			return true, "conflicting file, file sizes are equal, keeping local"
		}
	}

	// The conditionals above should cover all cases,
	// unless the programmer did something wrong,
	// in which case we default to being non-destructive and skip the file
	return true, "conflicting file, unhandled case"
}

func ensureNoLocalModifications(files []*changedFile) error {
	conflicts := findLocalConflicts(files)
	if len(conflicts) == 0 {
		return nil
	}

	buffer := bytes.NewBufferString("")
	formatConflicts(conflicts, buffer)
	return fmt.Errorf(buffer.String())
}
