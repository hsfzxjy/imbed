// Code generated by "github.com/hsfzxjy/imbed/schema/gen"; DO NOT EDIT.

package imgur

import (
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

var ApplierSchema = schema.StructFunc(func(prototype *Applier) *schema.StructBuilder[Applier] {
	return schema.Struct(prototype,
		schema.F("App", &prototype.App, AppSchema),
	).DebugName("upload.imgur#Applier")
})

func (x *Applier) EncodeMsg(w *fastbuf.W) {
	ApplierSchema.Build().EncodeMsg(w, x)
}

var AppSchema = schema.StructFunc(func(prototype *App) *schema.StructBuilder[App] {
	return schema.Struct(prototype,
		schema.F("clientId", &prototype.ClientId, schema.String()),
	).DebugName("App")
})

func (x *App) EncodeMsg(w *fastbuf.W) {
	AppSchema.Build().EncodeMsg(w, x)
}

var ConfigSchema = schema.StructFunc(func(prototype *Config) *schema.StructBuilder[Config] {
	return schema.Struct(prototype,
		schema.F("apps", &prototype.Apps, schema.Map(AppSchema)),
		schema.F("default", &prototype.Default, schema.String().Default("")),
	).DebugName("upload.imgur#Config")
})

func (x *Config) EncodeMsg(w *fastbuf.W) {
	ConfigSchema.Build().EncodeMsg(w, x)
}

var ParamsSchema = schema.StructFunc(func(prototype *Params) *schema.StructBuilder[Params] {
	return schema.Struct(prototype,
		schema.F("app", &prototype.AppName, schema.String().Default("")),
	).DebugName("upload.imgur#Params")
})

func (x *Params) EncodeMsg(w *fastbuf.W) {
	ParamsSchema.Build().EncodeMsg(w, x)
}
