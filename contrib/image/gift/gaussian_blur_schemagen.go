// Code generated by "github.com/hsfzxjy/imbed/schema/gen"; DO NOT EDIT.

package gift

import (
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

var gaussianBlurSchema = schema.StructFunc(func(prototype *gaussianBlur) *schema.StructBuilder[gaussianBlur] {
	return schema.Struct(prototype,
		schema.F("s", &prototype.Sigma, schema.Rat()),
	).DebugName("gaussianBlur")
})

func (x *gaussianBlur) EncodeMsg(w *msgp.Writer) error {
	return gaussianBlurSchema.Build().EncodeMsg(w, x)
}