package pos

import (
	"errors"

	"github.com/davecgh/go-spew/spew"
)

type Source struct {
	FmtFunc func(str string, start, end uint32) string
}

func (s *Source) checkEqual(other *Source) {
	if s != other {
		panic("pos: source not equal")
	}
}

type errorWithPos struct {
	error
	P
}

func (s *errorWithPos) Unwrap() error {
	return s.error
}

func (s *errorWithPos) GetPos() P {
	return s.P
}

type errorStringWithPos struct {
	string
	P
}

func (s *errorStringWithPos) Error() string {
	return s.string
}

func (s *errorStringWithPos) GetPos() P {
	return s.P
}

type PosGetter interface {
	GetPos() P
}

func FmtError(err error) string {
	var pos PosGetter
	if errors.As(err, &pos) {
		return pos.GetPos().FmtError(err)
	}
	return err.Error()
}

type P struct {
	src   *Source
	start uint32
	end   uint32
}

func (p P) GetPos() P {
	return p
}

func (p P) FmtError(err error) string {
	if err == nil {
		return ""
	}
	return p.src.FmtFunc(err.Error(), p.start, p.end)
}

func (p P) WrapErrorString(err string) error {
	if p.src == nil {
		return errors.New(err)
	}
	return &errorStringWithPos{err, p}
}

func (p P) WrapError(err error) error {
	if p.src == nil || err == nil {
		return err
	}
	if pos, ok := err.(PosGetter); ok {
		spew.Dump(pos)
		op := pos.GetPos()
		if op.src != nil {
			p = op
		}
	}
	return &errorWithPos{err, p}
}

func (p P) Or(other P) P {
	if p.src == nil {
		return other
	}
	return p
}

func (p P) Start() P {
	if p.src == nil {
		return P{}
	}
	return P{p.src, p.start, p.start}
}

func (p P) End() P {
	if p.src == nil {
		return P{}
	}
	return P{p.src, p.end, p.end}
}

func (p P) Sub(other P) P {
	p.src.checkEqual(other.src)
	if p.start == p.end && other.start == other.end {
		l, r := p.start, other.start
		if l > r {
			l, r = r, l
		}
		return P{p.src, l, r}
	}
	return p
}

func (p P) Add(other P) P {
	if p.src == nil {
		return other
	}
	if other.src == nil {
		return p
	}
	p.src.checkEqual(other.src)
	return P{p.src, min(p.start, other.start), max(p.end, other.end)}
}

func (p *P) ExtendEnd(off int) {
	p.end += uint32(off)
}

func New(src *Source, start, end uint32) P {
	return P{src, start, end}
}
