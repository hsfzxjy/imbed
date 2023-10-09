package lang

import (
	"fmt"

	"github.com/hsfzxjy/imbed/transform"
)

func (c *Context) parseTransSeq() (*transform.Graph, error) {
	var transforms []transform.Transform
	reader := transReader{Parser: c.parser}
	for !reader.EOF() {
		reader.Space()
		name, ok := reader.Ident()
		if !ok {
			return nil, reader.Expect("transform name")
		}
		m, ok := c.registry.Lookup(name)
		if !ok {
			return nil, reader.Error(fmt.Errorf("no transform named %q", name))
		}
		t, err := m.Parse(c.app.Config(), reader)
		if err != nil {
			return nil, err
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
	return transform.BuildGraph(transforms), nil
}
