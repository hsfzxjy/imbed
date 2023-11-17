package transform

import (
	"sync"

	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/db"
)

type transSeq struct {
	Seq        []*Transform
	Start, End int
	model      db.StepListTpl
	modelOnce  sync.Once
}

func (ts *transSeq) IsTerminal() bool {
	return ts.End-ts.Start == 1 && ts.Seq[ts.Start].Category.IsTerminal()
}

func (ts *transSeq) TagSpec() tag.Spec {
	return ts.Seq[ts.End-1].Tag
}

func (ts *transSeq) Model() db.StepListTpl {
	ts.modelOnce.Do(func() {
		if ts.model != nil {
			return
		}
		list := make([]*db.StepTpl, ts.End-ts.Start)
		for i := ts.Start; i < ts.End; i++ {
			list[i-ts.Start] = &ts.Seq[i].model
		}
		ts.model = db.NewStepListTpl(list, ts.Seq[ts.Start].IsTerminal())
	})
	return ts.model
}
