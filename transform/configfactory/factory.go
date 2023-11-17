package cfgf

import (
	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/configq"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type Factory interface {
	ConfigHash(ctx db.Context) ref.Sha256
	CreateConfig(cp core.ConfigProvider) (any, error)
}

type Opt func(name string, schema schema.GenericSchema) Factory

func Workspace() Opt {
	return func(name string, schema schema.GenericSchema) Factory {
		return &fromWorkspace{name, schema}
	}
}

type fromWorkspace struct {
	name string
	schema.GenericSchema
}

func (b *fromWorkspace) ConfigHash(ctx db.Context) (zero ref.Sha256) {
	return
}

func (b *fromWorkspace) CreateConfig(cp core.ConfigProvider) (result any, err error) {
	cfgR, err := cp.ProvideWorkspaceConfig(b.name)
	if err != nil {
		return
	}
	return b.GenericSchema.ScanFromAny(cfgR)
}

func Needle(needle ndl.Needle) Opt {
	return func(name string, schema schema.GenericSchema) Factory {
		return &withNeedle{schema, needle}
	}
}

type withNeedle struct {
	schema.GenericSchema
	needle ndl.Needle
}

func (b *withNeedle) ConfigHash(ctx db.Context) (zero ref.Sha256) {
	return
}

func (b *withNeedle) CreateConfig(cp core.ConfigProvider) (result any, err error) {
	buf, err := cp.ProvideStockConfig(b.needle)
	if err != nil {
		return
	}
	cfgR := fastbuf.R{Buf: buf}
	return b.GenericSchema.DecodeMsgAny(&cfgR)
}

func OID(oid ref.OID) Opt {
	return func(name string, schema schema.GenericSchema) Factory {
		return &withOID{schema, oid}
	}
}

type withOID struct {
	schema.GenericSchema
	oid ref.OID
}

func (b *withOID) ConfigHash(ctx db.Context) ref.Sha256 {
	sha, err := configq.SHAByOID(b.oid).RunR(ctx)
	if err != nil {
		return ref.Sha256{}
	}
	return sha
}

func (b *withOID) CreateConfig(cp core.ConfigProvider) (result any, err error) {
	buf, err := cp.ProvideConfigByOID(b.oid)
	if err != nil {
		return
	}
	cfgR := fastbuf.R{Buf: buf}
	return b.GenericSchema.DecodeMsgAny(&cfgR)
}
