package auth

import (
	"os"
	"path/filepath"
)

func mkdir(path string) error {
	dir := filepath.Dir(path)
	if fileExists(dir) {
		return nil
	}
	return os.Mkdir(dir, 0700)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}
