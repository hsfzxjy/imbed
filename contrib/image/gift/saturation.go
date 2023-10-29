package gift

import (
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type saturation struct {
	apHead[*saturation]
	p *big.Rat `imbed:"p"`
}

func (s *saturation) filter() gift.Filter {
	f, _ := s.p.Float32()
	return gift.Saturation(f)
}

func registerSaturation(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.saturation",
		schema.ZSTSchema.Build(),
		saturationSchema.Build(),
	).
		Alias("sat", "saturation").
		Category(Category)
}
