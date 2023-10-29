package gift

import (
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type sepia struct {
	apHead[*sepia]
	p *big.Rat `imbed:"p"`
}

func (s *sepia) filter() gift.Filter {
	f, _ := s.p.Float32()
	return gift.Sepia(f)
}

func registerSepia(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.sepia",
		schema.ZSTSchema.Build(),
		sepiaSchema.Build(),
	).
		Alias("sepia").
		Category(Category)
}
