package schema

import (
	"io"
	"unsafe"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type _List struct {
	// def        func() []E
	elemSchema genericSchema
	proto      listProto
	defaulter
}

func (s *_List) decodeMsg(r *fastbuf.R, target unsafe.Pointer) *schemaError {
	size, err := r.ReadInt64()
	if err != nil {
		return newError(err)
	}
	s.proto.ListInit(target, int(size))
	for i := 0; i < int(size); i++ {
		err := s.elemSchema.decodeMsg(r, s.proto.ListElem(target, i))
		if err != nil {
			return err.AppendIndex(i)
		}
	}
	return nil
}

func (s *_List) scanFrom(r Scanner, target unsafe.Pointer) (pos.P, *schemaError) {
	size, listPos, err := r.ListSize()
	if err != nil {
		return listPos, newError(err)
	}
	s.proto.ListInit(target, int(size))
	var elemPos pos.P
	err = r.IterElem(func(i int, r Scanner) error {
		p, err := s.elemSchema.scanFrom(r, s.proto.ListElem(target, i))
		if err != nil {
			err.AppendIndex(i)
		}
		elemPos = p
		listPos = listPos.Add(p)
		return err.AsError()
	})
	if err != nil {
		return elemPos, newError(err)
	}
	return listPos, nil
}

func (s *_List) encodeMsg(w *fastbuf.W, source unsafe.Pointer) {
	length := s.proto.ListLen(source)
	w.WriteInt64(int64(length))

	for i := 0; i < length; i++ {
		s.elemSchema.encodeMsg(w, s.proto.ListElem(source, i))
	}
}

func (s *_List) visit(v Visitor, source unsafe.Pointer) *schemaError {
	length := s.proto.ListLen(source)
	lv, ev, err := v.VisitList(length)
	if err != nil {
		return newError(err)
	}
	if err := lv.VisitListBegin(length); err != nil {
		return newError(err)
	}
	for i := 0; i < length; i++ {
		if err := lv.VisitListItemBegin(i); err != nil {
			return newError(err).AppendIndex(i)
		}
		if err := s.elemSchema.visit(ev, s.proto.ListElem(source, i)); err != nil {
			return err.AppendIndex(i)
		}
		if err := lv.VisitListItemEnd(i); err != nil {
			return newError(err).AppendIndex(i)
		}
	}
	if err := lv.VisitListEnd(length); err != nil {
		return newError(err)
	}
	return nil
}

// func (s *_List[E]) setDefault(target unsafe.Pointer) *schemaError {
// 	if s.def == nil {
// 		return newError(ErrRequired)
// 	}
// 	*(*[]E)(target) = s.def()
// 	return nil
// }

// func (s *_List[E]) hasDefault() bool {
// 	return s.def != nil
// }

// func (s *_List[E]) equal(a, b []E) bool {
// 	return false
// }

func (s *_List) writeTypeInfo(w io.Writer) error {
	_, err := w.Write([]byte("list"))
	if err != nil {
		return err
	}
	return s.elemSchema.writeTypeInfo(w)
}

func _() { var _ genericSchema = &_List{} }
