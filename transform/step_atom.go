package transform

import (
	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/db"
)

type stepAtom struct {
	Name string
	Applier
	Category

	Tag tag.Spec

	model db.StepTpl
}

func (t *stepAtom) ForceTerminal() bool {
	return t.Tag.Kind != tag.None
}

type StepAtomList []*stepAtom

func (l StepAtomList) Range(span span) StepAtomList {
	return l[span.Start:span.End]
}
