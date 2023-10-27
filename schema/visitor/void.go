package schemavisitor

import (
	"math/big"

	"github.com/hsfzxjy/imbed/schema"
)

type Void struct{}

// VisitStructBegin implements schema.StructVisitor.
func (Void) VisitStructBegin(size int) error {
	return nil
}

// VisitStructEnd implements schema.StructVisitor.
func (Void) VisitStructEnd(size int) error {
	return nil
}

// VisitStructFieldBegin implements schema.StructVisitor.
func (Void) VisitStructFieldBegin(name string) error {
	return nil
}

// VisitStructFieldEnd implements schema.StructVisitor.
func (Void) VisitStructFieldEnd(name string) error {
	return nil
}

// VisitMapBegin implements schema.MapVisitor.
func (Void) VisitMapBegin(size int) error {
	return nil
}

// VisitMapEnd implements schema.MapVisitor.
func (Void) VisitMapEnd(size int) error {
	return nil
}

// VisitMapItemBegin implements schema.MapVisitor.
func (Void) VisitMapItemBegin(key string) error {
	return nil
}

// VisitMapItemEnd implements schema.MapVisitor.
func (Void) VisitMapItemEnd(key string) error {
	return nil
}

// VisitListBegin implements schema.ListVisitor.
func (Void) VisitListBegin(size int) error {
	return nil
}

// VisitListEnd implements schema.ListVisitor.
func (Void) VisitListEnd(size int) error {
	return nil
}

// VisitListItemBegin implements schema.ListVisitor.
func (Void) VisitListItemBegin(i int) error {
	return nil
}

// VisitListItemEnd implements schema.ListVisitor.
func (Void) VisitListItemEnd(i int) error {
	return nil
}

// VisitBool implements schema.Visitor.
func (Void) VisitBool(x bool, isDefault bool) error {
	return nil
}

// VisitRat implements schema.Visitor.
func (Void) VisitRat(x *big.Rat, isDefault bool) error {
	return nil
}

// VisitInt64 implements schema.Visitor.
func (Void) VisitInt64(x int64, isDefault bool) error {
	return nil
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
func (Void) VisitString(x string, isDefault bool) error {
	return nil
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
