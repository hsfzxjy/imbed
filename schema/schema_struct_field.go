package schema

import (
	"io"
	"unsafe"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type goStructFieldProto struct {
	offset uintptr
}

func (p goStructFieldProto) FieldPtr(target unsafe.Pointer) unsafe.Pointer {
	return unsafe.Add(target, p.offset)
}

type _StructField struct {
	name string
	// offset     uintptr
	proto      structFieldProto
	elemSchema genericSchema
}

// func (s *_StructField) Name() string { return s.name }

func (s *_StructField) scanFrom(r Scanner, target unsafe.Pointer) (pos.P, *schemaError) {
	pos, err := s.elemSchema.
		scanFrom(r, s.proto.FieldPtr(target))
	return pos, err.AppendPath(s.name)
}

func (s *_StructField) decodeMsg(r *fastbuf.R, target unsafe.Pointer) *schemaError {
	return s.elemSchema.
		decodeMsg(r, s.proto.FieldPtr(target)).
		AppendPath(s.name)
}

func (s *_StructField) encodeMsg(w *fastbuf.W, source unsafe.Pointer) {
	s.elemSchema.
		encodeMsg(w, s.proto.FieldPtr(source))
}

func (s *_StructField) visit(v Visitor, source unsafe.Pointer) *schemaError {
	return s.elemSchema.
		visit(v, s.proto.FieldPtr(source)).
		AppendPath(s.name)
}

func (s *_StructField) setDefault(target unsafe.Pointer) *schemaError {
	return s.elemSchema.
		setDefault(s.proto.FieldPtr(target)).
		AppendPath(s.name)
}

func (s *_StructField) hasDefault() bool {
	return s.elemSchema.hasDefault()
}

func (s *_StructField) writeTypeInfo(w io.Writer) error {
	name := s.name
	var buf = make([]byte, 0, len(name)+int(unsafe.Sizeof(uint64(0))))
	buf = append(buf, name...)
	_, err := w.Write(buf)
	if err != nil {
		return err
	}
	return s.elemSchema.writeTypeInfo(w)
}
