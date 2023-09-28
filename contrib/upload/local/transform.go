package local

import (
	"errors"
	"path"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util"
)

type applier struct {
	Path string
}

func (t applier) Apply(a asset.Asset) (asset.Update, error) {
	filename := content.BuildFID(a.Content(), a.BaseName()).Humanize()
	filepath := path.Join(t.Path, filename)
	err := util.WriteFile(filepath, a.Content().BytesReader())
	if err != nil {
		return nil, err
	}
	return asset.UpdateUrl(filepath), nil
}

type Config struct {
	Path string
}

func (c *Config) Validate() error {
	if c.Path == "" {
		return errors.New("empty upload path")
	}
	return nil
}

type Params struct{}

func (Params) BuildTransform(c *Config) (asset.Applier, error) {
	return applier{c.Path}, nil
}

func Register(r transform.Registry) {
	var c Config
	var p Params
	transform.RegisterIn(
		r,
		"upload.local",
		schema.Struct(&c,
			schema.F("path", &c.Path, schema.String()),
		).Build(),
		schema.Struct(&p).Build(),
	).Kind(transform.KindPersist)
}
