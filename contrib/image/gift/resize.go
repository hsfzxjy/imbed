package gift

import (
	"errors"
	"fmt"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util"
)

//go:generate go run github.com/hsfzxjy/imbed/schema/gen

//imbed:schemagen
type resizeApplier struct {
	applierHead[*resizeApplier]
	resizeParams `imbed:""`
}

func (a *resizeApplier) filter() gift.Filter {
	switch a.Mode {
	case "default":
		return gift.Resize(
			int(a.W), int(a.H),
			util.Unwrap(Resampling(a.Resampling)),
		)
	case "fit":
		return gift.ResizeToFit(
			int(a.W), int(a.H),
			util.Unwrap(Resampling(a.Resampling)),
		)
	case "fill":
		return gift.ResizeToFill(
			int(a.W), int(a.H),
			util.Unwrap(Resampling(a.Resampling)),
			util.Unwrap(Anchor(a.Anchor)),
		)
	default:
		panic("unreachable")
	}
}

//imbed:schemagen
type resizeConfig struct{}

//imbed:schemagen
type resizeParams struct {
	H int64 `imbed:"h,0"`
	W int64 `imbed:"w,0"`

	Resampling string `imbed:"sample,\"cub\""`

	Anchor string `imbed:"anchor,\"c\""`
	Mode   string `imbed:"mode,\"default\""`
}

func (p *resizeParams) Validate() error {
	if p.H <= 0 && p.W <= 0 {
		return errors.New("at least one of h and w should be positive integer (mode=default)")
	}
	if _, err := Resampling(p.Resampling); err != nil {
		return err
	}
	switch p.Mode {
	case "default":
	case "fit":
		if _, err := Anchor(p.Anchor); err != nil {
			return err
		}
		fallthrough
	case "fill":
		if p.H <= 0 || p.W <= 0 {
			return fmt.Errorf("both h and w should be positive integers (mode=%s)", p.Mode)
		}
	default:
		return fmt.Errorf("unknown Mode: %q", p.Mode)
	}
	return nil
}

func (p *resizeParams) BuildTransform(*resizeConfig) (transform.Applier, error) {
	return &resizeApplier{
		resizeParams: *p,
	}, nil
}

func regsiterResize(r *transform.Registry) {
	transform.RegisterIn(
		r, "image.resize",
		resizeConfigSchema.Build(),
		resizeParamsSchema.Build(),
	).
		Alias("resize").
		Category(Category)
}
