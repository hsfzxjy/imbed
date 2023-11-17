package cfgf

import (
	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/configq"
)

type configProvider interface {
	configq.Provider
	core.WorkspaceConfigProvider
}

type Opt interface {
	ConfigHash(ctx db.Context) ref.Sha256
	FromDB() bool
	QueryModel(cp configProvider) (*db.ConfigModel, error)
}

func Workspace() Opt {
	return fromWorkspace{}
}

type fromWorkspace struct{}

func (b fromWorkspace) ConfigHash(ctx db.Context) (zero ref.Sha256) {
	return
}

func (b fromWorkspace) FromDB() bool {
	return false
}

func (b fromWorkspace) QueryModel(cp configProvider) (*db.ConfigModel, error) {
	panic("not implemented")
}

func SHANeedle(needle ndl.Needle) Opt {
	return withNeedle{needle}
}

type withNeedle struct {
	needle ndl.Needle
}

func (b withNeedle) ConfigHash(ctx db.Context) (zero ref.Sha256) {
	return
}

func (b withNeedle) FromDB() bool {
	return true
}

func (b withNeedle) QueryModel(cp configProvider) (*db.ConfigModel, error) {
	return cp.DBConfigBySHANeedle(b.needle)
}

func OID(oid ref.OID) Opt {
	return withOID{oid}
}

type withOID struct {
	oid ref.OID
}

func (b withOID) ConfigHash(ctx db.Context) ref.Sha256 {
	if ctx == nil {
		return ref.Sha256{}
	}
	sha, err := configq.SHAByOID(b.oid).RunR(ctx)
	if err != nil {
		return ref.Sha256{}
	}
	return sha
}

func (b withOID) FromDB() bool {
	return true
}

func (b withOID) QueryModel(cp configProvider) (*db.ConfigModel, error) {
	return cp.DBConfigByOID(b.oid)
}
