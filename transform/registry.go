package transform

type Registry struct {
	metadataTable map[string]genericMetadata
}

var defaultRegistry = NewRegistry()

func DefaultRegistry() Registry { return defaultRegistry }

func NewRegistry() Registry {
	return Registry{map[string]genericMetadata{}}
}

func (r Registry) Lookup(name string) (genericMetadata, bool) {
	m, ok := r.metadataTable[name]
	return m, ok
}
