package gift

import (
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type sigmoid struct {
	apHead[*sigmoid]
	Factor   *big.Rat `imbed:"factor"`
	Midpoint *big.Rat `imbed:"mid,@/util/rats!rats.R1_2"`
}

func (s *sigmoid) filter() gift.Filter {
	f, _ := s.Factor.Float32()
	m, _ := s.Midpoint.Float32()
	return gift.Sigmoid(m, f)
}

func registerSigmoid(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.sigmoid",
		schema.ZSTSchema.Build(),
		sigmoidSchema.Build(),
	).
		Alias("sigmoid").
		Category(Category)
}
