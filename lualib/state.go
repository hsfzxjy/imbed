package lualib

import (
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type lstate = lua.LState

type LState struct{ *lstate }

func newLState(options lua.Options) LState {
	return LState{lua.NewState(options)}
}

type LStatePool struct {
	mu    sync.Mutex
	saved map[lua.Options]*[]LState
}

func NewLStatePool() *LStatePool {
	return &LStatePool{saved: make(map[lua.Options]*[]LState)}
}

func (p *LStatePool) New(options lua.Options) LState {
	return newLState(options)
}

func (p *LStatePool) NewDefault() LState {
	return newLState(lua.Options{})
}

func (p *LStatePool) Get(options lua.Options) LState {
	p.mu.Lock()
	defer p.mu.Unlock()
	if Ls, ok := p.saved[options]; ok {
		if len(*Ls) > 0 {
			L := (*Ls)[len(*Ls)-1]
			(*Ls)[len(*Ls)-1] = LState{nil}
			*Ls = (*Ls)[:len(*Ls)-1]
			return L
		}
	}
	return newLState(options)
}

func (p *LStatePool) Put(L LState) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if Ls, ok := p.saved[L.Options]; ok {
		*Ls = append(*Ls, L)
	} else {
		Ls := make([]LState, 0, 4)
		Ls = append(Ls, L)
		p.saved[L.Options] = &Ls
	}
}

func (p *LStatePool) Shutdown() {
	for _, Ls := range p.saved {
		for _, L := range *Ls {
			L.Close()
		}
	}
}
