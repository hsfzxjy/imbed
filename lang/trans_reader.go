package lang

import (
	"math/big"

	"github.com/hsfzxjy/imbed/parser"
	"github.com/hsfzxjy/imbed/schema"
	schemareader "github.com/hsfzxjy/imbed/schema/reader"
)

type transReader struct {
	*parser.Parser
	schemareader.Void
}

func (p transReader) Bool() (bool, error) {
	v, ok := p.Parser.Bool()
	var err error
	if !ok {
		err = p.Expect(`'true' or 'false'`)
	}
	return v, err
}

func (p transReader) Rat() (*big.Rat, error) {
	v, ok := p.Parser.Rat()
	var err error
	if !ok {
		err = p.Expect(`Rat literal`)
	}
	return v, err
}

func (p transReader) Int64() (int64, error) {
	v, ok := p.Parser.Int64()
	var err error
	if !ok {
		err = p.Expect(`int literal`)
	}
	return v, err
}

func (p transReader) String() (string, error) {
	v, ok := p.Parser.String(",:=")
	var err error
	if !ok {
		err = p.Expect(`string literal`)
	}
	return v, err
}

func (p transReader) IterField(f func(name string, r schema.Reader) error) error {
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
			return p.Expect(`'='`)
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

func (p transReader) Error(e error) error {
	return p.Parser.Error(e)
}

func _() { var _ schema.Reader = transReader{} }
