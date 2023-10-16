package transform

import (
	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

type paramsWithMetadata[C any, P ParamStruct[C]] struct {
	metadata *metadata[C, P]
	params   P
}

func (pm *paramsWithMetadata[C, P]) Metadata() Metadata { return pm.metadata }
func (pm *paramsWithMetadata[C, P]) VisitParams(v schema.Visitor) error {
	return pm.metadata.paramsSchema.Visit(v, &pm.params)
}
func (pm *paramsWithMetadata[C, P]) BuildWith(cfgBuilder ConfigBuilder, cp core.ConfigProvider) (Transform, error) {
	if b, ok := cfgBuilder.(configBuilderTyped[C]); ok {
		cfg, err := b.buildConfig(cp)
		if err != nil {
			return nil, err
		}
		applier, err := (pm.params).BuildTransform(cfg)
		if err != nil {
			return nil, err
		}
		return newSingleTransform(pm.metadata, cfg, &pm.params, applier), nil
	} else {
		panic("configBuilder is of wrong type")
	}
}

type metadata[C any, P ParamStruct[C]] struct {
	name    string
	aliases []string

	configSchema schema.Schema[C]
	paramsSchema schema.Schema[P]

	kind Kind
}

func (m *metadata[C, P]) Name() string { return m.name }

func (m *metadata[C, P]) Parse(paramsR schema.Scanner) (ParamsWithMetadata, error) {
	var params P
	err := m.paramsSchema.ScanFrom(paramsR, &params)
	if err != nil {
		return nil, paramsR.Error(err)
	}
	return &paramsWithMetadata[C, P]{metadata: m, params: params}, nil
}

func (m *metadata[C, P]) decodeMsg(msgR *msgp.Reader) (ParamsWithMetadata, error) {
	var params P
	err := m.paramsSchema.DecodeMsg(msgR, &params)
	if err != nil {
		return nil, err
	}
	return &paramsWithMetadata[C, P]{metadata: m, params: params}, nil

}

func (m *metadata[C, P]) ConfigBuilderWorkspace() ConfigBuilder {
	return configBuilderWorkspace[C, P]{m}
}

func (m *metadata[C, P]) ConfigBuilderNeedle(n ndl.Needle) ConfigBuilder {
	return &configBuilderNeedle[C, P]{m, n}
}

func (m *metadata[C, P]) ConfigBuilderHash(h ref.Sha256Hash) ConfigBuilder {
	return &configBuilderHash[C, P]{m, h}
}
