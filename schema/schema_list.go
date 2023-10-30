package schema

import (
	"io"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _List[E any] struct {
	def        func() []E
	elemSchema schema[E]
}

func (s *_List[E]) decodeMsg(r *msgp.Reader, target unsafe.Pointer) *schemaError {
	var ret []E
	size, err := r.ReadArrayHeader()
	if err != nil {
		return newError(err)
	}
	ret = make([]E, size)
	for i := 0; i < int(size); i++ {
		err := s.elemSchema.decodeMsg(r, unsafe.Pointer(&ret[i]))
		if err != nil {
			return err.AppendIndex(i)
		}
	}
	*(*[]E)(target) = ret
	return nil
}

func (s *_List[E]) scanFrom(r Scanner, target unsafe.Pointer) *schemaError {
	var ret []E
	size, err := r.ListSize()
	if err != nil {
		return newError(err)
	}
	ret = make([]E, size)
	err = r.IterElem(func(i int, r Scanner) error {
		err := s.elemSchema.scanFrom(r, unsafe.Pointer(&ret[i]))
		if err != nil {
			err.AppendIndex(i)
		}
		return err.AsError()
	})
	if err != nil {
		return newError(err)
	}
	*(*[]E)(target) = ret
	return nil
}

func (s *_List[E]) encodeMsg(w *msgp.Writer, source unsafe.Pointer) *schemaError {
	value := *(*[]E)(source)
	err := w.WriteArrayHeader(uint32(len(value)))
	if err != nil {
		return newError(err)
	}
	for i := 0; i < len(value); i++ {
		err := s.elemSchema.encodeMsg(w, unsafe.Pointer(&value[i]))
		if err != nil {
			return err.AppendIndex(i)
		}
	}
	return nil
}

func (s *_List[E]) visit(v Visitor, source unsafe.Pointer) *schemaError {
	value := *(*[]E)(source)
	lv, ev, err := v.VisitList(len(value))
	if err != nil {
		return newError(err)
	}
	if err := lv.VisitListBegin(len(value)); err != nil {
		return newError(err)
	}
	for i := 0; i < len(value); i++ {
		if err := lv.VisitListItemBegin(i); err != nil {
			return newError(err).AppendIndex(i)
		}
		if err := s.elemSchema.visit(ev, unsafe.Pointer(&value[i])); err != nil {
			return err.AppendIndex(i)
		}
		if err := lv.VisitListItemEnd(i); err != nil {
			return newError(err).AppendIndex(i)
		}
	}
	if err := lv.VisitListEnd(len(value)); err != nil {
		return newError(err)
	}
	return nil
}

func (s *_List[E]) setDefault(target unsafe.Pointer) *schemaError {
	if s.def == nil {
		return newError(ErrRequired)
	}
	*(*[]E)(target) = s.def()
	return nil
}

func (s *_List[E]) hasDefault() bool {
	return s.def != nil
}

func (s *_List[E]) equal(a, b []E) bool {
	return false
}

func (s *_List[E]) writeTypeInfo(w io.Writer) error {
	_, err := w.Write([]byte("list"))
	if err != nil {
		return err
	}
	return s.elemSchema.writeTypeInfo(w)
}

func (s *_List[E]) _schema_stub([]E) {}

func _() { var _ schema[[]int] = &_List[int]{} }
