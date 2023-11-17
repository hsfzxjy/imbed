package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type ParamFor[C any] interface {
	BuildTransform(cfg C) (Applier, error)
}

type Applier interface {
	asset.Applier
	EncodeMsg(w *fastbuf.W)
}

type Composer interface {
	Compose(appliers []Applier) (Applier, error)
}
