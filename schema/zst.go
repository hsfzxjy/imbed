package schema

import (
	"unsafe"

	"github.com/tinylib/msgp/msgp"
)

type ZST struct{}

var ZSTSchema = Struct(&ZST{}).DebugName("ZST").buildStruct()

func (x ZST) EncodeMsg(w *msgp.Writer) error {
	return nil
}

func ZSTSchemaAs[T ~struct{} | ~struct{ ZST }]() *_Struct[T] {
	return (*_Struct[T])(unsafe.Pointer(ZSTSchema))
}
