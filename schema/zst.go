package schema

import (
	"unsafe"

	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type ZST struct{}

var ZSTSchema = Struct(&ZST{}).DebugName("ZST").buildStruct()

func (x ZST) EncodeMsg(w *fastbuf.W) {}

func ZSTSchemaAs[T ~struct{} | ~struct{ ZST }]() *_Struct[T] {
	return (*_Struct[T])(unsafe.Pointer(ZSTSchema))
}
