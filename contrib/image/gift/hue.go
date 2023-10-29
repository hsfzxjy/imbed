package gift

import (
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type hue struct {
	apHead[*hue]
	X *big.Rat `imbed:"x"`
}

func (h *hue) filter() gift.Filter {
	x, _ := h.X.Float32()
	return gift.Hue(x)
}

func registerHue(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.hue",
		schema.ZSTSchema.Build(),
		hueSchema.Build(),
	).
		Alias("hue").
		Category(Category)
}
