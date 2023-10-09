package schema

type Visitor interface {
	AtomVisitor
	VisitList(size int) (lv ListVisitor, elem Visitor, err error)
	VisitMap(size int) (mv MapVisitor, elem Visitor, err error)
	VisitStruct(size int) (sv StructVisitor, elem Visitor, err error)
}

type AtomVisitor interface {
	VisitInt64(x int64) error
	VisitFloat64(x float64) error
	VisitBool(x bool) error
	VisitString(x string) error
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
