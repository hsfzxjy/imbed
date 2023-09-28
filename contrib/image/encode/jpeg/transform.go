package jpeg

import (
	"fmt"
	"image/jpeg"
	"io"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

type applier struct {
	Quality int
}

func (x applier) Apply(app core.App, a asset.Asset) (asset.Update, error) {
	ic := content.NewImage(a.Content())
	c := content.New(func(w io.Writer) error {
		return jpeg.Encode(w, ic.Image(), &jpeg.Options{
			Quality: x.Quality,
		})
	})
	return asset.MergeUpdates(
		asset.UpdateContent(c),
		asset.UpdateFileExtension(".jpg"),
	), nil
}

type Config struct {
	DefaultQuality int64
}

func (c *Config) Validate() error {
	if c.DefaultQuality < 0 || c.DefaultQuality > 100 {
		return fmt.Errorf("bad default_quality: %d", c.DefaultQuality)
	}
	return nil
}

type Params struct {
	Quality int64
}

func (p *Params) Validate() error {
	if p.Quality < -1 || p.Quality > 100 {
		return fmt.Errorf("bad quality: %d", p.Quality)
	}
	return nil
}

func (p Params) BuildTransform(c *Config) (asset.Applier, error) {
	var q = p.Quality
	if q == -1 {
		q = c.DefaultQuality
	}
	return applier{int(q)}, nil
}

func Register(r transform.Registry) {
	var c Config
	var p Params
	transform.
		RegisterIn(
			r,
			"image.encode.jpeg",
			schema.Struct(&c,
				schema.F(
					"default_quality",
					&c.DefaultQuality,
					schema.Int().
						Default(jpeg.DefaultQuality)),
			).Build(),
			schema.Struct(&p,
				schema.F(
					"q",
					&p.Quality,
					schema.Int().
						Default(-1)),
			).Build(),
		).
		Alias("jpeg", "jpg").
		Kind(transform.KindChangeContent)
}
