package luao

import lua "github.com/hsfzxjy/gopher-lua"

var mapCDH *lua.CustomDataHelper[Map]

func init() {
	mt := lua.NewTable()
	mt.RawSetString("__index", (&lua.LFunction{
		IsG:       true,
		IsFast:    true,
		GFunction: map__index,
	}).AsLValue())
	mt.RawSetString("__newindex", (&lua.LFunction{
		IsG:       true,
		IsFast:    true,
		GFunction: map__newindex,
	}).AsLValue())
	mapCDH = lua.RegisterCustomData[Map](mt)
}

func map__index(L *lua.LState) int {
	m := mapCDH.Must(L.Get(1))
	k := L.CheckString(2)
	v, ok := m.Get(k)
	if !ok {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(v.AsLValue())
	return 1
}

func map__newindex(L *lua.LState) int {
	m := mapCDH.Must(L.Get(1))
	k := L.CheckString(2)
	v := L.Get(3)
	obj, err := m.valueP.L2G(v)
	if err != nil {
		L.RaiseError("invalid value: %s", err)
		return 0
	}
	m.Set(k, obj)
	return 0
}

// A Map is a mapping from string to Object. The value type is homogeneous.
type Map struct {
	valueP Protocol
	m      map[string]Object
}

func NewMap(valueP Protocol) *Map {
	return &Map{
		valueP: valueP,
		m:      make(map[string]Object),
	}
}

func (t *Map) Len() int {
	return len(t.m)
}

func (t *Map) ValueProtocol() Protocol {
	return t.valueP
}

func (t *Map) Get(k string) (Object, bool) {
	v, ok := t.m[k]
	return v, ok
}

func (t *Map) Set(k string, v Object) {
	t.m[k] = v
}

type MapP struct {
	valueP Protocol
}

func (p *MapP) L2G(v lua.LValue) (Object, error) {
	if _, ok := mapCDH.As(v); ok {
		return Object{v}, nil
	}
	tb, ok := v.AsLTable()
	if !ok {
		return Object{}, &typeError{lua.LTTable, v.Type()}
	}
	m := NewMap(p.valueP)
	var err error
	tb.ForEach(func(k, v lua.LValue) {
		if err != nil {
			return
		}
		var key lua.LString
		var value Object
		key, ok = k.AsLString()
		if !ok {
			err = &typeError{lua.LTString, k.Type()}
			return
		}
		value, err = p.valueP.L2G(v)
		if err != nil {
			return
		}
		m.Set(string(key), value)
	})
	if err != nil {
		return Object{}, err
	}
	return Object{mapCDH.AsLValue(m)}, nil
}
