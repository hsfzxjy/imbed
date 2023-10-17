package transform

import (
	"bytes"
	"slices"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/tinylib/msgp/msgp"
)

type singleTransform[C any, P IParam[C, A], A IApplier] struct {
	*metadata[C, P, A]
	config  *configWrapper[C]
	params  P
	applier A

	encodable
}

func (t *singleTransform[C, P, A]) Apply(app core.App, asset asset.Asset) (asset.Update, error) {
	return t.applier.Apply(app, asset)
}

func (t *singleTransform[C, P, A]) compute() {
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
	err = t.applierSchema.EncodeMsg(w, t.applier)
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

func (t *singleTransform[C, P, A]) AssociatedConfigs() []ref.EncodableObject {
	return []ref.EncodableObject{t.config}
}

func (t *singleTransform[C, P, A]) Name() string {
	return t.name
}

func (t *singleTransform[C, P, A]) Kind() Kind {
	return t.metadata.kind
}

func newSingleTransform[C any, P IParam[C, A], A IApplier](
	metadata *metadata[C, P, A],
	cfg C, params P,
	applier A,
) *singleTransform[C, P, A] {
	t := new(singleTransform[C, P, A])
	t.applier = applier
	t.metadata = metadata
	t.config = wrapConfig(metadata.configSchema, cfg)
	t.params = params
	t.Compute = t.compute
	return t
}
