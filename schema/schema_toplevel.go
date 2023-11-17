package schema

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type _TopLevel[S any] struct {
	sig sig
	*_Struct[S]
}

func new_Toplevel[S any](schema *_Struct[S]) *_TopLevel[S] {
	return &_TopLevel[S]{
		sig:     sigFor(schema),
		_Struct: schema,
	}
}

func (s *_TopLevel[S]) New() *S             { return new(S) }
func (s *_TopLevel[S]) NewAny() any         { return new(S) }
func (s *_TopLevel[S]) Struct() *_Struct[S] { return s._Struct }

func (s *_TopLevel[S]) ScanFrom(r Scanner) (*S, error) {
	var target = new(S)
	err := s.scanFrom(r, unsafe.Pointer(target))

	if err != nil && errors.Is(err.AsError(), ErrRequired) {
		err2 := s.setDefault(unsafe.Pointer(target))
		if err2 == nil {
			err = nil
		}
		goto RETURN
	}
RETURN:
	if err != nil {
		return nil, err.SetOp("ScanFrom").AsError()
	}
	return target, nil
}

func (s *_TopLevel[S]) ScanFromAny(r Scanner) (any, error) {
	return s.ScanFrom(r)
}

func (s *_TopLevel[S]) DecodeMsg(r *fastbuf.R) (*S, error) {
	var target = new(S)
	var sig sig
	err := r.ReadFull(sig[:])
	if err != nil {
		return nil, newError(err).SetOp("DecodeMsg").AsError()
	}
	if sig != s.sig {
		err = badSig(s.sig, sig)
		return nil, newError(err).SetOp("DecodeMsg").AsError()
	}
	err = s.decodeMsg(r, unsafe.Pointer(target)).SetOp("DecodeMsg").AsError()
	if err != nil {
		return nil, err
	}
	return target, nil
}

func (s *_TopLevel[S]) DecodeMsgAny(r *fastbuf.R) (any, error) {
	return s.DecodeMsg(r)
}

func (s *_TopLevel[S]) EncodeMsg(w *fastbuf.W, source *S) {
	w.AppendRaw(s.sig[:])
	s.encodeMsg(w, unsafe.Pointer(source))
}

func (s *_TopLevel[S]) EncodeMsgAny(w *fastbuf.W, source any) {
	v, ok := source.(*S)
	if !ok {
		panic(fmt.Errorf("expect %T, got %T", (*S)(nil), source))
	}
	s.EncodeMsg(w, v)
}

func (s *_TopLevel[S]) Visit(v Visitor, source *S) error {
	return s.visit(v, unsafe.Pointer(source)).SetOp("Visit").AsError()
}

func (s *_TopLevel[S]) VisitAny(v Visitor, source any) error {
	data, ok := source.(*S)
	if !ok {
		return fmt.Errorf("expect %T, got %T", (*S)(nil), source)
	}
	return s.Visit(v, data)
}

func (s *_TopLevel[S]) WrapAny(data any) GenericValue {
	return Wrap[*S](s, data.(*S))
}

func _() {
	type X struct{}
	var _ Schema[*X] = &_TopLevel[X]{}
}
