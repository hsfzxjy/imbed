package configq

import (
	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/bucketnames"
	"github.com/hsfzxjy/imbed/util"
	"github.com/hsfzxjy/imbed/util/iter"
	"github.com/hsfzxjy/tipe"
)

type provider struct {
	h internal.H
}

func (p provider) ProvideStockConfig(needle ndl.Needle) ([]byte, error) {
	cursor, err := p.h.Bucket(bucketnames.CONFIGS).Cursor(needle.Bytes())
	if err != nil {
		return nil, err
	}
	it := iter.FilterMap(cursor, func(kv util.KV) (r tipe.Result[[]byte]) {
		if needle.Match(kv.K) {
			return tipe.Ok(kv.V)
		} else {
			return r.FillErr(iter.Stop)
		}
	})
	return iter.One(it).Tuple()
}

type configProvider struct {
	provider
	core.App
}

func NewProvider(ctx db.Context, app core.App) core.ConfigProvider {
	return configProvider{
		provider: provider{
			h: ctx.(internal.H),
		},
		App: app,
	}
}
