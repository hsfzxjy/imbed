package gift

import (
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type gamma struct {
	apHead[*gamma]
	Gamma *big.Rat `imbed:"g"`
}

func (g *gamma) filter() gift.Filter {
	f, _ := g.Gamma.Float32()
	return gift.Gamma(f)
}

func registerGamma(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.gamma",
		schema.ZSTSchema.Build(),
		gammaSchema.Build(),
	).
		Alias("gamma").
		Category(Category)
}
