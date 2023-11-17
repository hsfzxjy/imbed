package lang

import (
	"math/big"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/parser"
	"github.com/hsfzxjy/imbed/schema"
	schemascanner "github.com/hsfzxjy/imbed/schema/scanner"
)

type transScanner struct {
	*parser.Parser
	schemascanner.Void
}

func (p transScanner) Bool() (bool, pos.P, error) {
	v, pos, ok := p.Parser.Bool()
	var err error
	if !ok {
		err = p.Error(nil)
	}
	return v, pos, err
}

func (p transScanner) Rat() (*big.Rat, pos.P, error) {
	v, pos, ok := p.Parser.Rat()
	var err error
	if !ok {
		err = p.Error(nil)
	}
	return v, pos, err
}

func (p transScanner) Int64() (int64, pos.P, error) {
	v, pos, ok := p.Parser.Int64()
	var err error
	if !ok {
		err = p.Error(nil)
	}
	return v, pos, err
}

func (p transScanner) String() (string, pos.P, error) {
	v, pos, ok := p.Parser.String(",:=")
	var err error
	if !ok {
		err = p.Error(nil)
	}
	return v, pos, err
}

func (p transScanner) IterField(f func(name string, r schema.Scanner, namePos pos.P) error) error {
	const FIELD_SEP = ':'
	const KV_SEP = '='
	const BOUNDARY = ','
	for {
		p.Parser.Space()
		p.Parser.Byte(FIELD_SEP)
		p.Parser.Space()
		field_name, pos, ok := p.Parser.Ident()
		if !ok {
			break
		}
		p.Parser.Space()
		if _, ok = p.Parser.Byte(KV_SEP); !ok {
			return p.Error(nil)
		}
		pos.ExtendEnd(1)
		p.Parser.Space()
		err := f(field_name, p, pos)
		if err != nil {
			return err
		}
		if p.EOF() || p.PeekByte() == BOUNDARY {
			break
		}
	}
	return nil
}

func (p transScanner) UnnamedField() schema.Scanner {
	p.Space()
	state := p.Snapshot()
	defer p.Reset(state)
	if _, _, ok := p.Ident(); !ok {
		return p
	}
	p.Space()
	if p.PeekByte() != '=' {
		return p
	}
	return nil
}

func (p transScanner) Error(e error) error {
	return p.Parser.Error(e)
}

func _() { var _ schema.Scanner = transScanner{} }
