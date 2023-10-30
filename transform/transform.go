package transform

import (
	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/core/ref"
)

type Transform struct {
	Name string
	Applier
	Category
	Data   *Data
	Config ref.EncodableObject

	Tag tag.Spec
}

func (t *Transform) ForceTerminal() bool {
	return t.Tag.Kind != tag.None
}
