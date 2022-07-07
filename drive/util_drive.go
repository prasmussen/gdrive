package drive

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"io"
)

func (self *Drive) listFiles(ctx context.Context, parentId string) ([]*drive.File, error) {
	args := &listAllFilesArgs{
		query:     fmt.Sprintf("'%s' in parents and trashed=false", parentId),
		fields:    []googleapi.Field{"nextPageToken", "files(id,name,md5Checksum,mimeType,size,createdTime,parents)"},
		sortOrder: "",
		maxFiles:  0,
	}
	return self.listAllFilesWithContext(ctx, args)
}

func (self *Drive) listAllFilesWithContext(ctx context.Context, args *listAllFilesArgs) ([]*drive.File, error) {
	var files []*drive.File

	var pageSize int64
	if args.maxFiles > 0 && args.maxFiles < 1000 {
		pageSize = args.maxFiles
	} else {
		pageSize = 1000
	}

	controlledStop := fmt.Errorf("Controlled stop")

	err := self.service.Files.List().
		SupportsAllDrives(true).
		IncludeItemsFromAllDrives(true).
		Q(args.query).
		Fields(args.fields...).
		OrderBy(args.sortOrder).
		PageSize(pageSize).
		Pages(ctx, func(fl *drive.FileList) error {
			files = append(files, fl.Files...)

			// Stop when we have all the files we need
			if args.maxFiles > 0 && len(files) >= int(args.maxFiles) {
				return controlledStop
			}

			return nil
		})

	if err != nil && err != controlledStop {
		return nil, err
	}

	if args.maxFiles > 0 {
		n := min(len(files), int(args.maxFiles))
		return files[:n], nil
	}

	return files, nil
}

func (self *Drive) Create(ctx context.Context, dstFile *drive.File, reader io.Reader, chunkSize googleapi.MediaOption) (*drive.File, error) {
	return self.service.Files.Create(dstFile).
		SupportsAllDrives(true).
		Fields("id", "name", "size", "md5Checksum", "webContentLink").
		Context(ctx).
		Media(reader, chunkSize).
		Do()
}

func (self *Drive) fileExists(ctx context.Context, dstFile *drive.File) (*drive.File, error) {
	parent := dstFile.Parents
	if len(parent) == 0 {
		parent = append(parent, "root")
	}
	for _, p := range parent {
		files, err := self.listFiles(ctx, p)
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			// exists same filename less than 1 parent
			if f.Name == dstFile.Name {
				return f, nil
			}
		}
	}
	return nil, nil
}
