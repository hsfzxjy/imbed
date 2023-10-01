package content

import (
	"image"
	"io"

	"github.com/hsfzxjy/imbed/core/ref"
)

type ContentOption func(*content)

func WithLoader(l Loader) ContentOption {
	return func(c *content) {
		c.loader = l
	}
}

func WithSizer(s Sizer) ContentOption {
	return func(c *content) {
		c.sizer = s
	}
}

func WithFilePath(fp string) ContentOption {
	return func(c *content) {
		ls := FromFile(fp)
		c.loader = ls
		c.sizer = ls
	}
}

type loadFunc func(io.Writer) error

func (f loadFunc) Load(w io.Writer) error {
	return f(w)
}

func WithLoadFunc(f loadFunc) ContentOption {
	return WithLoader(f)
}

func WithHash(hash ref.Murmur3Hash) ContentOption {
	return func(c *content) {
		c.hash = hash
	}
}

func WithBytes(buf []byte) ContentOption {
	return func(c *content) {
		c.buf = buf
	}
}

func New(opts ...ContentOption) Content {
	c := new(content)
	for _, o := range opts {
		o(c)
	}
	return c
}

func AsImage(ic Content) ImageContent {
	switch c := ic.(type) {
	case ImageContent:
		return c
	case *content:
		return &imageContent{content: c}
	default:
		panic("unsupported type")
	}
}

func NewImage(image image.Image) ImageContent {
	ic := &imageContent{imageCache: image}
	ic.once.Do(func() {})
	return ic
}
