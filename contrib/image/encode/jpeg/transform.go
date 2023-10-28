package jpeg

import (
	"fmt"
	"image/jpeg"
	"io"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type Applier struct {
	Quality int64 `imbed:""`
}

func (x *Applier) Apply(app core.App, a asset.Asset) (asset.Update, error) {
	ic := content.AsImage(a.Content())
	c := content.New(content.WithLoadFunc(func(w io.Writer) error {
		im, err := ic.Image()
		if err != nil {
			return err
		}
		return jpeg.Encode(w, im.Image, &jpeg.Options{
			Quality: int(x.Quality),
		})
	}))
	return asset.MergeUpdates(
		asset.UpdateContent(c),
		asset.UpdateFileExtension(".jpg"),
	), nil
}

//imbed:schemagen image.encode.jpeg#Config
type Config struct {
	DefaultQuality int64 `imbed:"default_quality,75"`
}

func (c *Config) Validate() error {
	if c.DefaultQuality < 0 || c.DefaultQuality > 100 {
		return fmt.Errorf("bad default_quality: %d", c.DefaultQuality)
	}
	return nil
}

//imbed:schemagen image.encode.jpeg#Params
type Params struct {
	Quality int64 `imbed:"q,-1"`
}

func (p *Params) Validate() error {
	if p.Quality < -1 || p.Quality > 100 {
		return fmt.Errorf("bad quality: %d", p.Quality)
	}
	return nil
}

func (p *Params) BuildTransform(c *Config) (transform.Applier, error) {
	var q = p.Quality
	if q == -1 {
		q = c.DefaultQuality
	}
	return &Applier{Quality: q}, nil
}

func Register(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.encode.jpeg",
		ConfigSchema.Build(),
		ParamsSchema.Build(),
	).
		Alias("jpeg", "jpg")
}
