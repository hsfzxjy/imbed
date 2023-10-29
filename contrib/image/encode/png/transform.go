package png

import (
	"image/png"
	"io"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

type PNG struct {
	schema.ZST
}

func (PNG) Apply(app core.App, a asset.Asset) (asset.Update, error) {
	ic := content.AsImage(a.Content())
	c := content.New(content.WithLoadFunc(func(w io.Writer) error {
		im, err := ic.Image()
		if err != nil {
			return err
		}
		return png.Encode(w, im.Image)
	}))
	return asset.MergeUpdates(
		asset.UpdateContent(c),
		asset.UpdateFileExtension(".png"),
	), nil
}

func (PNG) BuildTransform(*schema.ZST) (transform.Applier, error) {
	return PNG{}, nil
}

func Register(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.encode.png",
		schema.ZSTSchema.Build(),
		schema.ZSTSchemaAs[PNG]().Build(),
	).
		Alias("png")
}
