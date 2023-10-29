package bmp

import (
	"io"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
	"golang.org/x/image/bmp"
)

type BMP struct {
	schema.ZST
}

func (BMP) Apply(app core.App, a asset.Asset) (asset.Update, error) {
	ic := content.AsImage(a.Content())
	c := content.New(content.WithLoadFunc(func(w io.Writer) error {
		im, err := ic.Image()
		if err != nil {
			return err
		}
		return bmp.Encode(w, im.Image)
	}))
	return asset.MergeUpdates(
		asset.UpdateContent(c),
		asset.UpdateFileExtension(".bmp"),
	), nil
}

func (BMP) BuildTransform(*schema.ZST) (transform.Applier, error) {
	return BMP{}, nil
}

func Register(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.encode.bmp",
		schema.ZSTSchema.Build(),
		schema.ZSTSchemaAs[BMP]().Build(),
	).
		Alias("bmp")
}
