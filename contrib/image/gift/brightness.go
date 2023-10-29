package gift

import (
	"fmt"
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util/rats"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type brightnessApplier struct {
	applierHead[*brightnessApplier]
	percentage *big.Rat `imbed:""`
}

func (a *brightnessApplier) filter() gift.Filter {
	f, _ := a.percentage.Float32()
	return gift.Brightness(f)
}

//imbed:schemagen
type brightnessParams struct {
	Percentage *big.Rat `imbed:"p"`
}

func (p *brightnessParams) Validate() error {
	pc := p.Percentage
	if pc.Cmp(rats.R100) > 0 || pc.Cmp(rats.RN100) < 0 {
		return fmt.Errorf("percentage must be in [-100, 100], got %s", pc.FloatString(1))
	}
	return nil
}

func (p *brightnessParams) BuildTransform(*schema.ZST) (transform.Applier, error) {
	return &brightnessApplier{
		percentage: p.Percentage,
	}, nil
}

func registerBrightness(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.brightness",
		schema.ZSTSchema.Build(),
		brightnessParamsSchema.Build(),
	).
		Alias("brightness", "bright").
		Category(Category)
}
