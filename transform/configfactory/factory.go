package cfgf

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

type Factory interface {
	ConfigHash() ref.Sha256Hash
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

func (b *fromWorkspace) ConfigHash() (zero ref.Sha256Hash) {
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

func (b *withNeedle) ConfigHash() (zero ref.Sha256Hash) {
	return
}

func (b *withNeedle) CreateConfig(cp core.ConfigProvider) (result any, err error) {
	buf, err := cp.ProvideStockConfig(b.needle)
	if err != nil {
		return
	}
	cfgR := msgp.NewReader(bytes.NewReader(buf))
	return b.GenericSchema.DecodeMsgAny(cfgR)
}

func Hash(hash ref.Sha256Hash) Opt {
	return func(name string, schema schema.GenericSchema) Factory {
		return &withHash{schema, hash}
	}
}

type withHash struct {
	schema.GenericSchema
	hash ref.Sha256Hash
}

func (b *withHash) ConfigHash() ref.Sha256Hash {
	return b.hash
}

func (b *withHash) CreateConfig(cp core.ConfigProvider) (result any, err error) {
	buf, err := cp.ProvideStockConfig(ndl.RawFull(ref.AsRawString(b.hash)))
	if err != nil {
		return
	}
	cfgR := msgp.NewReader(bytes.NewReader(buf))
	return b.GenericSchema.DecodeMsgAny(cfgR)
}
