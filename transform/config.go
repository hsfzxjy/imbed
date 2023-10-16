package transform

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
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

type configBuilderTyped[C any] interface {
	ConfigBuilder
	buildConfig(cp core.ConfigProvider) (*C, error)
}

type ConfigBuilder interface {
	ConfigHash() ref.Sha256Hash
}

type configBuilderWorkspace[C any, P ParamStruct[C]] struct{ *metadata[C, P] }

func (b configBuilderWorkspace[C, P]) ConfigHash() (zero ref.Sha256Hash) { return }

func (b configBuilderWorkspace[C, P]) buildConfig(cp core.ConfigProvider) (*C, error) {
	cfgR, err := cp.ProvideWorkspaceConfig(b.metadata.name)
	if err != nil {
		return nil, err
	}
	var cfg C
	err = b.metadata.configSchema.DecodeValue(cfgR, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type configBuilderNeedle[C any, P ParamStruct[C]] struct {
	*metadata[C, P]
	needle ndl.Needle
}

func (b *configBuilderNeedle[C, P]) ConfigHash() (zero ref.Sha256Hash) { return }

func (b *configBuilderNeedle[C, P]) buildConfig(cp core.ConfigProvider) (*C, error) {
	buf, err := cp.ProvideStockConfig(b.needle)
	if err != nil {
		return nil, err
	}
	var cfg C
	cfgR := msgp.NewReader(bytes.NewReader(buf))
	err = b.metadata.configSchema.DecodeMsg(cfgR, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type configBuilderHash[C any, P ParamStruct[C]] struct {
	*metadata[C, P]
	hash ref.Sha256Hash
}

func (b *configBuilderHash[C, P]) ConfigHash() ref.Sha256Hash { return b.hash }

func (b *configBuilderHash[C, P]) buildConfig(cp core.ConfigProvider) (*C, error) {
	buf, err := cp.ProvideStockConfig(ndl.RawFull(ref.AsRawString(b.hash)))
	if err != nil {
		return nil, err
	}
	var cfg C
	cfgR := msgp.NewReader(bytes.NewReader(buf))
	err = b.metadata.configSchema.DecodeMsg(cfgR, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

type builder struct {
	ConfigBuilder
	ParamsWithMetadata
}

func (b *builder) Build(cp core.ConfigProvider) (Transform, error) {
	return b.BuildWith(b.ConfigBuilder, cp)
}
