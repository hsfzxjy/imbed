package schemavisitor

import (
	"math/big"

	"github.com/hsfzxjy/imbed/schema"
)

type Void struct{}

// VisitStructBegin implements schema.StructVisitor.
func (Void) VisitStructBegin(size int) error {
	panic("unimplemented")
}

// VisitStructEnd implements schema.StructVisitor.
func (Void) VisitStructEnd(size int) error {
	panic("unimplemented")
}

// VisitStructFieldBegin implements schema.StructVisitor.
func (Void) VisitStructFieldBegin(name string) error {
	panic("unimplemented")
}

// VisitStructFieldEnd implements schema.StructVisitor.
func (Void) VisitStructFieldEnd(name string) error {
	panic("unimplemented")
}

// VisitMapBegin implements schema.MapVisitor.
func (Void) VisitMapBegin(size int) error {
	panic("unimplemented")
}

// VisitMapEnd implements schema.MapVisitor.
func (Void) VisitMapEnd(size int) error {
	panic("unimplemented")
}

// VisitMapItemBegin implements schema.MapVisitor.
func (Void) VisitMapItemBegin(key string) error {
	panic("unimplemented")
}

// VisitMapItemEnd implements schema.MapVisitor.
func (Void) VisitMapItemEnd(key string) error {
	panic("unimplemented")
}

// VisitListBegin implements schema.ListVisitor.
func (Void) VisitListBegin(size int) error {
	panic("unimplemented")
}

// VisitListEnd implements schema.ListVisitor.
func (Void) VisitListEnd(size int) error {
	panic("unimplemented")
}

// VisitListItemBegin implements schema.ListVisitor.
func (Void) VisitListItemBegin(i int) error {
	panic("unimplemented")
}

// VisitListItemEnd implements schema.ListVisitor.
func (Void) VisitListItemEnd(i int) error {
	panic("unimplemented")
}

// VisitBool implements schema.Visitor.
func (Void) VisitBool(x bool) error {
	panic("unimplemented")
}

// VisitRat implements schema.Visitor.
func (Void) VisitRat(x *big.Rat) error {
	panic("unimplemented")
}

// VisitInt64 implements schema.Visitor.
func (Void) VisitInt64(x int64) error {
	panic("unimplemented")
}

// VisitList implements schema.Visitor.
func (Void) VisitList(size int) (lv schema.ListVisitor, elem schema.Visitor, err error) {
	panic("unimplemented")
}

// VisitMap implements schema.Visitor.
func (Void) VisitMap(size int) (mv schema.MapVisitor, elem schema.Visitor, err error) {
	panic("unimplemented")
}

// VisitString implements schema.Visitor.
func (Void) VisitString(x string) error {
	panic("unimplemented")
}

// VisitStruct implements schema.Visitor.
func (Void) VisitStruct(size int) (sv schema.StructVisitor, elem schema.Visitor, err error) {
	panic("unimplemented")
}

func _() {
	var _ schema.Visitor = Void{}
	var _ schema.ListVisitor = Void{}
	var _ schema.MapVisitor = Void{}
	var _ schema.StructVisitor = Void{}
}
