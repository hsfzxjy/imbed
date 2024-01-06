package luafp_test

import (
	"testing"

	lua "github.com/hsfzxjy/gopher-lua"
	"github.com/stretchr/testify/assert"
)

func AssertFunctionProtoEquals(t *testing.T, expected *lua.FunctionProto, actual *lua.FunctionProto) {
	assert.Equal(t, expected.NumUpvalues, actual.NumUpvalues)
	assert.Equal(t, expected.NumParameters, actual.NumParameters)
	assert.Equal(t, expected.IsVarArg, actual.IsVarArg)
	assert.Equal(t, expected.NumUsedRegisters, actual.NumUsedRegisters)
	assert.Equal(t, expected.Code, actual.Code)
	for i, c := range expected.Constants {
		assert.Truef(t, c.Equals(actual.Constants[i]), "expected %v, got %v", c.String(), actual.Constants[i].String())
	}
	for i, p := range expected.FunctionPrototypes {
		AssertFunctionProtoEquals(t, p, actual.FunctionPrototypes[i])
	}
	assert.Equal(t, expected.DbgSourcePositions, actual.DbgSourcePositions)
	assert.Equal(t, expected.DbgLocals, actual.DbgLocals)
	assert.Equal(t, expected.DbgCalls, actual.DbgCalls)
	assert.Equal(t, expected.DbgUpvalues, actual.DbgUpvalues)
}
