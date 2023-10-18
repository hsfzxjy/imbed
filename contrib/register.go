package contrib

import (
	"github.com/hsfzxjy/imbed/contrib/image"
	"github.com/hsfzxjy/imbed/contrib/upload"
	"github.com/hsfzxjy/imbed/transform"
)

func Register(r *transform.Registry) {
	image.Register(r)
	upload.Register(r)
}
