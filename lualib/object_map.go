package lualib

import (
	"errors"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

type mapMTBuilder struct {
	ReadOnly bool
}

func (mt mapMTBuilder) get() *lua.LTable {
	t := newTable()
	t.RawSetString("__index", newGFunction(mt.__index))
	t.RawSetString("__newindex", newGFunction(mt.__newindex))
	return t
}

func (mt mapMTBuilder) __index(L *lua.LState) int {
	m := L.CheckUserData(1).Value.(*Map)
	key := L.CheckString(2)
	if v, ok := m.m[key]; ok {
		L.Push(v.AsLValue(L, mt.ReadOnly))
		return 1
	}
	L.Push(lua.LNil)
	return 1
}

func (mt mapMTBuilder) __newindex(L *lua.LState) int {
	if mt.ReadOnly {
		L.RaiseError("read only")
		return 0
	}
	m := L.CheckUserData(1).Value.(*Map)
	key := L.CheckString(2)
	v := L.CheckAny(3)
	val, err := m.ValueChecker.check(L, v)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	m.m[key] = val
	return 0
}

var (
	mapMT   = sync.OnceValue(mapMTBuilder{ReadOnly: false}.get)
	mapMTRO = sync.OnceValue(mapMTBuilder{ReadOnly: true}.get)
)

type Map struct {
	ValueChecker objectChecker

	m map[string]Object

	readonly bool

	cacheOnce sync.Once
	ud, udro  *lua.LUserData
}

func (m *Map) AsLValue(L *lua.LState, readonly bool) lua.LValue {
	m.cacheOnce.Do(func() {
		m.ud = newUserData(m)
		m.ud.Metatable = mapMT()
		m.udro = newUserData(m)
		m.udro.Metatable = mapMTRO()
	})
	if readonly {
		return m.udro
	} else {
		return m.ud
	}
}

func (m *Map) IsIntegral() error {
	for _, v := range m.m {
		if err := v.IsIntegral(); err != nil {
			return err
		}
	}
	return nil
}

type MapChecker struct {
	ValueChecker objectChecker
}

func (mc MapChecker) check(L *lua.LState, v lua.LValue) (Object, error) {
	switch t := v.(type) {
	case *lua.LTable:
		m := &Map{
			ValueChecker: mc.ValueChecker,
			m:            make(map[string]Object),
		}
		var err error
		t.ForEach(func(k, v lua.LValue) {
			if err != nil {
				return
			}
			key, ok := k.(lua.LString)
			if !ok {
				err = errors.New("not a string key")
				return
			}
			val, err := mc.ValueChecker.check(L, v)
			if err != nil {
				return
			}
			m.m[string(key)] = val
		})
		if err != nil {
			return nil, err
		}
		return m, nil
	case *lua.LUserData:
		if m, ok := t.Value.(*Map); ok && m.ValueChecker == mc.ValueChecker {
			return m, nil
		} else {
			return nil, errors.New("not a desired map")
		}
	default:
		return nil, errors.New("not a table")
	}
}
