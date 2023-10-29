package gift

import (
	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type invert struct {
	apHead[*invert]
}

func (i *invert) filter() gift.Filter {
	return gift.Invert()
}

func registerInvert(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.invert",
		schema.ZSTSchema.Build(),
		invertSchema.Build(),
	).
		Alias("invert").
		Category(Category)
}
