package schema

import "math/big"

type Visitor interface {
	AtomVisitor
	VisitList(size int) (lv ListVisitor, elem Visitor, err error)
	VisitMap(size int) (mv MapVisitor, elem Visitor, err error)
	VisitStruct(size int) (sv StructVisitor, elem Visitor, err error)
	VisitPtr(isNil, isDefault bool) (elem Visitor, err error)
}

type AtomVisitor interface {
	VisitInt64(x int64, isDefault bool) error
	VisitRat(x *big.Rat, isDefault bool) error
	VisitBool(x bool, isDefault bool) error
	VisitString(x string, isDefault bool) error
}

type ListVisitor interface {
	VisitListBegin(size int) error
	VisitListEnd(size int) error
	VisitListItemBegin(i int) error
	VisitListItemEnd(i int) error
}

type MapVisitor interface {
	VisitMapBegin(size int) error
	VisitMapEnd(size int) error
	VisitMapItemBegin(key string) error
	VisitMapItemEnd(key string) error
}

type StructVisitor interface {
	VisitStructBegin(size int) error
	VisitStructEnd(size int) error
	VisitStructFieldBegin(name string) error
	VisitStructFieldEnd(name string) error
}
