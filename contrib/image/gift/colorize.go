package gift

import (
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type colorize struct {
	applierHead[*colorize]
	H *big.Rat `imbed:"h"`
	S *big.Rat `imbed:"s"`
	P *big.Rat `imbed:"p"`
}

func (c *colorize) filter() gift.Filter {
	h, _ := c.H.Float32()
	s, _ := c.S.Float32()
	p, _ := c.P.Float32()
	return gift.Colorize(h, s, p)
}

func (c *colorize) BuildTransform(*schema.ZST) (transform.Applier, error) {
	return c, nil
}

func registerColorize(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.colorize",
		schema.ZSTSchema.Build(),
		colorizeSchema.Build(),
	).
		Alias("colorize").
		Category(Category)
}
