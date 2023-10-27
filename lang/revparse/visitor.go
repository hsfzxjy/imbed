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
	currentField string
}

func NewVisitor(b *strings.Builder) *visitor {
	return &visitor{b: b}
}

func (v *visitor) writeField() {
	v.b.WriteString(":")
	v.b.WriteString(v.currentField)
	v.b.WriteByte('=')
}

func (v *visitor) VisitStructFieldBegin(name string) error {
	v.currentField = name
	return nil
}

func (v *visitor) VisitBool(x bool, isDefault bool) error {
	var s string
	if x {
		s = "true"
	} else {
		s = "false"
	}
	v.b.WriteString(s)
	return nil
}

func (v *visitor) VisitInt64(x int64, isDefault bool) error {
	if isDefault {
		return nil
	}
	v.writeField()
	s := strconv.FormatInt(x, 10)
	v.b.WriteString(s)
	return nil
}

func (v *visitor) VisitRat(x *big.Rat, isDefault bool) error {
	if isDefault {
		return nil
	}
	v.writeField()
	v.b.WriteString(fmtRat(x))
	return nil
}

func (v *visitor) VisitString(x string, isDefault bool) error {
	if isDefault {
		return nil
	}
	v.writeField()
	quoteString(v.b, x, '[')
	return nil
}

func (v *visitor) VisitStruct(size int) (sv schema.StructVisitor, elem schema.Visitor, err error) {
	return v, v, nil
}
