package util

import (
	"os"
	"path"
)

func ReplaceExt(filename, newExt string) string {
	if newExt == "" {
		return filename
	}
	ext := path.Ext(filename)
	return filename[:len(filename)-len(ext)] + newExt
}

func IsDir(p string) bool {
	info, err := os.Stat(p)
	return err == nil && info.IsDir()
}

func IsFile(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}
