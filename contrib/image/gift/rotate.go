package gift

import (
	"image/color"
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type rotateApplier struct {
	applierHead[*rotateApplier]
	rotateParams `imbed:""`
}

func (a *rotateApplier) filter() gift.Filter {
	deg, _ := a.Deg.Float32()
	return gift.Rotate(deg, color.Opaque, util.Unwrap(Interpolation(a.Interpolation)))
}

//imbed:schemagen
type rotateConfig struct{}

//imbed:schemagen
type rotateParams struct {
	Deg           *big.Rat `imbed:"deg"`
	Interpolation string   `imbed:"itpl,\"lin\""`
}

func (p *rotateParams) Validate() error {
	if _, err := Interpolation(p.Interpolation); err != nil {
		return err
	}
	return nil
}

func (p *rotateParams) BuildTransform(*rotateConfig) (transform.Applier, error) {
	return &rotateApplier{
		rotateParams: *p,
	}, nil
}

func registerRotate(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.rotate",
		rotateConfigSchema.Build(),
		rotateParamsSchema.Build(),
	).
		Alias("rotate", "rot").
		Category(Category)
}
