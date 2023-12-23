package lualib

import (
	"errors"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

var (
	structMT   = sync.OnceValue(structMTBuilder{ReadOnly: false}.get)
	structMTRO = sync.OnceValue(structMTBuilder{ReadOnly: true}.get)
)

type Struct struct {
	checker *structChecker
	objects []Object

	cacheOnce sync.Once
	ud, udro  *lua.LUserData
}

func NewStruct(checker *structChecker) *Struct {
	s := new(Struct)
	s.checker = checker
	s.objects = make([]Object, len(checker.fieldCheckers))
	return s
}

func (s *Struct) AsLValue(L *lua.LState, readonly bool) lua.LValue {
	s.cacheOnce.Do(func() {
		s.ud = newUserData(s)
		s.ud.Metatable = structMT()
		s.udro = newUserData(s)
		s.udro.Metatable = structMTRO()
	})
	if readonly {
		return s.udro
	} else {
		return s.ud
	}
}

func (s *Struct) IsIntegral() error {
	for idx, v := range s.objects {
		checker := s.checker.fieldCheckers[idx]
		if v == nil {
			if _, ok := checker.(PtrChecker); !ok {
				return errors.New("nil field")
			} else {
				continue
			}
		}
		if err := v.IsIntegral(); err != nil {
			return err
		}
	}
	return nil
}

type structMTBuilder struct {
	ReadOnly bool
}

func (b structMTBuilder) get() *lua.LTable {
	mt := newTable()
	mt.RawSetString("__index", newGFunction(b.__index))
	mt.RawSetString("__newindex", newGFunction(b.__newindex))
	mt.RawSetString("__call", newGFunction(b.__call))
	return mt
}

func (s structMTBuilder) __call(L *lua.LState) int {
	ss := L.CheckUserData(1).Value.(*Struct)
	param := L.Get(2)
	objs, err := ss.checker.checkNewObjects(L, param)
	if err != nil {
		L.RaiseError(err.Error())
	} else {
		ss.objects = objs
	}
	return 0
}

func (b structMTBuilder) __index(L *lua.LState) int {
	m := L.CheckUserData(1).Value.(*Struct)
	key := L.CheckString(2)
	idx, ok := m.checker.key2i[key]
	if !ok {
		L.RaiseError("no such field: %s", key)
		return 0
	}
	o := m.objects[idx]
	if o == nil {
		if _, ok := m.checker.fieldCheckers[idx].(PtrChecker); ok {
			L.Push(lua.LNil)
			return 1
		} else {
			L.RaiseError("uninitialized field")
			return 0
		}
	} else {
		L.Push(o.AsLValue(L, b.ReadOnly))
	}
	return 1
}

func (s structMTBuilder) __newindex(L *lua.LState) int {
	if s.ReadOnly {
		L.RaiseError("read only")
		return 0
	}
	m := L.CheckUserData(1).Value.(*Struct)
	key := L.CheckString(2)
	v := L.CheckAny(3)
	idx, ok := m.checker.key2i[key]
	if !ok {
		L.RaiseError("no such field: %s", key)
		return 0
	}
	val, err := m.checker.fieldCheckers[idx].check(L, v)
	if err != nil {
		L.RaiseError(err.Error())
		return 0
	}
	m.objects[idx] = val
	return 0
}

type FieldChecker struct {
	Name    string
	Checker objectChecker
}

type structChecker struct {
	fieldCheckers []objectChecker
	key2i         map[string]int
}

func NewStructChecker(checkers []FieldChecker) *structChecker {
	c := new(structChecker)
	c.fieldCheckers = make([]objectChecker, len(checkers))
	c.key2i = make(map[string]int, len(checkers))
	for i, checker := range checkers {
		c.fieldCheckers[i] = checker.Checker
		c.key2i[checker.Name] = i
	}
	return c
}

func (c *structChecker) checkNewObjects(L *lua.LState, v lua.LValue) ([]Object, error) {
	switch t := v.(type) {
	case *lua.LTable:
		objects := make([]Object, len(c.fieldCheckers))
		for k, idx := range c.key2i {
			val := t.RawGetString(k)
			o, err := c.fieldCheckers[idx].check(L, val)
			if err != nil {
				return nil, err
			}
			objects[idx] = o
		}
		return objects, nil
	case *lua.LUserData:
		if s, ok := t.Value.(*Struct); ok && s.checker == c {
			return s.objects, nil
		}
		return nil, errors.New("not a struct")
	}
	return nil, errors.New("not a struct")
}

func (c *structChecker) check(L *lua.LState, v lua.LValue) (Object, error) {
	objs, err := c.checkNewObjects(L, v)
	if err != nil {
		return nil, err
	}
	return &Struct{checker: c, objects: objs}, nil
}
