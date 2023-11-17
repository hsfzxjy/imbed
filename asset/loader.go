package asset

import (
	"path"

	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
	"github.com/hsfzxjy/imbed/util/iter"
)

func LoadFile(filepath string, pos pos.P) ExternalLoader {
	return func() (Asset, error) {
		a := new(asset)
		a.basename = path.Base(filepath)
		a.content = content.New(content.WithFilePath(filepath), content.WithPos(pos))
		return a, nil
	}
}

func FromDBModel(app core.App, model *db.AssetModel, upstream Asset) StockLoader {
	return func(db.Context) (StockAsset, error) {
		return fromDBModel(app, model, upstream)
	}
}

func fromDBModel(app core.App, model *db.AssetModel, upstream Asset) (*asset, error) {
	a := new(asset)
	a.model = model
	if upstream != nil {
		a.origin = upstream.(*asset)
	}
	filepath := app.FilePath(model.Filename())
	a.primaryInfo = primaryInfo{
		basename: model.Filename(),
		content:  content.New(content.WithFilePath(filepath)),
	}
	return a, nil
}

func FromQ(app core.App, query assetq.Query) StockLoader {
	return func(h db.Context) (StockAsset, error) {
		it, err := query.RunR(h)
		if err != nil {
			return nil, err
		}
		model, err := iter.One(it).Tuple()
		if err != nil {
			return nil, err
		}
		return fromDBModel(app, model, nil)
	}
}

type Loader interface {
	Do(db.Context) (Asset, error)
}

type ExternalLoader func() (Asset, error)

func (l ExternalLoader) Do(db.Context) (Asset, error) { return l() }

type StockLoader func(db.Context) (StockAsset, error)

func (l StockLoader) Do(c db.Context) (Asset, error) { return l(c) }
