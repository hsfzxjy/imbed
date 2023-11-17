package transform

import (
	"sync"

	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/db"
)

type step struct {
	globalOffset int
	atoms        []*stepAtom
	Appliers     []Applier
	Next         stepSet

	model     db.StepListTpl
	modelOnce sync.Once
}

func (s *step) Start() int {
	return s.globalOffset
}

func (s *step) End() int {
	return s.globalOffset + len(s.atoms)
}

func (s *step) IsTerminal() bool {
	return len(s.atoms) == 1 && s.atoms[0].Category.IsTerminal()
}

func (s *step) TagSpec() tag.Spec {
	return s.atoms[len(s.atoms)-1].Tag
}

func (s *step) Model() db.StepListTpl {
	s.modelOnce.Do(func() {
		if s.model != nil {
			return
		}
		list := make([]*db.StepTpl, len(s.atoms))
		for i, atom := range s.atoms {
			list[i] = &atom.model
		}
		s.model = db.NewStepListTpl(list, s.IsTerminal())
	})
	return s.model
}
