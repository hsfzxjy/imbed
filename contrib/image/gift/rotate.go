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
	applierHead[*rotate]
	Deg           *big.Rat `imbed:"deg"`
	Interpolation string   `imbed:"itpl,\"lin\""`
}

func (a *rotate) filter() gift.Filter {
	deg, _ := a.Deg.Float32()
	return gift.Rotate(deg, color.Opaque, util.Unwrap(Interpolation(a.Interpolation)))
}

func (p *rotate) Validate() error {
	if _, err := Interpolation(p.Interpolation); err != nil {
		return err
	}
	return nil
}

func (p *rotate) BuildTransform(*schema.ZST) (transform.Applier, error) {
	return p, nil
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
