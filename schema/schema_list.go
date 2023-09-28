package schema

import (
	"io"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _List[E any] struct {
	def        func() []E
	elemSchema schemaTyped[E]
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

func (s *_List[E]) decodeValue(r Reader, target unsafe.Pointer) *schemaError {
	var ret []E
	size, err := r.ListSize()
	if err != nil {
		return newError(err)
	}
	ret = make([]E, size)
	err = r.IterElem(func(i int, r Reader) error {
		err := s.elemSchema.decodeValue(r, unsafe.Pointer(&ret[i]))
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

func (s *_List[E]) setDefault(target unsafe.Pointer) *schemaError {
	if s.def == nil {
		return newError(ErrRequired)
	}
	*(*[]E)(target) = s.def()
	return nil
}

func (s *_List[E]) writeTypeInfo(w io.Writer) error {
	_, err := w.Write([]byte("list"))
	if err != nil {
		return err
	}
	return s.elemSchema.writeTypeInfo(w)
}

func (s *_List[E]) _schemaTyped_stub([]E) {}

func _() { var _ schemaTyped[[]int] = &_List[int]{} }
