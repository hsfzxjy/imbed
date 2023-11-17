package schema

import (
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type _StructField struct {
	name       string
	offset     uintptr
	elemSchema genericSchema
}

func (*_StructField) _fieldSchema_stub() {}
func (s *_StructField) Name() string     { return s.name }

func (s *_StructField) scanFrom(r Scanner, target unsafe.Pointer) (pos.P, *schemaError) {
	pos, err := s.elemSchema.
		scanFrom(r, unsafe.Add(target, s.offset))
	return pos, err.AppendPath(s.name)
}

func (s *_StructField) decodeMsg(r *fastbuf.R, target unsafe.Pointer) *schemaError {
	return s.elemSchema.
		decodeMsg(r, unsafe.Add(target, s.offset)).
		AppendPath(s.name)
}

func (s *_StructField) encodeMsg(w *fastbuf.W, source unsafe.Pointer) {
	s.elemSchema.
		encodeMsg(w, unsafe.Add(source, s.offset))
}

func (s *_StructField) visit(v Visitor, source unsafe.Pointer) *schemaError {
	return s.elemSchema.
		visit(v, unsafe.Add(source, s.offset)).
		AppendPath(s.name)
}

func (s *_StructField) setDefault(target unsafe.Pointer) *schemaError {
	return s.elemSchema.
		setDefault(unsafe.Add(target, s.offset)).
		AppendPath(s.name)
}

func (s *_StructField) hasDefault() bool {
	return s.elemSchema.hasDefault()
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
