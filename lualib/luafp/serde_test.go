package luafp_test

import (
	"strings"
	"testing"

	lua "github.com/hsfzxjy/gopher-lua"
	"github.com/hsfzxjy/imbed/lualib/luafp"
	"github.com/stretchr/testify/assert"
)

func compile(code string) *lua.FunctionProto {
	fn, err := luafp.Compile(code, "test")
	if err != nil {
		panic(err)
	}
	return fn
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
	s := luafp.Serialize(fn, nil).Full
	fn2, err := luafp.Deserialize(s)
	assert.ErrorIs(t, err, nil)
	assert.NotNil(t, fn2)
}

func Test_Serde_SameByteCode(t *testing.T) {
	s1 := luafp.Serialize(compile(code), nil).Code
	s2 := luafp.Serialize(compile(code2), nil).Code
	assert.Equal(t, s1, s2)
}
