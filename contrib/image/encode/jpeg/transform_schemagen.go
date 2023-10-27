// Code generated by "github.com/hsfzxjy/imbed/schema/gen"; DO NOT EDIT.

package jpeg

import (
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

var ApplierSchema = schema.StructFunc(func(prototype *Applier) *schema.StructBuilder[Applier] {
	return schema.Struct(prototype,
		schema.F("Quality", &prototype.Quality, schema.Int()),
	).DebugName("Applier")
})

func (x *Applier) EncodeMsg(w *msgp.Writer) error {
	return ApplierSchema.Build().EncodeMsg(w, x)
}

var ConfigSchema = schema.StructFunc(func(prototype *Config) *schema.StructBuilder[Config] {
	return schema.Struct(prototype,
		schema.F("default_quality", &prototype.DefaultQuality, schema.Int().Default(75)),
	).DebugName("image.encode.jpeg#Config")
})

func (x *Config) EncodeMsg(w *msgp.Writer) error {
	return ConfigSchema.Build().EncodeMsg(w, x)
}

var ParamsSchema = schema.StructFunc(func(prototype *Params) *schema.StructBuilder[Params] {
	return schema.Struct(prototype,
		schema.F("q", &prototype.Quality, schema.Int().Default(-1)),
	).DebugName("image.encode.jpeg#Params")
})

func (x *Params) EncodeMsg(w *msgp.Writer) error {
	return ParamsSchema.Build().EncodeMsg(w, x)
}
