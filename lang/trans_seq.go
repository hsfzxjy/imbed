package lang

import (
	"fmt"

	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/transform"
)

func (c *Context) parseTransSeq(cp core.ConfigProvider) (*transform.Graph, error) {
	var transforms []transform.Transform
	scanner := transScanner{Parser: c.parser}
	for !scanner.EOF() {
		scanner.Space()
		name, ok := scanner.Ident()
		if !ok {
			return nil, scanner.Expect("transform name")
		}
		m, ok := c.registry.Lookup(name)
		if !ok {
			return nil, scanner.Error(fmt.Errorf("no transform named %q", name))
		}
		var cb transform.ConfigBuilder
		if c.parser.Byte('@') {
			hex, ok := c.parser.String(" :")
			if !ok {
				return nil, scanner.Expect("config hash")
			}
			needle, err := ndl.HexPrefix(hex)
			if err != nil {
				return nil, scanner.Error(fmt.Errorf("bad config hash %q: %w", hex, err))
			}
			cb = m.ConfigBuilderNeedle(needle)
		} else {
			cb = m.ConfigBuilderWorkspace()
		}
		pm, err := m.ScanParams(scanner)
		if err != nil {
			return nil, err
		}
		t, err := pm.BuildWith(cb, cp)
		if err != nil {
			return nil, err
		}
		transforms = append(transforms, t)
		scanner.Space()
		if ok = scanner.Byte(','); !ok {
			scanner.Space()
			if !scanner.EOF() {
				return nil, scanner.Expect(`','`)
			}
		}
	}
	return transform.BuildGraph(transforms), nil
}
