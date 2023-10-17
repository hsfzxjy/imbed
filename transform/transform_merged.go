package transform

import (
	"bytes"
	"crypto/sha256"
	"slices"

	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
)

type mergedTransform struct {
	transforms []Transform
	encodable
}

func newMergedTransform(transforms []Transform) *mergedTransform {
	t := new(mergedTransform)
	t.transforms = transforms
	t.Compute = t.compute
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

func (m *mergedTransform) compute() {
	var (
		buf     bytes.Buffer
		err     error
		hash    ref.Sha256Hash
		encoded []byte
	)
	for _, t := range m.transforms {
		encoded, err = t.GetRawEncoded()
		if err != nil {
			goto ERROR
		}
		buf.Write(encoded)
	}
	encoded = slices.Clone(buf.Bytes())
	buf.Reset()

	if len(m.transforms) == 1 {
		hash, err = m.transforms[0].GetSha256Hash()
		if err != nil {
			goto ERROR
		}
	} else {
		buf.Grow(len(m.transforms) * sha256.Size)
		for _, t := range m.transforms {
			hash, err = t.GetSha256Hash()
			if err != nil {
				goto ERROR
			}
			buf.Write(ref.AsRaw(hash))
		}
		hash = ref.Sha256HashSum(buf.Bytes())
	}

	m.encoded, m.hash = encoded, hash
	return

ERROR:
	m.encodeError = err
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
