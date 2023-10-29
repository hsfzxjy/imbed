package gift

import (
	"image/color"
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type rotate struct {
	apHead[*rotate]
	Deg           *big.Rat `imbed:"deg"`
	Interpolation `imbed:"itpl!string,\"lin\""`
}

func (a *rotate) filter() gift.Filter {
	deg, _ := a.Deg.Float32()
	return gift.Rotate(deg, color.Opaque, util.Unwrap(a.Interpolation.Get()))
}

func (p *rotate) Validate() error {
	if _, err := p.Interpolation.Get(); err != nil {
		return err
	}
	return nil
}

func registerRotate(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.rotate",
		schema.ZSTSchema.Build(),
		rotateSchema.Build(),
	).
		Alias("rotate", "rot").
		Category(Category)
}
