package lualib_test

import (
	"strings"
	"testing"

	"github.com/hsfzxjy/imbed/lualib"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

func compile(code string) *lua.FunctionProto {
	reader := strings.NewReader(code)
	chunk, err := parse.Parse(reader, "input")
	if err != nil {
		panic(err)
	}
	proto, err := lua.Compile(chunk, "input")
	if err != nil {
		panic(err)
	}
	return proto
}

const code = `
local D = 42
local f = { a = 12, b = false }

local function add(a, b)
	local c = a
	local add2 = function(aa, bb)
		D = D + 1
		return D + c + aa + bb
	end

	return add2(a, b)
end

print(add(1, 2))`

// Same code yet with different whitespaces
var code2 = strings.Replace(code, "\n", "\n\n\t", 0)

func Test_Serde(t *testing.T) {
	fn := compile(code)
	s := lualib.Serialize(fn, nil).Full
	fn2, err := lualib.Deserialize(s)
	assert.ErrorIs(t, err, nil)
	assert.Equal(t, fn, fn2)
}

func Test_Serde_SameByteCode(t *testing.T) {
	s1 := lualib.Serialize(compile(code), nil).Code
	s2 := lualib.Serialize(compile(code2), nil).Code
	assert.Equal(t, s1, s2)
}

func Test_LookupLocalFunction(t *testing.T) {
	const code = `
	xxx = 0
	local yyy = 1
	local function add(a)
		if d == nil then
			d = 1
		end
		return a + d
	end
	`
	fn := compile(code)
	assert.Equal(t, fn.FunctionPrototypes[0], lualib.LookupLocalFunction(fn, "add"))
	assert.Nil(t, lualib.LookupLocalFunction(fn, "add2"))

	state := lua.NewState()
	f := state.NewFunctionFromProto(lualib.LookupLocalFunction(fn, "add"))
	state.Push(f)
	state.Push(lua.LNumber(1))
	state.PCall(1, 1, nil)
	assert.Equal(t, lua.LNumber(2), state.Get(-1))
	state.Pop(1)
	assert.Equal(t, lua.LNil, state.Get(-1))
	state.Push(f)
	state.Push(lua.LNumber(1))
	state.PCall(1, 1, nil)
	assert.Equal(t, lua.LNumber(2), state.Get(-1))
}
