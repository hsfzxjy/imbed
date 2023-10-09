package asset

import (
	"sync"

	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/util"
)

type primaryInfo struct {
	basename string
	content  content.Content
}

type updatingInfo struct {
	url       string
	ext       []byte
	transform Transform
}

type asset struct {
	mu sync.RWMutex

	origin *asset
	model  *db.AssetModel

	primaryInfo
	updatingInfo
}

func (a *asset) Model() *db.AssetModel {
	return a.model
}

func (a *asset) Content() content.Content {
	return a.content
}

func (a *asset) CompareCreated(other StockAsset) int {
	return a.model.CompareCreated(other.Model())
}

func (a *asset) BaseName() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.model != nil {
		return a.model.FID.Basename()
	}

	return a.basename
}

type Asset interface {
	Content() content.Content
	BaseName() string
	save(ctx db.Context) error
	saveFile(app core.App) (util.RevertFunc, error)
}

type StockAsset interface {
	Asset
	Model() *db.AssetModel
	CompareCreated(other StockAsset) int
}
