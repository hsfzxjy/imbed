package transform_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/asset/tag"
	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/transform"
	"github.com/hsfzxjy/imbed/util/fastbuf"
	"github.com/stretchr/testify/require"
)

type Applier struct {
	transform.Category
	X uint64
}

func (Applier) Apply(core.App, asset.Asset) (asset.Update, error) { return nil, nil }
func (Applier) EncodeMsg(*fastbuf.W)                              {}

type TFListSeed struct {
	Size, IsTerminal, ABCat, ForceTerminal uint64
}

func (s TFListSeed) Build() StepAtomList {
	size := s.Size % 65
	tfs := make(StepAtomList, size)
	for i := uint64(0); i < size; i++ {
		mask := uint64(1) << i
		var cat transform.Category
		if s.IsTerminal&mask != 0 {
			cat = transform.Terminal
		} else if s.ABCat&mask != 0 {
			cat = "A"
		} else {
			cat = "B"
		}
		var kind tag.Kind
		if s.ForceTerminal&mask != 0 {
			kind = tag.Normal
		}
		tfs[i] = &transform.StepAtom{
			Category: cat,
			Tag:      tag.Spec{Kind: kind},
			Applier:  &Applier{cat, i},
		}
	}
	return tfs
}

type StepAtomList transform.StepAtomList

func (l StepAtomList) String() string {
	var b strings.Builder
	for _, t := range l {
		b.WriteString(string(t.Category))
		b.WriteRune('\t')
		if t.ForceTerminal() {
			b.WriteString("ForceTerminal")
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (l StepAtomList) As() transform.StepAtomList {
	return transform.StepAtomList(l)
}

func FuzzPartition(f *testing.F) {
	f.Add(uint64(0), uint64(0b1), uint64(0b01), uint64(0b01))
	f.Fuzz(func(t *testing.T, size, ist, abcat, forcet uint64) {
		tfs := TFListSeed{size, ist, abcat, forcet}.Build()
		repr := tfs.String()
		result := (&transform.Scheduler{Sal: tfs.As()}).Partition()
		lasti := 0
		for _, span := range result {
			require.Equalf(t, lasti, span.Start, repr)
			require.Lessf(t, span.Start, span.End, repr)
			lasti = span.End

			startCat := tfs[span.Start].Category
			if startCat.IsTerminal() {
				require.Equalf(t, span.Start+1, span.End, repr)
			} else {
				for i := span.Start; i < span.End; i++ {
					tf := tfs[i]
					require.Equalf(t, startCat, tf.Category, repr)
					if i < span.End-1 {
						require.Falsef(t, tf.ForceTerminal(), repr)
					}
				}
			}
		}
		require.Equalf(t, lasti, len(tfs), repr)
	})
}

type Composer struct {
	Throws bool
	transform.Category
}

var ErrCompose = errors.New("compose error")

func (c *Composer) Compose(appliers []transform.Applier) (transform.Applier, error) {
	if c.Throws {
		return nil, ErrCompose
	}
	if len(appliers) <= 1 {
		return nil, errors.New("too short")
	}
	var x uint64
	for _, a := range appliers {
		a := a.(*Applier)
		if a.Category != c.Category {
			return nil, fmt.Errorf("inconsistent category %s != %s", a.Category, c.Category)
		}
		x += a.X
	}
	return &Applier{c.Category, x}, nil
}

type ComposerMSeed uint64

func (s ComposerMSeed) Build() ComposeM {
	m := make(map[transform.Category]transform.Composer)
	if s&(1<<0) != 0 {
		m["A"] = &Composer{
			s&(1<<1) != 0,
			"A",
		}
	}
	if s&(1<<2) != 0 {
		m["B"] = &Composer{
			s&(1<<3) != 0,
			"B",
		}
	}
	return m
}

func (s ComposerMSeed) HasError(tfs StepAtomList) bool {
	AThrows := s&0b11 == 0b11
	BThrows := s&0b1100 == 0b1100
	var maxA, maxB, A, B int
	reset := func() { A, B = 0, 0 }
	for _, t := range tfs {
		a := t.Applier.(*Applier)
		switch a.Category {
		case "A":
			B = 0
			A++
			maxA = max(A, maxA)
			if t.ForceTerminal() {
				reset()
			}
		case "B":
			A = 0
			B++
			maxB = max(B, maxB)
			if t.ForceTerminal() {
				reset()
			}
		default:
			reset()
		}
	}
	return AThrows && maxA > 1 ||
		BThrows && maxB > 1
}

type ComposeM map[transform.Category]transform.Composer

func (m ComposeM) String() string {
	return spew.Sdump(m)
}

func FuzzAssemble(f *testing.F) {
	f.Add(uint64(0), uint64(0b1), uint64(0b01), uint64(0b01), uint64(0b01))
	f.Fuzz(func(t *testing.T, size, ist, abcat, forcet, composem uint64) {
		tfs := TFListSeed{size, ist, abcat, forcet}.Build()
		cmSeed := ComposerMSeed(composem)
		cm := cmSeed.Build()
		repr := tfs.String() + "\n" + cm.String()
		rootSteps, err := transform.Scheduler{tfs.As(), cm}.Assemble()
		if cmSeed.HasError(tfs) {
			require.ErrorIsf(t, err, ErrCompose, repr)
			return
		} else {
			require.ErrorIsf(t, err, nil, repr)
		}
		var AX uint64
		lasti := 0
		steps := rootSteps
		for len(steps) > 0 {
			for i, step := range steps {
				for _, a := range step.Appliers {
					AX += a.(*Applier).X
				}
				require.Equalf(t, lasti, step.Start(), repr)
				lasti = step.End()
				if i < len(steps)-1 {
					require.Truef(t, step.IsTerminal(), repr)
					require.Emptyf(t, step.Next, repr)
				} else {
					steps = step.Next
				}
			}
		}
		n := uint64(len(tfs))
		sum := n * (n - 1) / 2
		require.Equalf(t, sum, AX, repr)
	})
}
