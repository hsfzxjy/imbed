package asset

import (
	"sync"

	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/db"
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

	tagSpecs []tag.Spec
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
		return a.model.Basename
	}

	return a.basename
}

type Asset interface {
	Content() content.Content
	BaseName() string
	save(ctx db.Context) (StockAsset, error)
}

type StockAsset interface {
	Asset
	Model() *db.AssetModel
	CompareCreated(other StockAsset) int
}
