package parser

import (
	"github.com/hsfzxjy/imbed/schema"
	schemareader "github.com/hsfzxjy/imbed/schema/reader"
)

type Reader struct {
	*Parser
	schemareader.Void
}

func NewReader(input []string) Reader {
	return Reader{Parser: New(input)}
}

func (p Reader) Bool() (bool, error) {
	v, ok := p.Parser.Bool()
	var err error
	if !ok {
		err = p.Expect(`'true' or 'false'`)
	}
	return v, err
}

func (p Reader) Float64() (float64, error) {
	v, ok := p.Parser.Float64()
	var err error
	if !ok {
		err = p.Expect(`float literal`)
	}
	return v, err
}

func (p Reader) Int64() (int64, error) {
	v, ok := p.Parser.Int64()
	var err error
	if !ok {
		err = p.Expect(`int literal`)
	}
	return v, err
}

func (p Reader) String() (string, error) {
	v, ok := p.Parser.String()
	var err error
	if !ok {
		err = p.Expect(`string literal`)
	}
	return v, err
}

func (p Reader) IterField(f func(name string, r schema.Reader) error) error {
	const FIELD_SEP = ':'
	const KV_SEP = '='
	for {
		p.Parser.Space()
		field_name, ok := p.Parser.Ident()
		if !ok {
			break
		}
		p.Parser.Space()
		if ok = p.Parser.Byte(KV_SEP); !ok {
			return p.Expect(`'='`)
		}
		p.Parser.Space()
		err := f(field_name, p)
		if err != nil {
			return err
		}
		p.Parser.Space()
		if ok = p.Parser.Byte(FIELD_SEP); !ok {
			break
		}
	}
	return nil
}

func (p Reader) Error(e error) error {
	return p.Parser.Error(e)
}

func _() { var _ schema.Reader = Reader{} }
