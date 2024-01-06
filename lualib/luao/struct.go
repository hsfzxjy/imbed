package luao

import lua "github.com/hsfzxjy/gopher-lua"

var structCDH *lua.CustomDataHelper[Struct]

func init() {
	mt := lua.NewTable()
	mt.RawSetString("__index", (&lua.LFunction{
		IsG:       true,
		IsFast:    true,
		GFunction: struct__index,
	}).AsLValue())
	mt.RawSetString("__newindex", (&lua.LFunction{
		IsG:       true,
		IsFast:    true,
		GFunction: struct__newindex,
	}).AsLValue())
	structCDH = lua.RegisterCustomData[Struct](mt)
}

func struct__index(L *lua.LState) int {
	s := structCDH.Must(L.Get(1))
	k := L.CheckString(2)
	idx, _, ok := s.p.fieldP(k)
	if !ok {
		L.RaiseError("invalid field: %s", k)
		return 1
	}
	L.Push(s.values[idx].AsLValue())
	return 1
}

func struct__newindex(L *lua.LState) int {
	s := structCDH.Must(L.Get(1))
	k := L.CheckString(2)
	v := L.Get(3)
	idx, fp, ok := s.p.fieldP(k)
	if !ok {
		L.RaiseError("invalid field: %s", k)
		return 0
	}
	obj, err := fp.P.L2G(v)
	if err != nil {
		L.RaiseError("invalid value: %w", err)
		return 0
	}
	s.values[idx] = obj
	return 0
}

type Struct struct {
	p      *StructP
	values []Object
}

func NewStruct(p *StructP) *Struct {
	return &Struct{
		p:      p,
		values: make([]Object, len(p.fields)),
	}
}

func (s *Struct) CheckIntegrity() error {
	return s.p.checkIntegrity(Object{structCDH.AsLValue(s)})
}

type StructP struct {
	name2i map[string]int
	fields []FieldP
}

type FieldP struct {
	Name string
	P    Protocol
}

func NewStructP(fields []FieldP) *StructP {
	sp := new(StructP)
	sp.name2i = make(map[string]int, len(fields))
	sp.fields = fields
	for i, f := range fields {
		_, ok := sp.name2i[f.Name]
		if ok {
			panic("duplicate field name: " + f.Name)
		}
		sp.name2i[f.Name] = i
	}
	return sp
}

func (sp *StructP) fieldP(name string) (int, *FieldP, bool) {
	i, ok := sp.name2i[name]
	if !ok {
		return 0, nil, false
	}
	return i, &sp.fields[i], true
}

func (sp *StructP) fieldPByIndex(i int) *FieldP {
	return &sp.fields[i]
}

func (sp *StructP) L2G(v lua.LValue) (Object, error) {
	if _, ok := structCDH.As(v); ok {
		return Object{v}, nil
	}
	tb, ok := v.AsLTable()
	if !ok {
		return Object{}, &typeError{lua.LTTable, v.Type()}
	}
	s := NewStruct(sp)
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
		idx, fp, ok := sp.fieldP(string(key))
		if !ok {
			err = &noSuchFieldError{string(key)}
			return
		}
		value, err = fp.P.L2G(v)
		if err != nil {
			return
		}
		s.values[idx] = value
	})
	if err != nil {
		return Object{}, err
	}
	return Object{structCDH.AsLValue(s)}, nil
}

func (sp *StructP) checkIntegrity(obj Object) error {
	s, ok := structCDH.As(obj.lv)
	if !ok {
		panic("invalid struct object")
	}
	if len(s.values) != len(sp.fields) {
		panic("invalid struct object")
	}
	for i, v := range s.values {
		if err := sp.fields[i].P.checkIntegrity(v); err != nil {
			return err
		}
	}
	return nil
}
