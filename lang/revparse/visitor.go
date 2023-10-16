package revparse

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/hsfzxjy/imbed/schema"
	schemavisitor "github.com/hsfzxjy/imbed/schema/visitor"
)

type visitor struct {
	b *strings.Builder
	schemavisitor.Void
}

func NewVisitor(b *strings.Builder) visitor {
	return visitor{b, schemavisitor.Void{}}
}

func (v visitor) VisitStructFieldBegin(name string) error {
	v.b.WriteString(":")
	v.b.WriteString(name)
	v.b.WriteByte('=')
	return nil
}

func (v visitor) VisitBool(x bool) error {
	var s string
	if x {
		s = "true"
	} else {
		s = "false"
	}
	v.b.WriteString(s)
	return nil
}

func (v visitor) VisitInt64(x int64) error {
	s := strconv.FormatInt(x, 10)
	v.b.WriteString(s)
	return nil
}

func (v visitor) VisitRat(x *big.Rat) error {
	v.b.WriteString(fmtRat(x))
	return nil
}

func (v visitor) VisitString(x string) error {
	quoteString(v.b, x, '[')
	return nil
}

func (v visitor) VisitStruct(size int) (sv schema.StructVisitor, elem schema.Visitor, err error) {
	return v, v, nil
}
