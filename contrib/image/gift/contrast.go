package gift

import (
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type contrast struct {
	apHead[*contrast]
	Percentage *big.Rat `imbed:"p"`
}

func (c *contrast) filter() gift.Filter {
	p, _ := c.Percentage.Float32()
	return gift.Contrast(p)
}

func registerContrast(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.contrast",
		schema.ZSTSchema.Build(),
		contrastSchema.Build(),
	).
		Alias("contrast").
		Category(Category)
}
