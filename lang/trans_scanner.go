package lang

import (
	"math/big"

	"github.com/hsfzxjy/imbed/parser"
	"github.com/hsfzxjy/imbed/schema"
	schemascanner "github.com/hsfzxjy/imbed/schema/scanner"
)

type transScanner struct {
	*parser.Parser
	schemascanner.Void
}

func (p transScanner) Bool() (bool, error) {
	v, ok := p.Parser.Bool()
	var err error
	if !ok {
		err = p.Error(nil)
	}
	return v, err
}

func (p transScanner) Rat() (*big.Rat, error) {
	v, ok := p.Parser.Rat()
	var err error
	if !ok {
		err = p.Error(nil)
	}
	return v, err
}

func (p transScanner) Int64() (int64, error) {
	v, ok := p.Parser.Int64()
	var err error
	if !ok {
		err = p.Error(nil)
	}
	return v, err
}

func (p transScanner) String() (string, error) {
	v, ok := p.Parser.String(",:=")
	var err error
	if !ok {
		err = p.Error(nil)
	}
	return v, err
}

func (p transScanner) IterField(f func(name string, r schema.Scanner) error) error {
	const FIELD_SEP = ':'
	const KV_SEP = '='
	const BOUNDARY = ','
	for {
		p.Parser.Space()
		p.Parser.Byte(FIELD_SEP)
		p.Parser.Space()
		field_name, ok := p.Parser.Ident()
		if !ok {
			break
		}
		p.Parser.Space()
		if ok = p.Parser.Byte(KV_SEP); !ok {
			return p.Error(nil)
		}
		p.Parser.Space()
		err := f(field_name, p)
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
	if _, ok := p.Ident(); !ok {
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
