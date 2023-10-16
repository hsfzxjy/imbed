package assetq

import (
	"slices"

	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/asset"
	"github.com/hsfzxjy/imbed/util/iter"
)

func Chain(targetModel *asset.AssetModel, depth int) internal.Runnable[[]*asset.AssetModel] {
	if depth < 0 {
		depth = 32768
	}
	return func(h internal.H) ([]*asset.AssetModel, error) {
		var results []*asset.AssetModel
		results = append(results, targetModel)
		model := targetModel
		for !model.OriginOID.IsZero() && depth > 0 {
			needle := core.StringFull(ref.AsRawString(model.OriginOID))
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
