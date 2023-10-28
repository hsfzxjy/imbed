package content

import (
	"bytes"
	"io"

	"github.com/hsfzxjy/imbed/core/ref"
)

type Loader interface {
	Load(w io.Writer) error
}

type Sizer interface {
	Size() (Size, error)
}

type LoadSizer interface {
	Loader
	Sizer
}

type Content interface {
	GetHash() (ref.Murmur3Hash, error)
	BytesReader() (*bytes.Reader, error)
	Sizer
}

type ImageContent interface {
	Content
	Image() (Image, error)
}

func BuildFID(content Content, basename string) (ref.FID, error) {
	hash, err := content.GetHash()
	if err != nil {
		return ref.FID{}, err
	}
	return ref.FIDFromParts(basename, hash), nil
}
