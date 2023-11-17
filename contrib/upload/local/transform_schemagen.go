// Code generated by "github.com/hsfzxjy/imbed/schema/gen"; DO NOT EDIT.

package local

import (
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

var ApplierSchema = schema.StructFunc(func(prototype *Applier) *schema.StructBuilder[Applier] {
	return schema.Struct(prototype,
		schema.F("Path", &prototype.Path, schema.String()),
	).DebugName("upload.local#Applier")
})

func (x *Applier) EncodeMsg(w *fastbuf.W) {
	ApplierSchema.Build().EncodeMsg(w, x)
}

var ConfigSchema = schema.StructFunc(func(prototype *Config) *schema.StructBuilder[Config] {
	return schema.Struct(prototype,
		schema.F("path", &prototype.Path, schema.String().Default("")),
	).DebugName("upload.local#Config")
})

func (x *Config) EncodeMsg(w *fastbuf.W) {
	ConfigSchema.Build().EncodeMsg(w, x)
}

var ParamsSchema = schema.StructFunc(func(prototype *Params) *schema.StructBuilder[Params] {
	return schema.Struct(prototype,
		schema.F("path", &prototype.Path, schema.String().Default("")),
	).DebugName("upload.local#Params")
})

func (x *Params) EncodeMsg(w *fastbuf.W) {
	ParamsSchema.Build().EncodeMsg(w, x)
}
