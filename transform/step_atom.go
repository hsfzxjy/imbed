package transform

import (
	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/db"
)

type stepAtom struct {
	Name string
	Applier
	Category

	paramsPos pos.P

	Tag tag.Spec

	model db.StepTpl
}

func (t *stepAtom) ForceTerminal() bool {
	return t.Tag.Kind != tag.None
}

type StepAtomList []*stepAtom

func (l StepAtomList) rangeSpan(span span) StepAtomList {
	return l[span.Start:span.End]
}

func (l StepAtomList) getPos() pos.P {
	if len(l) == 0 {
		return pos.P{}
	}
	var p pos.P
	for _, t := range l {
		p = p.Add(t.paramsPos)
	}
	return p
}
