package assetq

import (
	"slices"

	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/util/iter"
	"github.com/hsfzxjy/tipe"
)

func ChainForModel(targetModel *db.AssetModel, depth int) db.Task[[]*db.AssetModel] {
	if depth < 0 {
		depth = 32768
	}
	return func(tx *db.Tx) ([]*db.AssetModel, error) {
		var results []*db.AssetModel
		results = append(results, targetModel)
		model := targetModel
		for !model.OriginOID.IsZero() && depth > 0 {
			origModel, err := db.New(tx, model.OriginOID.Raw())
			if err != nil {
				return nil, err
			}
			model = origModel
			results = append(results, model)
			depth--
		}
		slices.Reverse(results)
		return results, nil
	}
}

func Chains(targetq Query, depth int) db.Task[iter.Ator[[]*Model]] {
	return func(tx *db.Tx) (iter.Ator[[]*Model], error) {
		it, err := targetq.RunR(tx)
		if err != nil {
			return nil, err
		}
		return iter.FilterMap(it, func(m *Model) (r tipe.Result[[]*Model]) {
			return tipe.MakeR(ChainForModel(m, depth).RunR(tx))
		}), nil
	}
}
