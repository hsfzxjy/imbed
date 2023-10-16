package transform

import (
	"github.com/hsfzxjy/imbed/schema"
)

type registrar[C any, P ParamStruct[C]] struct {
	registry Registry
	*metadata[C, P]
}

func Register[C any, P ParamStruct[C]](
	name string,
	configSchema schema.Schema[C],
	paramsSchema schema.Schema[P]) registrar[C, P] {
	return RegisterIn(defaultRegistry, name, configSchema, paramsSchema)
}

func RegisterIn[C any, P ParamStruct[C]](
	registry Registry,
	name string,
	configSchema schema.Schema[C],
	paramsSchema schema.Schema[P]) registrar[C, P] {
	m := new(metadata[C, P])
	m.name = name
	m.configSchema = configSchema
	m.paramsSchema = paramsSchema
	r := registrar[C, P]{registry, m}
	r.registerToName(name)
	return r
}

func (r registrar[C, P]) Alias(aliases ...string) registrar[C, P] {
	r.aliases = append(r.aliases, aliases...)
	for _, alias := range aliases {
		r.registerToName(alias)
	}
	return r
}

func (r registrar[C, P]) Kind(kind Kind) registrar[C, P] {
	r.kind = kind
	return r
}

func (r registrar[C, P]) registerToName(name string) {
	_, ok := r.registry.metadataTable[name]
	if ok {
		panic(name + " is already taken")
	}
	r.registry.metadataTable[name] = r.metadata
}
