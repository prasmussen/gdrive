package drive

import (
	"fmt"
	"io"
)

type DeleteArgs struct {
	Out       io.Writer
	Id        string
	Recursive bool
}

func (self *Drive) Delete(args DeleteArgs, try int) error {
	f, err := self.service.Files.Get(args.Id).Fields("name", "mimeType").Do()
	if err != nil {
		if isBackendOrRateLimitError(err) && try < MaxErrorRetries {
			exponentialBackoffSleep(try)
			try++
			return self.Delete(args, try)
		}
		return fmt.Errorf("Failed to get file: %s", err)
	}

	if isDir(f) && !args.Recursive {
		return fmt.Errorf("'%s' is a directory, use the 'recursive' flag to delete directories", f.Name)
	}

	err = self.service.Files.Delete(args.Id).Do()
	if err != nil {
		if isBackendOrRateLimitError(err) && try < MaxErrorRetries {
			exponentialBackoffSleep(try)
			try++
			return self.Delete(args, try)
		}
		return fmt.Errorf("Failed to delete file: %s", err)
	}

	fmt.Fprintf(args.Out, "Deleted '%s'\n", f.Name)
	return nil
}

func (self *Drive) deleteFile(fileId string, try int) error {
	err := self.service.Files.Delete(fileId).Do()
	if err != nil {
		if isBackendOrRateLimitError(err) && try < MaxErrorRetries {
			exponentialBackoffSleep(try)
			try++
			return self.deleteFile(fileId, try)
		}
		return fmt.Errorf("Failed to delete file: %s", err)
	}
	return nil
}
