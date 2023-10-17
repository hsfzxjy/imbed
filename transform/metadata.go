package transform

import (
	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

type paramsWithMetadata[C any, P IParam[C, A], A IApplier] struct {
	metadata *metadata[C, P, A]
	params   P
}

func (pm *paramsWithMetadata[C, P, A]) Metadata() Metadata { return pm.metadata }
func (pm *paramsWithMetadata[C, P, A]) VisitParams(v schema.Visitor) error {
	return pm.metadata.paramsSchema.Visit(v, pm.params)
}
func (pm *paramsWithMetadata[C, P, A]) BuildWith(cfgBuilder ConfigBuilder, cp core.ConfigProvider) (Transform, error) {
	if b, ok := cfgBuilder.(configBuilderTyped[C]); ok {
		cfg, err := b.buildConfig(cp)
		if err != nil {
			return nil, err
		}
		applier, err := (pm.params).BuildTransform(cfg)
		if err != nil {
			return nil, err
		}
		return newSingleTransform(pm.metadata, cfg, pm.params, applier), nil
	} else {
		panic("configBuilder is of wrong type")
	}
}

type metadata[C any, P IParam[C, A], A IApplier] struct {
	name    string
	aliases []string

	configSchema  schema.Schema[C]
	paramsSchema  schema.Schema[P]
	applierSchema schema.Schema[A]

	kind Kind
}

func (m *metadata[C, P, A]) Name() string {
	return m.name
}

func (m *metadata[C, P, A]) ScanParams(paramsR schema.Scanner) (ParamsWithMetadata, error) {
	params, err := m.paramsSchema.ScanFrom(paramsR)
	if err != nil {
		return nil, paramsR.Error(err)
	}
	return &paramsWithMetadata[C, P, A]{metadata: m, params: params}, nil
}

func (m *metadata[C, P, A]) decodeMsg(msgR *msgp.Reader) (ParamsWithMetadata, error) {
	params, err := m.paramsSchema.DecodeMsg(msgR)
	if err != nil {
		return nil, err
	}
	return &paramsWithMetadata[C, P, A]{metadata: m, params: params}, nil
}

func (m *metadata[C, P, A]) ConfigBuilderWorkspace() ConfigBuilder {
	return configBuilderWorkspace[C, P, A]{m}
}

func (m *metadata[C, P, A]) ConfigBuilderNeedle(n ndl.Needle) ConfigBuilder {
	return &configBuilderNeedle[C, P, A]{m, n}
}

func (m *metadata[C, P, A]) ConfigBuilderHash(h ref.Sha256Hash) ConfigBuilder {
	return &configBuilderHash[C, P, A]{m, h}
}
