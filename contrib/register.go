package contrib

import (
	"path"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/contrib/image"
	"github.com/hsfzxjy/imbed/contrib/upload"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util"
	"github.com/qiniu/open"
)

type openApplier struct{ schema.ZST }

func (openApplier) Apply(app core.App, a asset.Asset) (asset.Update, error) {
	c := a.Content()
	r, err := c.BytesReader()
	if err != nil {
		return nil, err
	}
	hash, err := c.GetHash()
	if err != nil {
		return nil, err
	}

	fid := hash.WithName(a.BaseName())
	filename := path.Join(app.TmpDir(), fid)
	_, err = util.SafeWriteFile(r, filename)
	return nil, open.Run(filename)
}

type openParams schema.ZST

func (*openParams) BuildTransform(*schema.ZST) (transform.Applier, error) {
	return openApplier{}, nil
}

func Register(r *transform.Registry) {
	image.Register(r)
	upload.Register(r)

	transform.RegisterIn(
		r, "core.open",
		schema.ZSTSchema.Build(),
		schema.ZSTSchemaAs[openParams]().Build(),
	).
		Alias("open", "debug").
		Category(transform.Terminal)
}
