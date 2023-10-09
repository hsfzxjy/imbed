package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/db/assetq"
)

type node struct {
	Parent    *node
	transform Transform
	Children  []*node
}

func (n *node) Apply(cache *assetCache, a asset.Asset, ret *[]asset.Asset) error {
	var err error
	if n.transform != nil {
		var cached asset.Asset
		if cached, err = cache.Lookup(a, n.transform); err != nil {
			return err
		}
		if cached != nil {
			a = cached
		} else {
			u, err := n.transform.Apply(cache.app, a)
			if err != nil {
				return err
			}
			a, err = asset.ApplyUpdate(a, n.transform, u)
			if err != nil {
				return err
			}
		}
	}
	if len(n.Children) == 0 {
		*ret = append(*ret, a)
	} else {
		for _, c := range n.Children {
			err := c.Apply(cache, a, ret)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Graph struct {
	root *node
}

func BuildGraph(transforms []Transform) *Graph {
	var parent, root *node
	root = &node{}
	parent = root
	addToParent := func(t Transform) *node {
		n := new(node)
		n.Parent = parent
		n.transform = t
		parent.Children = append(parent.Children, n)
		return n
	}
	var contentChanger []Transform
	for _, t := range transforms {
		k := t.Kind()
		switch k {
		case KindChangeContent:
			contentChanger = append(contentChanger, t)
		case KindPersist:
			if len(contentChanger) > 0 {
				var parentTransform Transform
				if len(contentChanger) == 1 {
					parentTransform = contentChanger[0]
				} else {
					parentTransform = newMergedTransform(contentChanger)
				}
				n := addToParent(parentTransform)
				parent = n
				contentChanger = nil
			}
			addToParent(t)
		}
	}
	return &Graph{root}
}

func (g *Graph) Compute(app db.App, initial asset.Asset) ([]asset.Asset, error) {
	var results []asset.Asset
	err := app.DB().RunR(func(ctx db.Context) error {
		cache := &assetCache{ctx, app}
		return g.root.Apply(cache, initial, &results)
	})
	if err != nil {
		return nil, err
	}
	return results, nil
}

type assetCache struct {
	ctx db.Context
	app db.App
}

func (c *assetCache) Lookup(a asset.Asset, transform Transform) (asset.Asset, error) {
	transSeqHash, err := transform.GetSha256Hash()
	if err != nil {
		return nil, err
	}
	hash, err := a.Content().GetHash()
	if err != nil {
		return nil, err
	}
	it, err := assetq.ByDependency(hash, transSeqHash).RunR(c.ctx)
	if err != nil {
		return nil, err
	}
	if !it.HasNext() {
		return nil, nil
	}
	model := it.Next()
	return asset.FromDBModel(c.app, model)(nil)
}
