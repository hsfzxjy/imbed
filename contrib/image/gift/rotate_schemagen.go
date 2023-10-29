// Code generated by "github.com/hsfzxjy/imbed/schema/gen"; DO NOT EDIT.

package gift

import (
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

var rotateSchema = schema.StructFunc(func(prototype *rotate) *schema.StructBuilder[rotate] {
	return schema.Struct(prototype,
		schema.F("deg", &prototype.Deg, schema.Rat()),
		schema.F("itpl", (*string)(&prototype.Interpolation), schema.String().Default("lin")),
	).DebugName("rotate")
})

func (x *rotate) EncodeMsg(w *msgp.Writer) error {
	return rotateSchema.Build().EncodeMsg(w, x)
}
