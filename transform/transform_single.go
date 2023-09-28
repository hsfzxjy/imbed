package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/tinylib/msgp/msgp"
)

type singleTransform[C any, P Params[C]] struct {
	asset.Applier
	*metadata[C, P]

	config *configInstance[C]

	params *P

	ref.LazyEncodableObject[*singleTransform[C, P]]
}

func newSingleTransform[C any, P Params[C]](
	metadata *metadata[C, P],
	cfg *C, params *P,
	applier asset.Applier,
) *singleTransform[C, P] {
	t := new(singleTransform[C, P])
	t.Applier = applier
	t.metadata = metadata
	t.config = newConfigInstance(metadata.configSchema, cfg)
	t.params = params
	t.LazyEncodableObject.Inner = t
	return t
}

func (t *singleTransform[C, P]) EncodeSelf(w *msgp.Writer) error {
	configHash, err := t.config.GetSha256Hash()
	if err != nil {
		return err
	}
	err = w.Append(ref.AsRaw(configHash)...)
	if err != nil {
		return err
	}
	err = w.WriteString(t.metadata.name)
	if err != nil {
		return err
	}
	err = t.paramsSchema.EncodeMsg(w, t.params)
	return err
}

func (t *singleTransform[C, P]) AssociatedConfigs() []ref.EncodableObject {
	return []ref.EncodableObject{t.config}
}

func (t *singleTransform[C, P]) Name() string {
	return t.name
}

func (t *singleTransform[C, P]) Kind() Kind { return t.metadata.kind }
