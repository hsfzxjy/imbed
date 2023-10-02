package util

import (
	"io"
	"os"
	"path"
	"sync"
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

type RevertFunc func()

func (f RevertFunc) Then(other RevertFunc) RevertFunc {
	return func() {
		f.Call()
		other.Call()
	}
}

func (f RevertFunc) Call() {
	if f != nil {
		f()
	}
}

func SafeWriteFile(r io.Reader, filename string) (revert RevertFunc, err error) {
	_, err = os.Stat(filename)
	switch {
	case err == nil:
		// file exists, do nothing and return
		return nil, nil
	case os.IsNotExist(err):
		// file does not exist, continue to write
	default:
		// some other error
		return nil, err
	}
	tmpf, err := os.CreateTemp("", "imbed-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpf.Name())
	err = tmpf.Chmod(0600)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(tmpf, r)
	if err != nil {
		return nil, err
	}
	err = tmpf.Close()
	if err != nil {
		return nil, err
	}
	err = os.Rename(tmpf.Name(), filename)
	if err != nil {
		return nil, err
	}
	var once sync.Once
	return func() { once.Do(func() { os.Remove(filename) }) }, nil
}
