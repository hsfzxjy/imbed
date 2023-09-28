package transform

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

type configInstance[C any] struct {
	schema.Schema[C]
	Value *C
	ref.LazyEncodableObject[*configInstance[C]]
}

func (v *configInstance[C]) EncodeSelf(w *msgp.Writer) error {
	return v.EncodeMsg(w, v.Value)
}

func newConfigInstance[C any](schema schema.Schema[C], value *C) *configInstance[C] {
	v := new(configInstance[C])
	v.Schema = schema
	v.Value = value
	v.Inner = v
	return v
}
