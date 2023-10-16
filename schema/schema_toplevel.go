package schema

import (
	"errors"
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type _TopLevel[S any] struct {
	sig sig
	*_Struct[S]
}

func (s *_TopLevel[S]) DecodeValue(r Reader, target *S) error {
	err := s.decodeValue(r, unsafe.Pointer(target))

	if err != nil && errors.Is(err.AsError(), ErrRequired) {
		err2 := s.setDefault(unsafe.Pointer(target))
		if err2 == nil {
			err = nil
		}
		goto RETURN
	}
RETURN:
	if err != nil {
		return err.SetOp("DecodeValue").AsError()
	}
	return nil
}

func (s *_TopLevel[S]) DecodeMsg(r *msgp.Reader, target *S) error {
	var sig sig
	_, err := r.ReadFull(sig[:])
	if err != nil {
		return newError(err).SetOp("DecodeMsg").AsError()
	}
	if sig != s.sig {
		err = badSig(s.sig, sig)
		return newError(err).SetOp("DecodeMsg").AsError()
	}
	return s.decodeMsg(r, unsafe.Pointer(target)).SetOp("DecodeMsg").AsError()
}

func (s *_TopLevel[S]) EncodeMsg(w *msgp.Writer, source *S) error {
	err := w.Append(s.sig[:]...)
	if err != nil {
		return newError(err).SetOp("EncodeMsg").AsError()
	}
	return s.encodeMsg(w, unsafe.Pointer(source)).SetOp("EncodeMsg").AsError()
}

func (s *_TopLevel[S]) Visit(v Visitor, source *S) error {
	return s.visit(v, unsafe.Pointer(source)).SetOp("Visit").AsError()
}

func _() {
	type X struct{}
	var _ Schema[X] = &_TopLevel[X]{}
}
