package schema

import (
	"errors"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _TopLevel[pS *S, S any] struct {
	sig sig
	*_Struct[S]
}

func (s *_TopLevel[pS, S]) New() pS { return new(S) }

func (s *_TopLevel[pS, S]) ScanFrom(r Scanner) (pS, error) {
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

func (s *_TopLevel[pS, S]) DecodeMsg(r *msgp.Reader) (pS, error) {
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

func (s *_TopLevel[pS, S]) EncodeMsg(w *msgp.Writer, source pS) error {
	err := w.Append(s.sig[:]...)
	if err != nil {
		return newError(err).SetOp("EncodeMsg").AsError()
	}
	return s.encodeMsg(w, unsafe.Pointer(source)).SetOp("EncodeMsg").AsError()
}

func (s *_TopLevel[pS, S]) Visit(v Visitor, source pS) error {
	return s.visit(v, unsafe.Pointer(source)).SetOp("Visit").AsError()
}

func _() {
	type X struct{}
	var _ Schema[*X] = &_TopLevel[*X, X]{}
}
