package schema

import (
	"io"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _Map[V any] struct {
	def         func() map[string]V
	valueSchema schemaTyped[V]
}

func (s *_Map[V]) decodeMsg(r *msgp.Reader, target unsafe.Pointer) *schemaError {
	n, err := r.ReadMapHeader()
	if err != nil {
		return newError(err)
	}
	size := int(n)
	ret := make(map[string]V, size)
	for i := 0; i < int(size); i++ {
		key, err := r.ReadString()
		if err != nil {
			return newError(err).AppendPath(key)
		}
		var value V
		err = s.valueSchema.decodeMsg(r, unsafe.Pointer(&value)).
			AppendPath(key).
			AsError()
		ret[key] = value
	}
	*(*map[string]V)(target) = ret
	return nil
}

func (s *_Map[V]) decodeValue(r Reader, target unsafe.Pointer) *schemaError {
	n, err := r.MapSize()
	if err != nil {
		return newError(err)
	}
	size := int(n)
	ret := make(map[string]V, size)
	err = r.IterKV(func(key string, r Reader) error {
		var value V
		err := s.valueSchema.decodeValue(r, unsafe.Pointer(&value))
		if err != nil {
			return err.AppendPath(key).AsError()
		}
		ret[key] = value
		return nil
	})
	if err != nil {
		return newError(err)
	}
	*(*map[string]V)(target) = ret
	return nil
}

func (s *_Map[V]) encodeMsg(w *msgp.Writer, source unsafe.Pointer) *schemaError {
	m := *(*map[string]V)(source)
	err := w.WriteMapHeader(uint32(len(m)))
	if err != nil {
		return newError(err)
	}
	for key, value := range m {
		err := w.WriteString(key)
		if err != nil {
			return newError(err).AppendPath(key)
		}
		{
			err := s.valueSchema.encodeMsg(w, unsafe.Pointer(&value))
			if err != nil {
				return err.AppendPath(key)
			}
		}
	}
	return nil
}

func (s *_Map[V]) setDefault(target unsafe.Pointer) *schemaError {
	if s.def == nil {
		return newError(ErrRequired)
	}
	*(*map[string]V)(target) = s.def()
	return nil
}

func (s *_Map[V]) writeTypeInfo(w io.Writer) error {
	_, err := w.Write([]byte("map[string]"))
	if err != nil {
		return err
	}
	return s.valueSchema.writeTypeInfo(w)
}

func (s *_Map[V]) _schemaTyped_stub(map[string]V) {}

func _() { var _ schemaTyped[map[string]int] = &_Map[int]{} }
