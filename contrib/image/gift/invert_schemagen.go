// Code generated by "github.com/hsfzxjy/imbed/schema/gen"; DO NOT EDIT.

package gift

import (
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

var invertSchema = schema.StructFunc(func(prototype *invert) *schema.StructBuilder[invert] {
	return schema.Struct(prototype).DebugName("invert")
})

func (x *invert) EncodeMsg(w *fastbuf.W) {
	invertSchema.Build().EncodeMsg(w, x)
}
