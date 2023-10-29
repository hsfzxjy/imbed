package gift

import (
	"fmt"
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type gaussianBlur struct {
	apHead[*gaussianBlur]
	Sigma *big.Rat `imbed:"s"`
}

func (g *gaussianBlur) filter() gift.Filter {
	s, _ := g.Sigma.Float32()
	return gift.GaussianBlur(s)
}

func (g *gaussianBlur) Validate() error {
	if g.Sigma.Sign() <= 0 {
		return fmt.Errorf("sigma must be >0, got %s", g.Sigma.FloatString(1))
	}
	return nil
}

func registerGaussianBlur(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.gaussian_blur",
		schema.ZSTSchema.Build(),
		gaussianBlurSchema.Build(),
	).
		Alias("gaussian", "blur").
		Category(Category)
}
