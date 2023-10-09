package asset

import (
	"path"

	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
	"github.com/hsfzxjy/imbed/util/iter"
)

func LoadFile(filepath string) ExternalLoader {
	return func() (Asset, error) {
		a := new(asset)
		a.basename = path.Base(filepath)
		a.content = content.New(content.WithFilePath(filepath))
		return a, nil
	}
}

func FromDBModel(app core.App, model *db.AssetModel) StockLoader {
	return func(db.Context) (StockAsset, error) {
		return fromDBModel(app, model)
	}
}

func fromDBModel(app core.App, model *db.AssetModel) (*asset, error) {
	a := new(asset)
	a.model = model
	filepath := app.FilePath(model.FID.Humanize())
	a.primaryInfo = primaryInfo{
		basename: model.FID.Basename(),
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
		model, err := iter.One(it)
		if err != nil {
			return nil, err
		}
		return fromDBModel(app, model)
	}
}

type Loader interface {
	Do(db.Context) (Asset, error)
}

type ExternalLoader func() (Asset, error)

func (l ExternalLoader) Do(db.Context) (Asset, error) { return l() }

type StockLoader func(db.Context) (StockAsset, error)

func (l StockLoader) Do(c db.Context) (Asset, error) { return l(c) }
