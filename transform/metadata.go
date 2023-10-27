package transform

import (
	"github.com/hsfzxjy/imbed/schema"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
)

type _constraint[C any, P ParamFor[C]] struct{}

func (_constraint[C, P]) buildApplier(params, config any) (Applier, error) {
	return params.(P).BuildTransform(config.(C))
}

type constraint interface {
	buildApplier(params, config any) (Applier, error)
}

type metadata struct {
	*Registry
	name    string
	aliases []string

	constraint

	configSchema schema.GenericSchema
	paramsSchema schema.GenericSchema

	category Category
}

func (m *metadata) Name() string {
	return m.name
}

func (m *metadata) Category() Category {
	return m.category
}

func (m *metadata) ConfigFactory(opt cfgf.Opt) cfgf.Factory {
	return opt(m.name, m.configSchema)
}
