package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/asset/tag"
)

func (s *step) Run(cache *assetCache, upstream asset.Asset, ret *[]asset.Asset) error {
	var result asset.Asset
	var err error
	if result, err = cache.Lookup(upstream, s); err != nil {
		return err
	}
	if result == nil {
		a := upstream
		updates := make([]asset.Update, 0, len(s.Appliers))
		for _, applier := range s.Appliers {
			update, err := applier.Apply(cache.app, a)
			if err != nil {
				return err
			}
			if update == nil {
				continue
			}
			a, err = asset.ApplyUpdate(a, s, update)
			if err != nil {
				return err
			}
			updates = append(updates, update)
		}
		if len(updates) > 0 {
			result, err = asset.ApplyUpdate(upstream, s, asset.MergeUpdates(updates...))
			if err != nil {
				return err
			}
		} else {
			result = upstream
		}
	}
	spec := s.TagSpec()
	if spec.Kind == tag.Auto {
		spec.Name = result.BaseName()
	}
	result = asset.Tag(result, spec)
	if s.IsTerminal() && result != upstream {
		*ret = append(*ret, result)
	}
	return s.Next.Run(cache, result, ret)
}

type stepSet []*step

func (ss stepSet) Run(cache *assetCache, upstream asset.Asset, ret *[]asset.Asset) error {
	for _, step := range ss {
		err := step.Run(cache, upstream, ret)
		if err != nil {
			return err
		}
	}
	return nil
}
