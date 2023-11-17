package asset

import (
	"fmt"

	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
	"github.com/hsfzxjy/imbed/util"
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

	var revertSaveFile bool

	for _, a := range assets {
		revert, err := a.saveFile(app)
		defer func() {
			if revertSaveFile {
				revert.Call()
			}
		}()
		if err != nil {
			revertSaveFile = true
			return nil, err
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

	fhash, err := a.content.GetHash()
	if err != nil {
		return nil, err
	}

	var transSeq db.StepListTpl
	if a.transform != nil {
		transSeq = a.transform.Model()
	} else {
		// this is an external asset, try not duplicate in DB
		it, err := assetq.ByFHash(ndl.RawFull(fhash.RawString())).RunR(ctx)
		if err == nil && it.HasNext() {
			model := it.Next()
			if model.IsErr() {
				return nil, model.UnwrapErr()
			}
			a.model = model.Unwrap()
			return a, nil
		}
	}

	model, err := db.AssetTpl{
		Origin:   originModel,
		FHash:    fhash,
		Basename: a.basename,
		Url:      a.url,
		ExtData:  a.ext,
		TransSeq: transSeq,
	}.Create().RunRW(ctx)

	if err != nil {
		return nil, err
	}

	a.model = model
	return a, nil
}

func (a *asset) saveFile(app core.App) (revert util.RevertFunc, err error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if a.model == nil {
		return nil, fmt.Errorf("model not saved")
	}
	if a.origin != nil {
		revert, err = a.origin.saveFile(app)
		if err != nil {
			return revert, err
		}
	}
	r, err := a.content.BytesReader()
	if err != nil {
		return revert, err
	}
	revert2, err := util.SafeWriteFile(
		r,
		app.FilePath(a.model.Filename()),
	)
	revert = revert.Then(revert2)
	return revert, err
}
