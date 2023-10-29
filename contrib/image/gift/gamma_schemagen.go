// Code generated by "github.com/hsfzxjy/imbed/schema/gen"; DO NOT EDIT.

package gift

import (
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

var gammaSchema = schema.StructFunc(func(prototype *gamma) *schema.StructBuilder[gamma] {
	return schema.Struct(prototype,
		schema.F("g", &prototype.Gamma, schema.Rat()),
	).DebugName("gamma")
})

func (x *gamma) EncodeMsg(w *msgp.Writer) error {
	return gammaSchema.Build().EncodeMsg(w, x)
}
