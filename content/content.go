package content

import (
	"bytes"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"image"
	"io"
	"sync"

	ref "github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/util"
)

type Loader func(io.Writer) error

func FromFile(filepath string) Loader {
	return func(w io.Writer) error {
		f, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer f.Close()
		return util.UnwrapErr(io.Copy(w, f))
	}
}

type Content interface {
	GetHash() ref.Murmur3Hash
	BytesReader() *bytes.Reader
}

type ImageContent interface {
	Content
	Image() image.Image
}

type content struct {
	loader Loader
	buf    []byte
	hash   ref.Murmur3Hash
	once   sync.Once
}

func New(loader Loader) Content {
	return &content{loader: loader}
}

func NewBytes(bytes []byte) Content {
	return &content{buf: bytes}
}

func (c *content) GetHash() ref.Murmur3Hash {
	c.load()
	return c.hash
}

func (c *content) load() {
	if c.loader != nil {
		c.once.Do(func() {
			var buf bytes.Buffer
			util.Check(c.loader(&buf))
			c.buf = buf.Bytes()

			c.hash = util.Unwrap(ref.Murmur3HashFromReader(bytes.NewReader(c.buf)))

			c.loader = nil
		})
	}
}

func (c *content) BytesReader() *bytes.Reader {
	c.load()
	return bytes.NewReader(c.buf)
}

func (c *content) GetFID(basename string) ref.FID {
	c.load()
	return ref.FIDFromParts(basename, c.hash)
}

func BuildFID(content Content, basename string) ref.FID {
	return ref.FIDFromParts(basename, content.GetHash())
}

type contentWithKnownHash struct {
	loader Loader
	buf    []byte
	hash   ref.Murmur3Hash
	once   sync.Once
}

func NewWithKnownHash(loader Loader, hash ref.Murmur3Hash) Content {
	return &contentWithKnownHash{
		loader: loader,
		hash:   hash,
	}
}

func (c *contentWithKnownHash) GetHash() ref.Murmur3Hash { return c.hash }
func (c *contentWithKnownHash) load() {
	if c.loader != nil {
		c.once.Do(func() {
			var buf bytes.Buffer
			util.Check(c.loader(&buf))
			c.buf = buf.Bytes()

			hash := util.Unwrap(ref.Murmur3HashFromReader(bytes.NewReader(c.buf)))

			if hash != c.hash {
				panic("corrupted file: not match with known hash")
			}

			c.loader = nil
		})
	}
}
func (c *contentWithKnownHash) BytesReader() *bytes.Reader {
	c.load()
	return bytes.NewReader(c.buf)
}

type imageContent struct {
	Content
	imageCache image.Image
	once       sync.Once
}

func NewImage(content Content) ImageContent {
	if ic, ok := content.(ImageContent); ok {
		return ic
	}
	return &imageContent{Content: content}
}

func NewImageWith(image image.Image) ImageContent {
	ic := &imageContent{imageCache: image}
	ic.once.Do(func() {})
	return ic
}

func (i *imageContent) Image() image.Image {
	i.once.Do(func() {
		im, _, err := image.Decode(i.BytesReader())
		util.Check(err)
		i.imageCache = im
	})
	return i.imageCache
}

func (i *imageContent) BytesReader() *bytes.Reader {
	if i.Content != nil {
		return i.Content.BytesReader()
	}
	panic("pure image must be encoded")
}
