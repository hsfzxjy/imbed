package lang

import (
	"fmt"

	"github.com/hsfzxjy/imbed/asset/tag"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/transform"
	cfgf "github.com/hsfzxjy/imbed/transform/configfactory"
)

func (c *Context) parseTransSeq(cp *configProvider) (*transform.Graph, error) {
	var stepAtomList transform.StepAtomList
	scanner := transScanner{Parser: c.parser}
	for !scanner.EOF() {
		scanner.Space()
		name, pos, ok := scanner.Ident()
		if !ok {
			return nil, pos.WrapErrorString("unknown transform")
		}
		var copt cfgf.Opt
		if _, ok := c.parser.Byte('@'); ok {
			hex, pos, ok := c.parser.String(" :,")
			if !ok {
				return nil, pos.WrapErrorString("expect config SHA")
			}
			needle, err := ndl.HexPrefix(hex, pos)
			if err != nil {
				return nil, pos.WrapError(fmt.Errorf("invalid config SHA %q: %w", hex, err))
			}
			copt = cfgf.SHANeedle(needle)
		} else {
			copt = cfgf.Workspace()
		}
		scanner.Space()
		scanner.Byte(':')
		view, err := c.registry.ScanFrom(name, scanner, copt)
		if err != nil {
			return nil, err
		}
		t, err := view.Build(cp)
		if err != nil {
			return nil, err
		}
		stepAtomList = append(stepAtomList, t)
		scanner.Space()
		spec, pos, err := c.parseTagSpec()
		if err != nil {
			return nil, pos.WrapError(err)
		}
		t.Tag = spec
		if pos, ok = scanner.Byte(','); !ok {
			scanner.Space()
			if !scanner.EOF() {
				return nil, pos.WrapErrorString(`expect ','`)
			}
		}
	}
	return transform.Schedule(c.registry, stepAtomList)
}

func (c *Context) parseTagSpec() (spec tag.Spec, pos pos.P, err error) {
	pos = c.parser.Pos(0)
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
	name, pos2, ok := c.parser.Tag()
	pos = pos.Add(pos2)
	if !ok {
		spec.Kind = tag.None
		err = c.parser.Error(nil)
		return
	}
	spec.Name = name
	return
}
