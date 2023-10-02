package asset

import (
	"fmt"

	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/util"
)

func SaveAll(ctx db.Context, app core.App, assets []Asset) error {
	for _, a := range assets {
		if err := a.save(ctx); err != nil {
			return err
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
			return err
		}
	}
	return nil
}

func (a *asset) saveConfigs(ctx db.Context) ([]ref.Sha256Hash, error) {
	if a.transform == nil {
		return nil, nil
	}
	cfgs := a.transform.AssociatedConfigs()
	var seen = make(map[ref.Sha256Hash]struct{}, len(cfgs))
	var ret = make([]ref.Sha256Hash, 0, len(cfgs))
	for _, cfg := range cfgs {
		hash, err := cfg.GetSha256Hash()
		if err != nil {
			return nil, err
		}
		if _, ok := seen[hash]; ok {
			continue
		}
		seen[hash] = struct{}{}
		ret = append(ret, hash)
		raw, err := cfg.GetRawEncoded()
		if err != nil {
			return nil, err
		}
		_, err = db.ConfigModel{
			Raw:  raw,
			Hash: hash,
		}.Create().RunRW(ctx)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (a *asset) save(ctx db.Context) error {
	var err error
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.model != nil {
		return nil
	}
	var originModel *db.AssetModel
	if a.origin != nil {
		err = a.origin.save(ctx)
		if err != nil {
			return err
		}
		originModel = a.origin.model
	}

	var transSeq db.TransSeq
	if a.transform != nil {
		transSeq.ConfigHashes, err = a.saveConfigs(ctx)
		if err != nil {
			return err
		}
		transSeq.Raw, err = a.transform.GetRawEncoded()
		if err != nil {
			return err
		}
		transSeq.Hash, err = a.transform.GetSha256Hash()
		if err != nil {
			return err
		}
	}

	fid, err := content.BuildFID(a.content, a.basename)
	if err != nil {
		return err
	}

	model, err := db.AssetTemplate{
		Origin:   originModel,
		FID:      fid,
		Url:      a.url,
		ExtData:  a.ext,
		TransSeq: transSeq,
	}.Create().RunRW(ctx)

	if err != nil {
		return err
	}

	a.model = model
	return nil
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
		app.FilePath(a.model.FID.Humanize()),
	)
	revert = revert.Then(revert2)
	return revert, err
}
