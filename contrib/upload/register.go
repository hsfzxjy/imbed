package upload

import (
	"github.com/hsfzxjy/imbed/contrib/upload/imgur"
	"github.com/hsfzxjy/imbed/contrib/upload/local"
	"github.com/hsfzxjy/imbed/transform"
)

func Register(r transform.Registry) {
	local.Register(r)
	imgur.Register(r)
}
