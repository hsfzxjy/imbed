package util

import (
	"io"
	"os"
)

func WriteFile(filepath string, r io.Reader) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	return UnwrapErr(io.Copy(f, r))
}
