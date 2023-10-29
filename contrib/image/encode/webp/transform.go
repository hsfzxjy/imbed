package webp

import (
	"fmt"
	"io"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

type WebpHint string

func (h WebpHint) Get() (encoder.EncodingPreset, error) {
	switch h {
	case "picture":
		return encoder.PresetPicture, nil
	case "photo":
		return encoder.PresetPhoto, nil
	case "drawing":
		return encoder.PresetDrawing, nil
	case "icon":
		return encoder.PresetIcon, nil
	case "text":
		return encoder.PresetText, nil
	default:
		return 0, fmt.Errorf("unknown webp hint: %q", h)
	}
}

type Applier struct {
	Config
}

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
		return webp.Encode(w, im.Image, cfg)
	}))
	return asset.MergeUpdates(
		asset.UpdateContent(c),
		asset.UpdateFileExtension(".webp"),
	), nil
}

//imbed:schemagen
type Config struct {
	WebpHint `imbed:"default_hint!string,\"photo\""`
	Quality  int64 `imbed:"default_quality,75"`
}

func (o *Config) Validate() error {
	if _, err := o.WebpHint.Get(); err != nil {
		return err
	}
	if o.Quality <= 0 || o.Quality > 110 {
		return fmt.Errorf("quality must be in range [1, 110], got %d", o.Quality)
	}
	return nil
}

func (o *Config) Get() (*encoder.Options, error) {
	var ret *encoder.Options
	preset, err := o.WebpHint.Get()
	if err != nil {
		return nil, err
	}
	if o.Quality <= 100 {
		ret, err = encoder.NewLossyEncoderOptions(preset, float32(o.Quality))
	} else {
		ret, err = encoder.NewLosslessEncoderOptions(preset, int(o.Quality-101))
	}
	if err != nil {
		return nil, err
	}
	ret.UseSharpYuv = true
	return ret, nil
}

//imbed:schemagen
type Params struct {
	WebpHint `imbed:"hint!string,\"\""`
	Quality  int64 `imbed:"q,-1"`
}

func (p *Params) Validate() error {
	if p.WebpHint != "" {
		if _, err := p.WebpHint.Get(); err != nil {
			return err
		}
	}
	if p.Quality != -1 {
		if p.Quality <= 0 || p.Quality > 110 {
			return fmt.Errorf("quality must be in range [1, 110], got %d", p.Quality)
		}
	}
	return nil
}

func (p *Params) BuildTransform(config *Config) (transform.Applier, error) {
	applier := new(Applier)
	applier.Config = (Config)(*p)
	if applier.Quality == -1 {
		applier.Quality = config.Quality
	}
	if applier.WebpHint == "" {
		applier.WebpHint = config.WebpHint
	}
	return applier, nil
}

func Register(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.encode.webp",
		ConfigSchema.Build(),
		ParamsSchema.Build(),
	).
		Alias("webp")
}
