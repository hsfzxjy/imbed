package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/db"
)

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

type step struct {
	transSeq
	Appliers []Applier
	Next     stepSet
}

func (s *step) Run(cache *assetCache, a asset.Asset, ret *[]asset.Asset) error {
	var cached asset.Asset
	var err error
	if cached, err = cache.Lookup(a, &s.transSeq); err != nil {
		return err
	}
	if cached != nil {
		a = cached
	} else {
		for _, applier := range s.Appliers {
			update, err := applier.Apply(cache.app, a)
			if err != nil {
				return err
			}
			a, err = asset.ApplyUpdate(a, s, update)
			if err != nil {
				return nil
			}
		}
	}
	if s.IsTerminal() {
		*ret = append(*ret, a)
	}
	return s.Next.Run(cache, a, ret)
}

type Graph struct {
	reg  *Registry
	root stepSet
}

func (g *Graph) Run(h db.Context, app db.App, initial asset.Asset) ([]asset.Asset, error) {
	var result []asset.Asset
	var cache = assetCache{h, app}
	err := g.root.Run(&cache, initial, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Schedule(reg *Registry, tfList []*Transform) (*Graph, error) {
	root, err := assemble(reg.composerTable, tfList)
	if err != nil {
		return nil, err
	}
	return &Graph{root: root, reg: reg}, nil
}

// span represents an interval [Start, End)
type span struct {
	Start, End int
	Category
}

func (s *span) Len() int {
	return s.End - s.Start
}

func partition(tfList []*Transform) []span {
	var spans = make([]span, 0, len(tfList))
	var n = -1
	var lastOpen = false
	for i, t := range tfList {
		cat := t.Category
		if lastOpen {
			if cat.IsTerminal() || cat != spans[n].Category {
				spans[n].End = i
				lastOpen = false
			} else {
				if t.ForceTerminal {
					spans[n].End = i + 1
					lastOpen = false
				}
				continue
			}
		}
		var end = 0
		lastOpen = true
		if cat.IsTerminal() || t.ForceTerminal {
			end = i + 1
			lastOpen = false
		}
		n++
		spans = append(spans, span{Start: i, End: end, Category: t.Category})
	}
	if lastOpen {
		spans[n].End = len(tfList)
	}
	return spans
}

func assemble(composerm map[Category]Composer, tfList []*Transform) (root stepSet, err error) {
	spans := partition(tfList)
	var lastNT *step
	var lastIsTerminal = true
	saveStep := func(step *step, isTerminal bool) {
		if lastNT == nil {
			root = append(root, step)
		} else {
			lastNT.Next = append(lastNT.Next, step)
		}
		if !isTerminal {
			lastNT = step
		}
		step.Compute = step.compute
	}
	for _, span := range spans {
		if span.IsTerminal() {
			if span.Start+1 != span.End {
				panic("buildGraph: span.start+1 != span.end for terminal step")
			}
			step := &step{
				transSeq: transSeq{
					Seq:   tfList,
					Start: span.Start,
					End:   span.End,
				},
				Appliers: []Applier{tfList[span.Start].Applier},
			}
			saveStep(step, true)
		} else {
			var appliers = make([]Applier, 0, span.Len())
			for _, t := range tfList[span.Start:span.End] {
				appliers = append(appliers, t.Applier)
			}
			if span.Len() > 1 {
				if composer, ok := composerm[span.Category]; ok {
					applier, err := composer.Compose(appliers)
					if err != nil {
						return nil, err
					}
					appliers = appliers[:0]
					appliers = append(appliers, applier)
				}
			}
			if lastIsTerminal {
				saveStep(&step{
					transSeq: transSeq{
						Seq:   tfList,
						Start: span.Start,
						End:   span.End,
					},
					Appliers: appliers,
				}, false)
			} else {
				lastNT.Appliers = append(lastNT.Appliers, appliers...)
				lastNT.End += span.Len()
			}
		}
		lastIsTerminal = span.IsTerminal()
	}
	return root, nil
}
