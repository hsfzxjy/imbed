package transform

import (
	"fmt"

	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/parser"
)

type Registry struct {
	metadataTable map[string]genericMetadata
}

var defaultRegistry = NewRegistry()

func DefaultRegistry() Registry { return defaultRegistry }

func NewRegistry() Registry {
	return Registry{map[string]genericMetadata{}}
}

func (r Registry) Parse(cp core.ConfigProvider, input []string) (*Graph, error) {
	reader := parser.NewReader(input)
	var transforms []Transform
	for !reader.EOF() {
		name, ok := reader.Ident()
		if !ok {
			return nil, reader.Expect("transform name")
		}
		m, ok := r.metadataTable[name]
		if !ok {
			return nil, reader.Error(fmt.Errorf("no transform named %q", name))
		}
		t, err := m.parse(cp, reader)
		if err != nil {
			return nil, reader.Error(err)
		}
		transforms = append(transforms, t)
		reader.Space()
		if ok = reader.Byte(','); !ok {
			reader.Space()
			if !reader.EOF() {
				return nil, reader.Expect(`','`)
			}
		}
	}
	return buildGraph(transforms), nil
}
