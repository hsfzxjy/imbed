package configq

import (
	"fmt"

	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/util"
	"github.com/hsfzxjy/imbed/util/iter"
	"github.com/hsfzxjy/tipe"
)

type provider struct {
	*db.Tx
}

func (p provider) ProvideStockConfig(needle ndl.Needle) ([]byte, error) {
	cursor := db.NewCursor(p.C_SHA__OID().Cursor(), needle.Bytes())
	it := iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[[]byte]) {
		if !needle.Match(kv.K) {
			return r.FillErr(iter.Stop)
		}
		data := p.Tx.CONFIGS().Get(kv.V)
		_, err, rest := ref.FromRaw[ref.Sha256](data)
		return tipe.MakeR(rest, err)
	})
	return iter.One(it).Tuple()
}

func (p provider) ProvideConfigByOID(oid ref.OID) ([]byte, error) {
	return nil, nil
}

type configProvider struct {
	provider
	core.App
}

func NewProvider(tx *db.Tx, app core.App) core.ConfigProvider {
	return configProvider{
		provider: provider{tx},
		App:      app,
	}
}

func SHAByOID(oid ref.OID) db.Task[ref.Sha256] {
	return func(tx *db.Tx) (ref.Sha256, error) {
		data := tx.CONFIGS().Get(oid.Raw())
		if data == nil {
			return ref.Sha256{}, fmt.Errorf("config not found with oid=%s", oid.FmtString())
		}
		sha, err, _ := ref.FromRaw[ref.Sha256](data)
		return sha, err
	}
}
