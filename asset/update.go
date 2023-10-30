package asset

import (
	"errors"

	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/content"
	"github.com/hsfzxjy/imbed/util"
)

type Update interface {
	apply(a *asset) error
}

type updateFileExtension struct{ ext string }

func (x updateFileExtension) apply(a *asset) error {
	a.basename = util.ReplaceExt(a.basename, x.ext)
	return nil
}

func UpdateFileExtension(ext string) Update {
	return updateFileExtension{ext}
}

type updateContent struct {
	content content.Content
}

func (x updateContent) apply(a *asset) error {
	a.content = x.content
	return nil
}

func UpdateContent(content content.Content) Update {
	return updateContent{content}
}

type updateExt struct {
	ext []byte
}

func (x updateExt) apply(a *asset) error {
	a.ext = x.ext
	return nil
}

func UpdateExt(ext []byte) Update {
	return updateExt{ext}
}

type updateUrl struct{ url string }

func (x updateUrl) apply(a *asset) error {
	if a.url != "" {
		return errors.New("UpdateUrl() on single asset more than once")
	}
	a.url = x.url
	return nil
}

func UpdateUrl(url string) Update { return updateUrl{url} }

type updates struct {
	list []Update
}

func (x updates) apply(a *asset) error {
	for _, u := range x.list {
		err := u.apply(a)
		if err != nil {
			return err
		}
	}
	return nil
}

func MergeUpdates(updateList ...Update) Update {
	return updates{updateList}
}

func ApplyUpdate(ia Asset, transform Transform, up Update) (Asset, error) {
	a := ia.(*asset)
	newAsset := new(asset)
	newAsset.origin = a
	newAsset.primaryInfo = a.primaryInfo
	newAsset.updatingInfo = updatingInfo{
		transform: transform,
	}
	err := up.apply(newAsset)
	if err != nil {
		return nil, err
	}
	return newAsset, nil
}

func Tag(ia Asset, spec tag.Spec) Asset {
	if spec.Kind == tag.None {
		return ia
	}
	a := ia.(*asset)
	a.tagSpecs = append(a.tagSpecs, spec)
	return a
}
