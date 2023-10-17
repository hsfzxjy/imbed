package transform

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/tinylib/msgp/msgp"
)

type configBuilderTyped[C any] interface {
	ConfigBuilder
	buildConfig(cp core.ConfigProvider) (C, error)
}

type ConfigBuilder interface {
	ConfigHash() ref.Sha256Hash
}

type configBuilderWorkspace[C any, P IParam[C, A], A IApplier] struct{ *metadata[C, P, A] }

func (b configBuilderWorkspace[C, P, A]) ConfigHash() (zero ref.Sha256Hash) { return }

func (b configBuilderWorkspace[C, P, A]) buildConfig(cp core.ConfigProvider) (result C, err error) {
	cfgR, err := cp.ProvideWorkspaceConfig(b.metadata.name)
	if err != nil {
		return
	}
	cfg, err := b.metadata.configSchema.ScanFrom(cfgR)
	if err != nil {
		return
	}
	return cfg, nil
}

type configBuilderNeedle[C any, P IParam[C, A], A IApplier] struct {
	*metadata[C, P, A]
	needle ndl.Needle
}

func (b *configBuilderNeedle[C, P, A]) ConfigHash() (zero ref.Sha256Hash) { return }

func (b *configBuilderNeedle[C, P, A]) buildConfig(cp core.ConfigProvider) (result C, err error) {
	buf, err := cp.ProvideStockConfig(b.needle)
	if err != nil {
		return
	}
	cfgR := msgp.NewReader(bytes.NewReader(buf))
	cfg, err := b.metadata.configSchema.DecodeMsg(cfgR)
	if err != nil {
		return
	}
	return cfg, nil
}

type configBuilderHash[C any, P IParam[C, A], A IApplier] struct {
	*metadata[C, P, A]
	hash ref.Sha256Hash
}

func (b *configBuilderHash[C, P, A]) ConfigHash() ref.Sha256Hash { return b.hash }

func (b *configBuilderHash[C, P, A]) buildConfig(cp core.ConfigProvider) (result C, err error) {
	buf, err := cp.ProvideStockConfig(ndl.RawFull(ref.AsRawString(b.hash)))
	if err != nil {
		return
	}
	cfgR := msgp.NewReader(bytes.NewReader(buf))
	cfg, err := b.metadata.configSchema.DecodeMsg(cfgR)
	if err != nil {
		return
	}
	return cfg, nil
}

type builder struct {
	ConfigBuilder
	ParamsWithMetadata
}

func (b *builder) Build(cp core.ConfigProvider) (Transform, error) {
	return b.BuildWith(b.ConfigBuilder, cp)
}
