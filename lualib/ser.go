package lualib

import (
	"github.com/hsfzxjy/imbed/util/fastbuf"
	lua "github.com/yuin/gopher-lua"
)

type serializer struct {
	nProtos uint32
	codes   fastbuf.W
	dbgs    fastbuf.W
}

func (s *serializer) doFunctionProto(proto *lua.FunctionProto) uint32 {
	protoIndices := make([]uint32, len(proto.FunctionPrototypes))
	for i, p := range proto.FunctionPrototypes {
		protoIndices[i] = s.doFunctionProto(p)
	}
	s.codes.
		WriteUint8(proto.NumUpvalues).
		WriteUint8(proto.NumParameters).
		WriteUint8(proto.IsVarArg).
		WriteUint8(proto.NumUsedRegisters).
		WriteUint32Array(proto.Code).
		WriteUint32Array(protoIndices)
	s.codes.WriteArrayHeader(uint32(len(proto.Constants)))
	for _, c := range proto.Constants {
		switch v := c.(type) {
		case lua.LString:
			s.codes.WriteUint8(lconstString).WriteString(string(v))
		case lua.LNumber:
			s.codes.WriteUint8(lconstNumber).WriteFloat64(float64(v))
		default:
			panic("lualib: unknown constant type " + v.Type().String())
		}
	}

	s.dbgs.
		WriteString(proto.SourceName).
		WriteInt(proto.LineDefined).
		WriteInt(proto.LastLineDefined).
		WriteIntArray(proto.DbgSourcePositions)
	s.dbgs.WriteArrayHeader(uint32(len(proto.DbgLocals)))
	for _, l := range proto.DbgLocals {
		s.dbgs.
			WriteString(l.Name).
			WriteInt(l.StartPc).
			WriteInt(l.EndPc)
	}
	s.dbgs.WriteArrayHeader(uint32(len(proto.DbgCalls)))
	for _, c := range proto.DbgCalls {
		s.dbgs.
			WriteString(c.Name).
			WriteInt(c.Pc)
	}
	s.dbgs.WriteArrayHeader(uint32(len(proto.DbgUpvalues)))
	for _, u := range proto.DbgUpvalues {
		s.dbgs.WriteString(u)
	}
	index := s.nProtos
	s.nProtos++
	return index
}

func Serialize(proto *lua.FunctionProto, sourceContent []byte) struct{ Full, Code []byte } {
	var s serializer
	s.doFunctionProto(proto)
	codes, dbgs := s.codes.Result(), s.dbgs.Result()
	var sz fastbuf.Size
	w := sz.
		ReserveBytes(dbgs).
		ReserveBytes(codes).
		Reserve(8).
		Build()
	w.
		AppendUsize(uint64(len(codes))).
		AppendRaw(codes).
		AppendRaw(dbgs)
	r := w.Result()
	return struct{ Full, Code []byte }{
		Full: r,
		Code: r[8 : 8+len(codes)],
	}
}
