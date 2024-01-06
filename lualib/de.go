package lualib

import (
	"fmt"

	lua "github.com/hsfzxjy/gopher-lua"
	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type deError struct {
	index  int
	detail string
	err    error
}

func (e *deError) Error() string {
	return fmt.Sprintf("lualib.Deserialize: error at %d-th proto: %s: %p", e.index, e.detail, e.err)
}

func (e *deError) Unwrap() error {
	return e.err
}

type deserializer struct {
	protos []*luaFunctionProto
	r      fastbuf.R
	dbgR   fastbuf.R
}

func (d *deserializer) doFunctionProto() error {
	proto := &luaFunctionProto{}
	d.protos = append(d.protos, proto)
	index := len(d.protos) - 1
	var err error
	proto.NumUpvalues, err = d.r.ReadUint8()
	if err != nil {
		return &deError{index: index, detail: "NumUpvalues", err: err}
	}
	proto.NumParameters, err = d.r.ReadUint8()
	if err != nil {
		return &deError{index: index, detail: "NumParameters", err: err}
	}
	proto.IsVarArg, err = d.r.ReadUint8()
	if err != nil {
		return &deError{index: index, detail: "IsVarArg", err: err}
	}
	proto.NumUsedRegisters, err = d.r.ReadUint8()
	if err != nil {
		return &deError{index: index, detail: "NumUsedRegisters", err: err}
	}
	proto.Code, err = d.r.ReadUint32Array()
	if err != nil {
		return &deError{index: index, detail: "Code", err: err}
	}
	protoIndices, err := d.r.ReadUint32Array()
	if err != nil {
		return &deError{index: index, detail: "protoIndices", err: err}
	}
	proto.FunctionPrototypes = make([]*lua.FunctionProto, len(protoIndices))
	for i, idx := range protoIndices {
		if idx >= uint32(index) {
			return &deError{index: index, detail: "protoIndices", err: fmt.Errorf("invalid index")}
		}
		proto.FunctionPrototypes[i] = d.protos[idx].asLua()
	}
	constantsLen, err := d.r.ReadArrayHeader()
	if err != nil {
		return &deError{index: index, detail: "constantsLen", err: err}
	}
	proto.Constants = make([]lua.LValue, constantsLen)
	for i := range proto.Constants {
		kind, err := d.r.ReadUint8()
		if err != nil {
			return &deError{index: index, detail: "constant kind", err: err}
		}
		switch kind {
		case lconstNumber:
			f, err := d.r.ReadFloat64()
			if err != nil {
				return &deError{index: index, detail: "constant number", err: err}
			}
			proto.Constants[i] = lua.LNumber(f).AsLValue()
		case lconstString:
			s, err := d.r.ReadString()
			if err != nil {
				return &deError{index: index, detail: "constant string", err: err}
			}
			proto.Constants[i] = lua.LString(s).AsLValue()
		default:
			return &deError{index: index, detail: "constant kind", err: fmt.Errorf("unknown kind %d", kind)}
		}
	}
	proto.SourceName, err = d.dbgR.ReadString()
	if err != nil {
		return &deError{index: index, detail: "SourceName", err: err}
	}
	proto.LineDefined, err = d.dbgR.ReadInt()
	if err != nil {
		return &deError{index: index, detail: "LineDefined", err: err}
	}
	proto.LastLineDefined, err = d.dbgR.ReadInt()
	if err != nil {
		return &deError{index: index, detail: "LastLineDefined", err: err}
	}
	proto.DbgSourcePositions, err = d.dbgR.ReadIntArray()
	if err != nil {
		return &deError{index: index, detail: "DbgSourcePositions", err: err}
	}
	localsLen, err := d.dbgR.ReadArrayHeader()
	if err != nil {
		return &deError{index: index, detail: "localsLen", err: err}
	}
	proto.DbgLocals = make([]*lua.DbgLocalInfo, localsLen)
	for i := range proto.DbgLocals {
		proto.DbgLocals[i] = &lua.DbgLocalInfo{}

		proto.DbgLocals[i].Name, err = d.dbgR.ReadString()
		if err != nil {
			return &deError{index: index, detail: "DbgLocals.Name", err: err}
		}
		proto.DbgLocals[i].StartPc, err = d.dbgR.ReadInt()
		if err != nil {
			return &deError{index: index, detail: "DbgLocals.StartPc", err: err}
		}
		proto.DbgLocals[i].EndPc, err = d.dbgR.ReadInt()
		if err != nil {
			return &deError{index: index, detail: "DbgLocals.EndPc", err: err}
		}
	}
	callsLen, err := d.dbgR.ReadArrayHeader()
	if err != nil {
		return &deError{index: index, detail: "callsLen", err: err}
	}
	proto.DbgCalls = make([]lua.DbgCall, callsLen)
	for i := range proto.DbgCalls {
		c := &proto.DbgCalls[i]
		c.Name, err = d.dbgR.ReadString()
		if err != nil {
			return &deError{index: index, detail: "DbgCalls.Name", err: err}
		}
		c.Pc, err = d.dbgR.ReadInt()
		if err != nil {
			return &deError{index: index, detail: "DbgCalls.Pc", err: err}
		}
	}
	upvaluesLen, err := d.dbgR.ReadArrayHeader()
	if err != nil {
		return &deError{index: index, detail: "upvaluesLen", err: err}
	}
	proto.DbgUpvalues = make([]string, upvaluesLen)
	for i := range proto.DbgUpvalues {
		proto.DbgUpvalues[i], err = d.dbgR.ReadString()
		if err != nil {
			return &deError{index: index, detail: "DbgUpvalues", err: err}
		}
	}
	proto.stringConstants = make([]string, len(proto.Constants))
	for i, c := range proto.Constants {
		sv := ""
		if slv, ok := c.AsLString(); ok {
			sv = string(slv)
		}
		proto.stringConstants[i] = sv
	}
	return nil
}

func Deserialize(full []byte) (*lua.FunctionProto, error) {
	var r = fastbuf.R{Buf: full}
	n, err := r.ReadUsize()
	if err != nil {
		return nil, err
	}
	r1, r2, err := r.SplitAt(n)
	if err != nil {
		return nil, err
	}
	var de = deserializer{
		protos: nil,
		r:      r1,
		dbgR:   r2,
	}
	for !de.r.EOF() && !de.dbgR.EOF() {
		if err := de.doFunctionProto(); err != nil {
			return nil, err
		}
	}
	if !de.r.EOF() || !de.dbgR.EOF() {
		return nil, fmt.Errorf("lualib.Deserialize: unexpected EOF")
	}
	return de.protos[len(de.protos)-1].asLua(), nil
}
