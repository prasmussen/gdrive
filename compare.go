package main

import (
	"encoding/json"
	"github.com/prasmussen/gdrive/drive"
	"os"
)

const MinCacheFileSize = 5 * 1024 * 1024

type Md5Comparer struct{}

func (self Md5Comparer) Changed(local *drive.LocalFile, remote *drive.RemoteFile) bool {
	return remote.Md5() != md5sum(local.AbsPath())
}

type CachedFileInfo struct {
	Size     int64  `json:"size"`
	Modified int64  `json:"modified"`
	Md5      string `json:"md5"`
}

func NewCachedMd5Comparer(path string) CachedMd5Comparer {
	cache := map[string]*CachedFileInfo{}

	f, err := os.Open(path)
	if err == nil {
		json.NewDecoder(f).Decode(&cache)
	}
	f.Close()
	return CachedMd5Comparer{path, cache}
}

type CachedMd5Comparer struct {
	path  string
	cache map[string]*CachedFileInfo
}

func (self CachedMd5Comparer) Changed(local *drive.LocalFile, remote *drive.RemoteFile) bool {
	return remote.Md5() != self.md5(local)
}

func (self CachedMd5Comparer) md5(local *drive.LocalFile) string {
	// See if file exist in cache
	cached, found := self.cache[local.AbsPath()]

	// If found and modification time and size has not changed, return cached md5
	if found && local.Modified().UnixNano() == cached.Modified && local.Size() == cached.Size {
		return cached.Md5
	}

	// Calculate new md5 sum
	md5 := md5sum(local.AbsPath())

	// Cache file info if file meets size criteria
	if local.Size() > MinCacheFileSize {
		self.cacheAdd(local, md5)
		self.persist()
	}

	return md5
}

func (self CachedMd5Comparer) cacheAdd(lf *drive.LocalFile, md5 string) {
	self.cache[lf.AbsPath()] = &CachedFileInfo{
		Size:     lf.Size(),
		Modified: lf.Modified().UnixNano(),
		Md5:      md5,
	}
}

func (self CachedMd5Comparer) persist() {
	writeJson(self.path, self.cache)
}
