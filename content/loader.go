package content

import (
	"io"
	"os"

	"github.com/hsfzxjy/imbed/util"
)

type fileLoader struct {
	filename string
}

func (l *fileLoader) Load(w io.Writer) error {
	f, err := os.Open(l.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return util.UnwrapErr(io.Copy(w, f))
}

func (l *fileLoader) Size() (Size, error) {
	fi, err := os.Stat(l.filename)
	if err != nil {
		return 0, err
	}
	return Size(fi.Size()), nil
}

func FromFile(filepath string) LoadSizer {
	return &fileLoader{filepath}
}
