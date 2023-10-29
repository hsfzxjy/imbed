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
type opacity struct {
	apHead[*opacity]
	Percentage *big.Rat `imbed:"p"`
}

func (a *opacity) filter() gift.Filter {
	return a
}

func (a *opacity) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	f, _ := a.Percentage.Float64()
	bit := math.Floor(f / 100 * 255)
	// 0 is fully transparent and 255 is opaque.
	alpha := uint8(bit)
	mask := image.NewUniform(color.Alpha{alpha})
	draw.DrawMask(dst, dst.Bounds(), src, image.Point{}, mask, image.Point{}, draw.Over)
}

func (a *opacity) Bounds(srcBounds image.Rectangle) image.Rectangle {
	return image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
}

func (p *opacity) Validate() error {
	v := p.Percentage
	if v.Sign() < 0 || v.Cmp(rats.R100) > 0 {
		return fmt.Errorf("opacity must be in [0, 100], got %s", v.FloatString(1))
	}
	return nil
}

func registerOpacity(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.opacity",
		schema.ZSTSchema.Build(),
		opacitySchema.Build(),
	).
		Alias("opacity").
		Category(Category)
}
