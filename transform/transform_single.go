package transform

import (
	"bytes"
	"slices"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/tinylib/msgp/msgp"
)

type singleTransform[C any, P ParamFor[C]] struct {
	*metadata[C, P]
	config *configWrapper[C]
	params P
	Applier

	encodable
}

func (t *singleTransform[C, P]) compute() {
	var (
		buf     bytes.Buffer
		w       *msgp.Writer
		encoded []byte
		hash    ref.Sha256Hash
	)
	configHash, err := t.config.GetSha256Hash()
	if err != nil {
		goto ERROR
	}

	// compute encoded
	w = msgp.NewWriter(&buf)
	err = w.Append(ref.AsRaw(configHash)...)
	if err != nil {
		goto ERROR
	}
	err = w.WriteString(t.metadata.name)
	if err != nil {
		goto ERROR
	}
	err = t.paramsSchema.EncodeMsg(w, t.params)
	if err != nil {
		goto ERROR
	}
	err = w.Flush()
	if err != nil {
		goto ERROR
	}
	encoded = slices.Clone(buf.Bytes())

	// compute hash
	buf.Reset()
	err = w.WriteString(t.metadata.name)
	if err != nil {
		goto ERROR
	}
	err = t.Applier.EncodeMsg(t.Registry, w)
	if err != nil {
		goto ERROR
	}
	err = w.Flush()
	if err != nil {
		goto ERROR
	}
	hash = ref.Sha256HashSum(buf.Bytes())

	t.encoded, t.hash = encoded, hash
	return

ERROR:
	t.encodeError = err

}

func (t *singleTransform[C, P]) AssociatedConfigs() []ref.EncodableObject {
	return []ref.EncodableObject{t.config}
}

func (t *singleTransform[C, P]) Name() string {
	return t.name
}

func (t *singleTransform[C, P]) Kind() Kind {
	return t.metadata.kind
}

func newSingleTransform[C any, P ParamFor[C]](
	metadata *metadata[C, P],
	cfg C, params P,
	applier Applier,
) *singleTransform[C, P] {
	t := new(singleTransform[C, P])
	t.Applier = applier
	t.metadata = metadata
	t.config = wrapConfig(metadata.configSchema, cfg)
	t.params = params
	t.Compute = t.compute
	return t
}
