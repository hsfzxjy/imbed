package gift

import (
	"fmt"
	"image"
	"image/draw"
	"unsafe"

	"github.com/disintegration/gift"
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/tinylib/msgp/msgp"
)

type filterer interface {
	filter() gift.Filter
}

const Category transform.Category = "gift"

// T must be a pointer type
type applierHead[T filterer] struct{}

func (ptr *applierHead[T]) Apply(app core.App, ass asset.Asset) (asset.Update, error) {
	filter := (*(*T)(unsafe.Pointer(&ptr))).filter()
	src, err := content.AsImage(ass.Content()).Image()
	if err != nil {
		return nil, err
	}
	dst := newImage(src.Image, filter.Bounds(src.Image.Bounds()))
	filter.Draw(dst, src.Image, nil)
	src.Image = dst
	return asset.UpdateContent(content.NewImage(src)), nil
}

// T must be a pointer type
type paramsHead[T transform.Applier] struct{}

func (ptr *paramsHead[T]) BuildTransform(*schema.ZST) (transform.Applier, error) {
	applier := (*(*T)(unsafe.Pointer(&ptr)))
	return applier, nil
}

// T must be a pointer type
type apHead[T interface {
	filterer
	transform.Applier
}] struct {
	applierHead[T]
	paramsHead[T]
}

func newImage(src image.Image, newBounds image.Rectangle) draw.Image {
	var dst draw.Image
	switch src.(type) {
	case *image.RGBA:
		dst = image.NewRGBA(newBounds)
	case *image.NRGBA:
		dst = image.NewNRGBA(newBounds)
	case *image.Gray:
		dst = image.NewGray(newBounds)
	default:
		dst = image.NewNRGBA(newBounds)
	}
	return dst
}

var resamplingMap = map[string]gift.Resampling{
	"box": gift.BoxResampling,
	"cub": gift.CubicResampling,
	"lan": gift.LanczosResampling,
	"lin": gift.LinearResampling,
	"nn":  gift.NearestNeighborResampling,
}

func Resampling(x string) (gift.Resampling, error) {
	r, ok := resamplingMap[x]
	if !ok {
		return nil, fmt.Errorf("unknown Resampling: %q", x)
	}
	return r, nil
}

var anchorMap = map[string]gift.Anchor{
	"c":  gift.CenterAnchor,
	"t":  gift.TopAnchor,
	"l":  gift.LeftAnchor,
	"r":  gift.RightAnchor,
	"b":  gift.BottomAnchor,
	"tl": gift.TopLeftAnchor,
	"tr": gift.TopRightAnchor,
	"bl": gift.BottomLeftAnchor,
	"br": gift.BottomRightAnchor,
}

func Anchor(x string) (gift.Anchor, error) {
	a, ok := anchorMap[x]
	if !ok {
		return 0, fmt.Errorf("unknown Anchor: %q", x)
	}
	return a, nil
}

var interpolationMap = map[string]gift.Interpolation{
	"nn":  gift.NearestNeighborInterpolation,
	"lin": gift.LinearInterpolation,
	"cub": gift.CubicInterpolation,
}

func Interpolation(x string) (gift.Interpolation, error) {
	i, ok := interpolationMap[x]
	if !ok {
		return 0, fmt.Errorf("unknown Interpolation: %q", x)
	}
	return i, nil
}

type composer struct{}

type composedApplier struct {
	applierHead[*composedApplier]
	*gift.GIFT
}

func (a *composedApplier) EncodeMsg(*msgp.Writer) error {
	panic("unreachable")
}

func (a *composedApplier) Draw(dst draw.Image, src image.Image, options *gift.Options) {
	a.GIFT.Draw(dst, src)
}

func (a *composedApplier) filter() gift.Filter {
	return a
}

func (composer) Compose(applier []transform.Applier) (transform.Applier, error) {
	filters := make([]gift.Filter, 0, len(applier))
	for _, a := range applier {
		filters = append(filters, a.(filterer).filter())
	}
	return &composedApplier{GIFT: gift.New(filters...)}, nil
}

func Register(r *transform.Registry) {
	r.RegisterComposer(Category, composer{})
	regsiterResize(r)
	registerRotate(r)
	registerOpacity(r)
	registerBrightness(r)
	registerColorBalance(r)
	registerColorize(r)
	registerContrast(r)
	registerGamma(r)
	registerGaussianBlur(r)
	registerGrayScale(r)
	registerHue(r)
}
