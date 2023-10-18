package schema

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
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

func (s *_TopLevel[S]) New() *S { return new(S) }

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

func (s *_TopLevel[S]) DecodeMsg(r *msgp.Reader) (*S, error) {
	var target = new(S)
	var sig sig
	_, err := r.ReadFull(sig[:])
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

func (s *_TopLevel[S]) EncodeMsg(w *msgp.Writer, source *S) error {
	err := w.Append(s.sig[:]...)
	if err != nil {
		return newError(err).SetOp("EncodeMsg").AsError()
	}
	return s.encodeMsg(w, unsafe.Pointer(source)).SetOp("EncodeMsg").AsError()
}

func (s *_TopLevel[S]) EncodeMsgAny(w *msgp.Writer, source any) error {
	v, ok := source.(*S)
	if !ok {
		return fmt.Errorf("expect %T, got %T", (*S)(nil), source)
	}
	return s.EncodeMsg(w, v)
}

func (s *_TopLevel[S]) Visit(v Visitor, source *S) error {
	return s.visit(v, unsafe.Pointer(source)).SetOp("Visit").AsError()
}

func _() {
	type X struct{}
	var _ Schema[*X] = &_TopLevel[X]{}
}
