package transform

func Schedule(reg *Registry, sal StepAtomList) (*Graph, error) {
	root, err := scheduler{sal, reg.composerTable}.Assemble()
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

type schedulerState struct {
	*scheduler
	lastNT   *step
	root     stepSet
	LastIsNT bool
}

func (ss *schedulerState) saveStep(span span, appliers []Applier) {
	step := &step{
		atoms:        ss.Sal.Range(span),
		globalOffset: span.Start,
		Appliers:     appliers,
	}
	if ss.lastNT == nil {
		ss.root = append(ss.root, step)
	} else {
		ss.lastNT.Next = append(ss.lastNT.Next, step)
	}
	lastIsNT := !step.IsTerminal()
	if lastIsNT {
		ss.lastNT = step
	}
	ss.LastIsNT = lastIsNT
}

func (ss *schedulerState) extendLastNT(span span, appliers []Applier) {
	if ss.lastNT == nil {
		panic("extendLastNT: lastNT == nil")
	}
	if !ss.LastIsNT {
		panic("extendLastNT: lastIsNT == false")
	}
	lastNT := ss.lastNT
	lastNT.Appliers = append(lastNT.Appliers, appliers...)
	lastNT.atoms = lastNT.atoms[:len(lastNT.atoms)+span.Len()]
}

type scheduler struct {
	Sal    StepAtomList
	CTable map[Category]Composer
}

func (s *scheduler) Partition() []span {
	var spans = make([]span, 0, len(s.Sal))
	var n = -1
	var lastOpen = false
	for i, t := range s.Sal {
		cat := t.Category
		if lastOpen {
			if cat.IsTerminal() || cat != spans[n].Category {
				spans[n].End = i
				lastOpen = false
			} else {
				if t.ForceTerminal() {
					spans[n].End = i + 1
					lastOpen = false
				}
				continue
			}
		}
		var end = 0
		lastOpen = true
		if cat.IsTerminal() || t.ForceTerminal() {
			end = i + 1
			lastOpen = false
		}
		n++
		spans = append(spans, span{Start: i, End: end, Category: t.Category})
	}
	if lastOpen {
		spans[n].End = len(s.Sal)
	}
	return spans
}

func (s scheduler) Assemble() (stepSet, error) {
	spans := s.Partition()
	var state = schedulerState{scheduler: &s}
	for _, span := range spans {
		if span.IsTerminal() {
			if span.Start+1 != span.End {
				panic("buildGraph: span.start+1 != span.end for terminal step")
			}
			state.saveStep(span, []Applier{s.Sal[span.Start].Applier})
		} else {
			appliers, err := s.resolveAppliers(span)
			if err != nil {
				return nil, err
			}
			if state.LastIsNT {
				state.extendLastNT(span, appliers)
			} else {
				state.saveStep(span, appliers)
			}
		}
	}
	return state.root, nil
}

func (s *scheduler) resolveAppliers(span span) ([]Applier, error) {
	var ret = make([]Applier, 0, span.Len())
	for _, t := range s.Sal.Range(span) {
		ret = append(ret, t.Applier)
	}
	if span.Len() > 1 {
		if composer, ok := s.CTable[span.Category]; ok {
			applier, err := composer.Compose(ret)
			if err != nil {
				return nil, err
			}
			clear(ret)
			ret = ret[:1]
			ret[0] = applier
		}
	}
	return ret, nil
}
