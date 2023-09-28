package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

type Params[C any] interface {
	BuildTransform(cfg *C) (asset.Applier, error)
}

type genericMetadata interface {
	parse(cp core.ConfigProvider, paramsR schema.Reader) (Transform, error)
	decodeMsg(cp core.ConfigProvider, paramsR *msgp.Reader) (Transform, error)
}

type Transform interface {
	asset.Transform
	Kind() Kind
}
