package lang

import (
	"fmt"

	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/transform"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
)

func (c *Context) parseTransSeq(cp core.ConfigProvider) (*transform.Graph, error) {
	var transforms []*transform.Transform
	scanner := transScanner{Parser: c.parser}
	for !scanner.EOF() {
		scanner.Space()
		name, ok := scanner.Ident()
		if !ok {
			return nil, scanner.ErrorString("unknown transform")
		}
		data, err := c.registry.ScanFrom(name, scanner)
		if err != nil {
			return nil, err
		}
		var cb cfgf.Factory
		if c.parser.Byte('@') {
			hex, ok := c.parser.String(" :")
			if !ok {
				return nil, scanner.ErrorString("expect config SHA")
			}
			needle, err := ndl.HexPrefix(hex)
			if err != nil {
				return nil, scanner.Error(fmt.Errorf("invalid config SHA %q: %w", hex, err))
			}
			cb = data.ConfigFactory(cfgf.Needle(needle))
		} else {
			cb = data.ConfigFactory(cfgf.Workspace())
		}
		t, err := data.AsBuilder(cb).Build(cp)
		if err != nil {
			return nil, err
		}
		transforms = append(transforms, t)
		scanner.Space()
		if ok = scanner.Byte(','); !ok {
			scanner.Space()
			if !scanner.EOF() {
				return nil, scanner.ErrorString(`expect ','`)
			}
		}
	}
	return transform.Schedule(c.registry, transforms)
}
