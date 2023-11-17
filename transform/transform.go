package transform

import (
	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/db"
)

type Transform struct {
	Name string
	Applier
	Category

	Tag tag.Spec

	model db.StepTpl
}

func (t *Transform) ForceTerminal() bool {
	return t.Tag.Kind != tag.None
}
