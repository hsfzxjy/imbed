package schema

import (
	"io"
	"slices"
	"unsafe"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type goMapProto[V any] struct{}

func (p goMapProto[V]) MapInit(target unsafe.Pointer, size int) {
	*(*map[string]V)(target) = make(map[string]V, size)
}

func (p goMapProto[V]) MapLen(target unsafe.Pointer) int {
	return len(*(*map[string]V)(target))
}

func (p goMapProto[V]) MapNewValue() unsafe.Pointer {
	return unsafe.Pointer(new(V))
}

func (p goMapProto[V]) MapSetValue(target unsafe.Pointer, key string, value unsafe.Pointer) {
	(*(*map[string]V)(target))[key] = *(*V)(value)
}

func (p goMapProto[V]) MapIter(target unsafe.Pointer, f func(key string, value unsafe.Pointer) bool) {
	for k, v := range *(*map[string]V)(target) {
		if !f(k, unsafe.Pointer(&v)) {
			break
		}
	}
}

func (p goMapProto[V]) MapIterOrdered(target unsafe.Pointer, f func(key string, value unsafe.Pointer) bool) {
	m := *(*map[string]V)(target)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, k := range keys {
		v := m[k]
		if !f(k, unsafe.Pointer(&v)) {
			break
		}
	}
}

type goMapDefaulter[V any] struct {
	def func() map[string]V
}

func (d goMapDefaulter[V]) hasDefault() bool {
	return d.def != nil
}

func (d goMapDefaulter[V]) setDefault(target unsafe.Pointer) *schemaError {
	if d.def == nil {
		return newError(ErrRequired)
	}
	*(*map[string]V)(target) = d.def()
	return nil
}

type _Map struct {
	valueSchema genericSchema
	proto       mapProto
	defaulter
}

func (s *_Map) decodeMsg(r *fastbuf.R, target unsafe.Pointer) *schemaError {
	n, err := r.ReadInt64()
	if err != nil {
		return newError(err)
	}
	size := int(n)
	s.proto.MapInit(target, size)
	for i := 0; i < int(size); i++ {
		key, err := r.ReadString()
		if err != nil {
			return newError(err).AppendPath(key)
		}
		var value = s.proto.MapNewValue()
		{
			err := s.valueSchema.decodeMsg(r, value).
				AppendPath(key)
			if err != nil {
				return err
			}
		}
		s.proto.MapSetValue(target, key, value)
	}
	return nil
}

func (s *_Map) scanFrom(r Scanner, target unsafe.Pointer) (pos.P, *schemaError) {
	n, mapPos, err := r.MapSize()
	if err != nil {
		return mapPos, newError(err)
	}
	size := int(n)
	s.proto.MapInit(target, size)
	var elemPos pos.P
	err = r.IterKV(func(key string, r Scanner) error {
		var value = s.proto.MapNewValue()
		p, err := s.valueSchema.scanFrom(r, value)
		if err != nil {
			return err.AppendPath(key).AsError()
		}
		elemPos = p
		mapPos = mapPos.Add(p)
		s.proto.MapSetValue(target, key, value)
		return nil
	})
	if err != nil {
		return elemPos, newError(err)
	}
	return mapPos, nil
}

func (s *_Map) encodeMsg(w *fastbuf.W, source unsafe.Pointer) {
	length := s.proto.MapLen(source)
	w.WriteInt64(int64(length))
	s.proto.MapIterOrdered(source, func(key string, value unsafe.Pointer) bool {
		w.WriteString(key)
		s.valueSchema.encodeMsg(w, value)
		return true
	})
}

func (s *_Map) visit(v Visitor, source unsafe.Pointer) *schemaError {
	length := s.proto.MapLen(source)
	mv, ev, err := v.VisitMap(length)
	if err != nil {
		return newError(err)
	}
	if err := mv.VisitMapBegin(length); err != nil {
		return newError(err)
	}
	var serr *schemaError
	s.proto.MapIter(source, func(key string, value unsafe.Pointer) bool {
		err := mv.VisitMapItemBegin(key)
		if err != nil {
			serr = newError(err).AppendPath(key)
			return false
		}

		serr = s.valueSchema.visit(ev, unsafe.Pointer(&v))
		if serr != nil {
			serr = serr.AppendPath(key)
			return false
		}

		err = mv.VisitMapItemEnd(key)
		if err != nil {
			serr = newError(err).AppendPath(key)
			return false
		}

		return true
	})
	if serr != nil {
		return serr
	}
	if err := mv.VisitMapEnd(length); err != nil {
		return newError(err)
	}
	return nil
}

func (s *_Map) writeTypeInfo(w io.Writer) error {
	_, err := w.Write([]byte("map[string]"))
	if err != nil {
		return err
	}
	return s.valueSchema.writeTypeInfo(w)
}

func _() { var _ genericSchema = &_Map{} }
