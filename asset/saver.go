package asset

import (
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/db"
)

func SaveAll(ctx db.Context, app core.App, assets []Asset) ([]StockAsset, error) {
	var result = make([]StockAsset, 0, len(assets))
	for _, a := range assets {
		if sa, err := a.save(ctx); err != nil {
			return nil, err
		} else {
			result = append(result, sa)
		}
	}
	return result, nil
}

func (a *asset) save(ctx db.Context) (stock StockAsset, retErr error) {
	var err error
	a.mu.Lock()
	defer a.mu.Unlock()
	defer func() {
		if a.model == nil {
			return
		}
		tags, err := db.AddTags(a.model.OID, a.tagSpecs).RunRW(ctx)
		if err != nil {
			stock, retErr = nil, err
			return
		}
		a.model.Tags = tags
	}()
	var originModel *db.AssetModel
	if a.origin != nil {
		_, err = a.origin.save(ctx)
		if err != nil {
			return nil, err
		}
		originModel = a.origin.model
	}

	if a.model != nil {
		return a, nil
	}

	var transSeq db.StepListTpl
	if a.transform != nil {
		transSeq = a.transform.Model()
	}

	model, err := db.AssetTpl{
		Origin:   originModel,
		Basename: a.basename,
		Url:      a.url,
		ExtData:  a.ext,
		TransSeq: transSeq,
		Content:  a.content,
	}.Create().RunRW(ctx)

	if err != nil {
		return nil, err
	}

	a.model = model
	return a, nil
}
