package local

import (
	"errors"
	"path"
	"path/filepath"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen upload.local#Applier
type Applier struct {
	Path string `imbed:""`
}

func (t *Applier) Apply(app core.App, a asset.Asset) (asset.Update, error) {
	fhash, err := a.Content().GetHash()
	if err != nil {
		return nil, err
	}
	filename := fhash.WithName(a.BaseName())
	filepath := path.Join(t.Path, filename)
	r, err := a.Content().BytesReader()
	if err != nil {
		return nil, err
	}
	err = util.WriteFile(filepath, r)
	if err != nil {
		return nil, err
	}
	return asset.UpdateUrl(filepath), nil
}

//imbed:schemagen upload.local#Config
type Config struct {
	Path string `imbed:"path,\"\""`
}

//imbed:schemagen upload.local#Params
type Params struct {
	Path string `imbed:"path,\"\""`
}

func (p *Params) BuildTransform(c *Config) (transform.Applier, error) {
	var path = p.Path
	if path == "" {
		path = c.Path
	}
	if path == "" {
		return nil, errors.New("empty upload path")
	}
	var err error
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return &Applier{Path: path}, nil
}

func Register(r *transform.Registry) {
	transform.RegisterIn(
		r, "upload.local",
		ConfigSchema.Build(),
		ParamsSchema.Build(),
	).
		Category(transform.Terminal)
}
