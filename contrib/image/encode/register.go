package encode

import (
	"github.com/hsfzxjy/imbed/contrib/image/encode/bmp"
	"github.com/hsfzxjy/imbed/contrib/image/encode/jpeg"
	"github.com/hsfzxjy/imbed/contrib/image/encode/png"
	"github.com/hsfzxjy/imbed/contrib/image/encode/tiff"
	"github.com/hsfzxjy/imbed/contrib/image/encode/webp"
	"github.com/hsfzxjy/imbed/transform"
)

func Register(r *transform.Registry) {
	jpeg.Register(r)
	png.Register(r)
	bmp.Register(r)
	webp.Register(r)
	tiff.Register(r)
}
