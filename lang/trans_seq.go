package lang

import (
	"fmt"

	"github.com/hsfzxjy/imbed/asset/tag"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/transform"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
)

func (c *Context) parseTransSeq(cp *configProvider) (*transform.Graph, error) {
	var stepAtomList transform.StepAtomList
	scanner := transScanner{Parser: c.parser}
	for !scanner.EOF() {
		scanner.Space()
		name, ok := scanner.Ident()
		if !ok {
			return nil, scanner.ErrorString("unknown transform")
		}
		var copt cfgf.Opt
		if c.parser.Byte('@') {
			hex, ok := c.parser.String(" :,")
			if !ok {
				return nil, scanner.ErrorString("expect config SHA")
			}
			needle, err := ndl.HexPrefix(hex)
			if err != nil {
				return nil, scanner.Error(fmt.Errorf("invalid config SHA %q: %w", hex, err))
			}
			copt = cfgf.SHANeedle(needle)
		} else {
			copt = cfgf.Workspace()
		}
		scanner.Space()
		scanner.Byte(':')
		view, err := c.registry.ScanFrom(name, scanner, copt)
		if err != nil {
			return nil, scanner.Error(err)
		}
		t, err := view.Build(cp)
		if err != nil {
			return nil, scanner.Error(err)
		}
		stepAtomList = append(stepAtomList, t)
		scanner.Space()
		spec, err := c.parseTagSpec()
		if err != nil {
			return nil, err
		}
		t.Tag = spec
		if ok = scanner.Byte(','); !ok {
			scanner.Space()
			if !scanner.EOF() {
				return nil, scanner.ErrorString(`expect ','`)
			}
		}
	}
	return transform.Schedule(c.registry, stepAtomList)
}

func (c *Context) parseTagSpec() (spec tag.Spec, err error) {
	if c.parser.PeekByte() != '+' {
		return
	}
	c.parser.Byte('+')
	spec.Kind = tag.Normal
	if c.parser.PeekByte() != '+' {
		goto PARSE_TAG
	}
	c.parser.Byte('+')
	spec.Kind = tag.Override
	if c.parser.PeekByte() != '+' {
		goto PARSE_TAG
	}
	c.parser.Byte('+')
	spec.Kind = tag.Auto
	return
PARSE_TAG:
	name, ok := c.parser.Tag()
	if !ok {
		spec.Kind = tag.None
		err = c.parser.Error(err)
		return
	}
	spec.Name = name
	return
}
