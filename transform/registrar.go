package transform

import (
	"github.com/hsfzxjy/imbed/schema"
)

type registrar[C any, P ParamFor[C]] struct {
	registry *Registry
	metadata *metadata[C, P]
}

func RegisterIn[C any, P ParamFor[C]](
	registry *Registry,
	name string,
	configSchema schema.Schema[C],
	paramsSchema schema.Schema[P]) registrar[C, P] {
	m := new(metadata[C, P])
	m.Registry = registry
	m.name = name
	m.configSchema = configSchema
	m.paramsSchema = paramsSchema
	r := registrar[C, P]{registry, m}
	r.registerToName(name)
	return r
}

func (r registrar[C, P]) Alias(aliases ...string) registrar[C, P] {
	r.metadata.aliases = append(r.metadata.aliases, aliases...)
	for _, alias := range aliases {
		r.registerToName(alias)
	}
	return r
}

func (r registrar[C, P]) Category(cat Category) registrar[C, P] {
	r.metadata.category = cat
	return r
}

func (r registrar[C, P]) registerToName(name string) {
	_, ok := r.registry.metadataTable[name]
	if ok {
		panic(name + " is already taken")
	}
	r.registry.metadataTable[name] = r.metadata
}
