package content

import (
	"fmt"
	"io"

	"image/jpeg"
	"image/png"

	webpencoder "github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"golang.org/x/image/bmp"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/core/ref"
)

type ContentOption func(*content)

func WithPos(pos pos.P) ContentOption {
	return func(c *content) {
		c.pos = pos
	}
}

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

func WithHash(hash ref.FileHash) ContentOption {
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

func imageLoadFunc(image Image) loadFunc {
	return func(w io.Writer) error {
		switch image.SourceFormat {
		case "jpeg":
			return jpeg.Encode(w, image.Image, &jpeg.Options{Quality: 100})
		case "png":
			return png.Encode(w, image.Image)
		case "bmp":
			return bmp.Encode(w, image.Image)
		case "webp":
			return webp.Encode(w, image.Image, &webpencoder.Options{
				Lossless: true,
			})
		default:
			return fmt.Errorf("unknown SourceFormat: %q", image.SourceFormat)
		}
	}
}

func NewImage(image Image) ImageContent {
	ic := &imageContent{
		imageCache: image,
		content: &content{
			loader: imageLoadFunc(image),
		},
	}
	ic.once.Do(func() {})
	return ic
}
