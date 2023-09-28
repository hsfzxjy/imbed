package image

import (
	"github.com/hsfzxjy/imbed/contrib/image/encode"
	"github.com/hsfzxjy/imbed/transform"
)

func Register(r transform.Registry) {
	encode.Register(r)
}
