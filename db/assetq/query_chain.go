package assetq

import (
	"slices"

	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/asset"
	"github.com/hsfzxjy/imbed/util/iter"
	"github.com/hsfzxjy/tipe"
)

func ChainForModel(targetModel *db.AssetModel, depth int) internal.Runnable[[]*asset.AssetModel] {
	if depth < 0 {
		depth = 32768
	}
	return func(h internal.H) ([]*asset.AssetModel, error) {
		var results []*asset.AssetModel
		results = append(results, targetModel)
		model := targetModel
		for !model.OriginOID.IsZero() && depth > 0 {
			needle := ndl.RawFull(ref.AsRawString(model.OriginOID))
			it, err := ByOid(needle).RunR(h)
			if err != nil {
				return nil, err
			}
			origModel := iter.One(it)
			if origModel.IsErr() {
				return nil, origModel.UnwrapErr()
			}
			model = origModel.Unwrap()
			results = append(results, model)
			depth--
		}
		slices.Reverse(results)
		return results, nil
	}
}

func Chains(targetq Query, depth int) internal.Runnable[iter.Ator[[]*Model]] {
	return func(h internal.H) (iter.Ator[[]*Model], error) {
		it, err := targetq.RunR(h)
		if err != nil {
			return nil, err
		}
		return iter.FilterMap(it, func(m *Model) (r tipe.Result[[]*Model]) {
			return tipe.MakeR(ChainForModel(m, depth).RunR(h))
		}), nil
	}
}
