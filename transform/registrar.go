package transform

import (
	"github.com/hsfzxjy/imbed/schema"
)

type registrar[C any, P IParam[C, A], A IApplier] struct {
	registry Registry
	*metadata[C, P, A]
}

func RegisterIn[C any, P IParam[C, A], A IApplier](
	registry Registry,
	name string,
	configSchema schema.Schema[C],
	paramsSchema schema.Schema[P],
	applierSchema schema.Schema[A]) registrar[C, P, A] {
	m := new(metadata[C, P, A])
	m.name = name
	m.configSchema = configSchema
	m.paramsSchema = paramsSchema
	m.applierSchema = applierSchema
	r := registrar[C, P, A]{registry, m}
	r.registerToName(name)
	return r
}

func (r registrar[C, P, A]) Alias(aliases ...string) registrar[C, P, A] {
	r.aliases = append(r.aliases, aliases...)
	for _, alias := range aliases {
		r.registerToName(alias)
	}
	return r
}

func (r registrar[C, P, A]) Kind(kind Kind) registrar[C, P, A] {
	r.kind = kind
	return r
}

func (r registrar[C, P, A]) registerToName(name string) {
	_, ok := r.registry.metadataTable[name]
	if ok {
		panic(name + " is already taken")
	}
	r.registry.metadataTable[name] = r.metadata
}
