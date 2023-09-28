package schema

import (
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _StructField struct {
	name       string
	offset     uintptr
	elemSchema schema
}

func (*_StructField) _fieldSchema_stub() {}
func (s *_StructField) Name() string     { return s.name }

func (s *_StructField) decodeValue(r Reader, target unsafe.Pointer) *schemaError {
	return s.elemSchema.
		decodeValue(r, unsafe.Add(target, s.offset)).
		AppendPath(s.name)
}

func (s *_StructField) decodeMsg(r *msgp.Reader, target unsafe.Pointer) *schemaError {
	return s.elemSchema.
		decodeMsg(r, unsafe.Add(target, s.offset)).
		AppendPath(s.name)
}

func (s *_StructField) encodeMsg(w *msgp.Writer, source unsafe.Pointer) *schemaError {
	return s.elemSchema.
		encodeMsg(w, unsafe.Add(source, s.offset)).
		AppendPath(s.name)
}

func (s *_StructField) setDefault(target unsafe.Pointer) *schemaError {
	return s.elemSchema.
		setDefault(unsafe.Add(target, s.offset)).
		AppendPath(s.name)
}

func (s *_StructField) writeTypeInfo(w io.Writer) error {
	var buf = make([]byte, 0, len(s.name)+int(unsafe.Sizeof(uint64(0))))
	buf = append(buf, s.name...)
	offset := s.offset
	buf = binary.AppendUvarint(buf, uint64(offset))
	_, err := w.Write(buf)
	if err != nil {
		return err
	}
	return s.elemSchema.writeTypeInfo(w)
}

func _() {
	var _ fieldSchema = &_StructField{}
}
