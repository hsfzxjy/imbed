package lualib

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
)

func TestObject(t *testing.T) {
	L := lua.NewState()
	var checker = NewStructChecker([]FieldChecker{
		{"int", IntChecker{}},
		{"oint", &PtrChecker{IntChecker{}}},
		{"bool", BoolChecker{}},
		{"string", StringChecker{}},
		{"intMap", &MapChecker{IntChecker{}}},
	})
	for _, tc := range []struct {
		Code     string
		Validate bool
		Error    string
	}{
		{`
		s.int = false
		`, false, "not a number"},
		{`
		s.int = 42
		`, false, ""},
		{`
		s.bool = 42
		`, false, "not a bool"},
		{`
		s.bool = true
		`, false, ""},
		{`
		s.string = true
		`, false, "not a string"},
		{`
		s.string = ""
		`, false, ""},
		{`
		s.intMap = 42
		`, false, "not a table"},
		{`
		s.intMap = {[1]=2}
		`, false, "not a string key"},
		{`
		s.intMap = {a=2}
		`, false, ""},
		{`
		s.oint = {a=2}
		`, false, "not a number"},
		{`
		s.oint = nil
		`, false, ""},
		{`
		s.oint = 42
		`, false, ""},
		{`
		s{bool=true,int = 42,string="",intMap={}}
		`, true, ""},
		{`
		s{bool=true,int = 42,string="",intMap={}}
		s(s)
		`, true, ""},
		{`
		s.bool = true
		`, true, "nil field"},
	} {
		t.Run(tc.Code, func(t *testing.T) {
			s := NewStruct(checker)
			code, err := Compile(fmt.Sprintf(`
			local myfunc = function(s)
				%s
			end`, tc.Code), "input")
			assert.ErrorIs(t, err, nil)
			fn := LookupLocalFunction(code, "myfunc")
			L.Push(L.NewFunctionFromProto(fn))
			L.Push(s.AsLValue(L, false))
			checkErr := func(err error) {
				if tc.Error != "" {
					assert.ErrorContains(t, err, tc.Error)
				} else {
					assert.ErrorIs(t, err, nil)
				}
			}
			err1 := L.PCall(1, 0, nil)
			if tc.Validate {
				assert.ErrorIs(t, err1, nil)
				checkErr(s.IsIntegral())
			} else {
				checkErr(err1)
			}
		})
	}
}
