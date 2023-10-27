package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/tinylib/msgp/msgp"
)

type ParamFor[C any] interface {
	BuildTransform(cfg C) (Applier, error)
}

type Applier interface {
	asset.Applier
	EncodeMsg(w *msgp.Writer) error
}

type Composer interface {
	Compose(appliers []Applier) (Applier, error)
}
