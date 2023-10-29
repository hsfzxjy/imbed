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
type brightness struct {
	applierHead[*brightness]
	Percentage *big.Rat `imbed:"p"`
}

func (a *brightness) filter() gift.Filter {
	f, _ := a.Percentage.Float32()
	return gift.Brightness(f)
}

func (p *brightness) Validate() error {
	pc := p.Percentage
	if pc.Cmp(rats.R100) > 0 || pc.Cmp(rats.RN100) < 0 {
		return fmt.Errorf("percentage must be in [-100, 100], got %s", pc.FloatString(1))
	}
	return nil
}

func (p *brightness) BuildTransform(*schema.ZST) (transform.Applier, error) {
	return p, nil
}

func registerBrightness(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.brightness",
		schema.ZSTSchema.Build(),
		brightnessSchema.Build(),
	).
		Alias("brightness", "bright").
		Category(Category)
}
