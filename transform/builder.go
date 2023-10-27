package transform

import (
	"github.com/hsfzxjy/imbed/core"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
)

type cfgFactory = cfgf.Factory
type data = Data

type Builder struct {
	*data
	cfgFactory
}

func (b *Builder) Build(cp core.ConfigProvider) (*Transform, error) {
	cfg, err := b.cfgFactory.CreateConfig(cp)
	if err != nil {
		return nil, err
	}
	applier, err := b.buildApplier(b.params, cfg)
	if err != nil {
		return nil, err
	}
	return &Transform{
		Name:          b.name,
		Applier:       applier,
		Category:      b.category,
		Data:          b.data,
		Config:        newEncodable(b.configSchema.WrapAny(cfg)),
		ForceTerminal: false,
	}, nil
}
