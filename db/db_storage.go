package db

import (
	"io"
	"os"
	"path"

	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/util"
)

func AssetPathFor(app core.App, fhash ref.FileHash) string {
	return storage{app}.AssetPathFor(fhash)
}

type storage struct {
	app core.App
}

func (s storage) AssetPathFor(fhash ref.FileHash) string {
	if fhash.IsZero() {
		panic("storage: file hash with zero value")
	}
	pretty := fhash.FmtString()
	pretty = path.Join(pretty[:2], pretty)
	return s.app.FilePath(pretty)
}

func (s *storage) CreateFile(tx *Tx, fhash ref.FileHash, r io.Reader) error {
	fpath := s.AssetPathFor(fhash)
	fdir := path.Dir(fpath)
	for {
		info, err := os.Stat(fdir)
		if os.IsNotExist(err) {
			goto MKDIR
		}
		if !info.IsDir() {
			return os.ErrExist
		}
		break
	MKDIR:
		if err = os.MkdirAll(fdir, 0o755); err != nil {
			return err
		}
	}
	revert, err := util.SafeWriteFile(r, fpath)
	tx.onRollback(revert)
	return err
}

func (s *storage) RemoveFile(fhash ref.FileHash) error {
	fpath := s.AssetPathFor(fhash)
	err := os.Remove(fpath)
	if os.IsNotExist(err) {
		err = nil
	}
	return err
}
