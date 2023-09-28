package asset

import (
	"path"

	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
)

func LoadFile(filepath string) *asset {
	a := new(asset)
	a.basename = path.Base(filepath)
	a.content = content.New(content.FromFile(filepath))
	return a
}

func FromDBModel(app core.App, model *db.AssetModel) *asset {
	a := new(asset)
	a.model = model
	filepath := app.FilePath(model.FID.Humanize())
	a.primaryInfo = primaryInfo{
		basename: model.FID.Basename(),
		content:  content.New(content.FromFile(filepath)),
	}
	return a
}
