package gift

import (
	"fmt"
	"math/big"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util/rats"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type colorBalance struct {
	applierHead[*colorBalance]
	r *big.Rat `imbed:"r,@/util/rats!rats.R0"`
	g *big.Rat `imbed:"g,@/util/rats!rats.R0"`
	b *big.Rat `imbed:"b,@/util/rats!rats.R0"`
}

func (x *colorBalance) Validate() error {
	for _, x := range [...]struct {
		value *big.Rat
		name  string
	}{
		{x.r, "r"},
		{x.g, "g"},
		{x.b, "b"},
	} {
		if x.value.Cmp(rats.RN100) < 0 || x.value.Cmp(rats.RN500) > 0 {
			return fmt.Errorf("%s must be in [-100, 500], got %s", x.name, x.value.FloatString(1))
		}
	}
	return nil
}

func (a *colorBalance) filter() gift.Filter {
	r, _ := a.r.Float32()
	g, _ := a.g.Float32()
	b, _ := a.b.Float32()
	return gift.ColorBalance(r, g, b)
}

func (a *colorBalance) BuildTransform(*schema.ZST) (transform.Applier, error) {
	return a, nil
}

func registerColorBalance(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.color_balance",
		schema.ZSTSchema.Build(),
		colorBalanceSchema.Build(),
	).
		Alias("balance").
		Category(Category)
}
