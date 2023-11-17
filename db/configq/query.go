package configq

import (
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/util"
	"github.com/hsfzxjy/imbed/util/iter"
	"github.com/hsfzxjy/tipe"
)

type Provider interface {
	DBConfigBySHANeedle(needle ndl.Needle) (*db.ConfigModel, error)
	DBConfigByOID(oid ref.OID) (*db.ConfigModel, error)
}

type provider struct {
	*db.Tx
}

func (p provider) DBConfigBySHANeedle(needle ndl.Needle) (*db.ConfigModel, error) {
	cursor := db.NewCursor(p.C_SHA__OID().Cursor(), needle.Bytes())
	it := iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[*db.ConfigModel]) {
		if !needle.Match(kv.K) {
			return r.FillErr(iter.Stop)
		}
		oid, err := ref.FromRawExact[ref.OID](kv.V)
		if err != nil {
			return r.FillErr(err)
		}
		return tipe.MakeR(p.DBConfigByOID(oid))
	})
	return iter.One(it).Tuple()
}

func (p provider) DBConfigByOID(oid ref.OID) (*db.ConfigModel, error) {
	encoded := p.Tx.CONFIGS().Get(oid.Raw())
	return db.DecodeConfigModel(oid, encoded)
}

func NewProvider(tx *db.Tx) provider {
	return provider{tx}
}

func SHAByOID(oid ref.OID) db.Task[ref.Sha256] {
	return func(tx *db.Tx) (ref.Sha256, error) {
		encoded := tx.CONFIGS().Get(oid.Raw())
		return db.DecodeConfigModelSHA(encoded)
	}
}
