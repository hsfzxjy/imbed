package db

import (
	"github.com/hsfzxjy/imbed/db/internal"
	"github.com/hsfzxjy/imbed/db/internal/asset"
	"github.com/hsfzxjy/imbed/db/internal/config"
)

type AssetModel = asset.AssetModel
type AssetTemplate = asset.AssetTemplate
type TransSeq = asset.TransSeq
type ConfigModel = config.ConfigModel

type Context = internal.Context
type Service = internal.Service

type App = internal.App

func Open(app App) (Service, error) {
	return internal.Open(app)
}
