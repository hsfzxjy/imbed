package assetq

import (
	"slices"

	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/asset"
	"github.com/hsfzxjy/imbed/util/iter"
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
			origModel, err := iter.One2(ByOid(needle).RunR(h))
			if err != nil {
				return nil, err
			}
			results = append(results, origModel)
			model = origModel
			depth--
		}
		slices.Reverse(results)
		return results, nil
	}
}

func Chains(targetq Query, depth int) internal.Runnable[core.Iterator[[]*asset.AssetModel]] {
	return func(h internal.H) (core.Iterator[[]*asset.AssetModel], error) {
		it, err := targetq.RunR(h)
		if err != nil {
			return nil, err
		}
		return iter.FilterMap(it, func(m *asset.AssetModel) ([]*asset.AssetModel, bool) {
			chain, err := ChainForModel(m, depth).RunR(h)
			if err != nil {
				return nil, false
			}
			return chain, true
		}), nil
	}
}
