package tiff

import (
	"fmt"
	"io"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/transform"
	"golang.org/x/image/tiff"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

type CompressionType string

func (c CompressionType) Get() (tiff.CompressionType, error) {
	switch c {
	case "none", "no":
		return tiff.Uncompressed, nil
	case "deflate":
		return tiff.Deflate, nil
	case "lzw":
		return tiff.LZW, nil
	case "ccittg3":
		return tiff.CCITTGroup3, nil
	case "ccittg4":
		return tiff.CCITTGroup4, nil
	default:
		return 0, fmt.Errorf("unknown tiff compression type: %q", c)
	}
}

type Applier struct{ Config }

func (ap *Applier) Apply(app core.App, a asset.Asset) (asset.Update, error) {
	ic := content.AsImage(a.Content())
	c := content.New(content.WithLoadFunc(func(w io.Writer) error {
		im, err := ic.Image()
		if err != nil {
			return err
		}
		cfg, err := ap.Get()
		if err != nil {
			return err
		}
		return tiff.Encode(w, im.Image, cfg)
	}))
	return asset.MergeUpdates(
		asset.UpdateContent(c),
		asset.UpdateFileExtension(".tiff"),
	), nil
}

//imbed:schemagen
type Config struct {
	Predictor bool `imbed:"predictor,false"`

	Compression CompressionType `imbed:"compression!string,\"deflate\""`
}

func (c *Config) Get() (*tiff.Options, error) {
	ct, err := c.Compression.Get()
	if err != nil {
		return nil, err
	}
	return &tiff.Options{
		Compression: ct,
		Predictor:   c.Predictor,
	}, nil
}

func (c *Config) Validate() error {
	_, err := c.Compression.Get()
	return err
}

//imbed:schemagen
type Params struct {
	Predictor   *bool           `imbed:"pred,nil"`
	Compression CompressionType `imbed:"c!string,\"\""`
}

func (p *Params) Validate() error {
	if p.Compression != "" {
		_, err := p.Compression.Get()
		return err
	}
	return nil
}

func (p *Params) BuildTransform(config *Config) (transform.Applier, error) {
	applier := new(Applier)
	if p.Predictor == nil {
		applier.Predictor = config.Predictor
	} else {
		applier.Predictor = *p.Predictor
	}
	if p.Compression != "" {
		applier.Compression = p.Compression
	} else {
		applier.Compression = config.Compression
	}
	return applier, nil
}

func Register(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.encode.tiff",
		ConfigSchema.Build(),
		ParamsSchema.Build(),
	).
		Alias("tiff")
}
