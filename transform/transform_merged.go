package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/tinylib/msgp/msgp"
)

type mergedTransform struct {
	transforms []Transform
	ref.LazyEncodableObject[*mergedTransform]
}

func newMergedTransform(transforms []Transform) *mergedTransform {
	t := new(mergedTransform)
	t.transforms = transforms
	t.Inner = t
	return t
}

func _() { var _ Transform = &mergedTransform{} }

func (m *mergedTransform) Name() string { return "mergedTransform" }

func (m *mergedTransform) Apply(app core.App, a asset.Asset) (asset.Update, error) {
	var updates []asset.Update
	var tmpAsset = a
	for _, t := range m.transforms {
		u, err := t.Apply(app, tmpAsset)
		if err != nil {
			return nil, err
		}
		tmpAsset, err = asset.ApplyUpdate(tmpAsset, nil, u)
		if err != nil {
			return nil, err
		}
		updates = append(updates, u)
	}
	return asset.MergeUpdates(updates...), nil
}

func (m *mergedTransform) EncodeSelf(w *msgp.Writer) error {
	for _, t := range m.transforms {
		encoded, err := t.GetRawEncoded()
		if err != nil {
			return err
		}
		err = w.Append(encoded...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mergedTransform) Kind() Kind {
	var firstKind = m.transforms[0].Kind()
	for _, t := range m.transforms[1:] {
		if firstKind != t.Kind() {
			panic("inconsistent kind")
		}
	}
	return firstKind
}

func (m *mergedTransform) AssociatedConfigs() []ref.EncodableObject {
	var ret = make([]ref.EncodableObject, 0, len(m.transforms))
	for _, t := range m.transforms {
		ret = append(ret, t.AssociatedConfigs()...)
	}
	return ret
}
