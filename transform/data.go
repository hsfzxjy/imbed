package transform

import (
	"github.com/hsfzxjy/imbed/schema"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
	"github.com/tinylib/msgp/msgp"
)

type Data struct {
	*metadata
	params any
}

func (d *Data) Metadata() *metadata {
	return d.metadata
}

func (d *Data) Params() any {
	return d.params
}

func (d *Data) EncodeMsg(w *msgp.Writer) error {
	return d.paramsSchema.EncodeMsgAny(w, d.params)
}

func (d *Data) Visit(v schema.Visitor) error {
	return d.paramsSchema.VisitAny(v, d.params)
}

func (d *Data) AsBuilder(f cfgf.Factory) *Builder {
	return &Builder{d, f}
}
