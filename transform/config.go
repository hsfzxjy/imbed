package transform

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
)

type configWrapper[C any] struct {
	schema.Schema[C]
	Value C

	encodable
}

func (v *configWrapper[C]) compute() {
	buf, err := schema.EncodeBytes(v.Schema, v.Value)
	if err != nil {
		v.encodeError = err
		return
	}
	v.encoded = buf
	v.hash = ref.Sha256HashSum(buf)
}

func (v *configWrapper[C]) GetRawEncoded() ([]byte, error) {
	v.compute()
	return v.encoded, v.encodeError
}

func (v *configWrapper[C]) GetSha256Hash() (ref.Sha256Hash, error) {
	v.compute()
	return v.hash, v.encodeError
}

func wrapConfig[C any](schema schema.Schema[C], value C) *configWrapper[C] {
	v := new(configWrapper[C])
	v.Schema = schema
	v.Value = value
	v.Compute = v.compute
	return v
}
