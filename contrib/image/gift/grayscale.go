package gift

import (
	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type grayScale struct {
	apHead[*grayScale]
}

func (g *grayScale) filter() gift.Filter {
	return gift.Grayscale()
}

func registerGrayScale(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.gray_scale",
		schema.ZSTSchema.Build(),
		grayScaleSchema.Build(),
	).
		Alias("gray", "grey").
		Category(Category)
}
