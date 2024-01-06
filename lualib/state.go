package lualib

import (
	"sync"

	lua "github.com/hsfzxjy/gopher-lua"
)

type LStatePool struct {
	options lua.Options
	mu      sync.Mutex
	saved   []*lua.LState
}

func NewLStatePool(options lua.Options) *LStatePool {
	return &LStatePool{options: options}
}

func (p *LStatePool) New() *lua.LState {
	return lua.NewState(p.options)
}

func (p *LStatePool) Get() *lua.LState {
	p.mu.Lock()
	defer p.mu.Unlock()

	saved := p.saved
	if len(saved) > 0 {
		L := saved[len(saved)-1]
		(saved)[len(saved)-1] = nil
		saved = saved[:len(saved)-1]
		p.saved = saved
		return L
	}
	return p.New()
}

func (p *LStatePool) Put(L *lua.LState) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if L.Options != p.options {
		panic("LStatePool: L.Options != p.options")
	}
	p.saved = append(p.saved, L)
}

func (p *LStatePool) Shutdown() {
	for _, L := range p.saved {
		L.Close()
	}
}
