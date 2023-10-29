package gift

import (
	"fmt"
	"math"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type pixelate struct {
	apHead[*pixelate]
	Size int64 `imbed:"size"`
}

func (p *pixelate) filter() gift.Filter {
	return gift.Pixelate(int(p.Size))
}

func (p *pixelate) Validate() error {
	if p.Size <= 0 {
		return fmt.Errorf("size must be positive integer")
	}
	if p.Size > math.MaxInt {
		return fmt.Errorf("size too large")
	}
	return nil
}

func registerPixelate(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.pixelate",
		schema.ZSTSchema.Build(),
		pixelateSchema.Build(),
	).
		Alias("pixelate", "pixel").
		Category(Category)
}
