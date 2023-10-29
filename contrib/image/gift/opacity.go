package gift

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util/rats"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type opacityApplier struct {
	applierHead[*opacityApplier]
	alphaBit int64 `imbed:""`
}

func (a *opacityApplier) filter() gift.Filter {
	return a
}

func (a *opacityApplier) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	// 0 is fully transparent and 255 is opaque.
	alpha := uint8(a.alphaBit)
	mask := image.NewUniform(color.Alpha{alpha})
	draw.DrawMask(dst, dst.Bounds(), src, image.Point{}, mask, image.Point{}, draw.Over)
}

func (a *opacityApplier) Bounds(srcBounds image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
}

//imbed:schemagen
type opacityParams struct {
	Value *big.Rat `imbed:"v"`
}

func (p *opacityParams) Validate() error {
	v := p.Value
	if v.Sign() < 0 || v.Cmp(rats.R100) > 0 {
		return fmt.Errorf("opacity must be in [0, 100], got %s", v.FloatString(1))
	}
	return nil
}

func (p *opacityParams) BuildTransform(*schema.ZST) (transform.Applier, error) {
	f, _ := p.Value.Float64()
	bit := math.Floor(f / 100 * 255)
	return &opacityApplier{
		alphaBit: int64(bit),
	}, nil
}

func registerOpacity(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.opacity",
		schema.ZSTSchema.Build(),
		opacityParamsSchema.Build(),
	).
		Alias("opacity").
		Category(Category)
}
