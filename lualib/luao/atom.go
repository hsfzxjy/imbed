package luao

import lua "github.com/hsfzxjy/gopher-lua"

type IntP struct{}

func isInt64(v lua.LNumber) bool {
	return v == lua.LNumber(int64(v))
}

func (IntP) L2G(v lua.LValue) (Object, error) {
	n, ok := v.AsLNumber()
	if !ok {
		return Object{}, &typeError{lua.LTNumber, v.Type()}
	}
	if !isInt64(n) {
		return Object{}, &intError{n}
	}
	return Object{v}, nil
}

func (IntP) checkIntegrity(obj Object) error {
	if obj.Type() != lua.LTNumber {
		return &typeError{lua.LTNumber, obj.Type()}
	}
	return nil
}

type StringP struct{}

func (StringP) L2G(v lua.LValue) (Object, error) {
	if typ := v.Type(); typ != lua.LTString {
		return Object{}, &typeError{lua.LTString, typ}
	}
	return Object{v}, nil
}

func (StringP) checkIntegrity(obj Object) error {
	if obj.Type() != lua.LTString {
		return &typeError{lua.LTString, obj.Type()}
	}
	return nil
}

type BoolP struct{}

func (BoolP) L2G(v lua.LValue) (Object, error) {
	if typ := v.Type(); typ != lua.LTBool {
		return Object{}, &typeError{lua.LTBool, typ}
	}
	return Object{v}, nil
}

func (BoolP) checkIntegrity(obj Object) error {
	if obj.Type() != lua.LTBool {
		return &typeError{lua.LTBool, obj.Type()}
	}
	return nil
}

type OptP struct{ Protocol }

func (p *OptP) L2G(v lua.LValue) (Object, error) {
	if v.Type() == lua.LTNil {
		return Object{lua.LNil}, nil
	}
	return p.Protocol.L2G(v)
}

func (p *OptP) checkIntegrity(obj Object) error {
	if obj.Type() == lua.LTNil {
		return nil
	}
	return p.Protocol.checkIntegrity(obj)
}
