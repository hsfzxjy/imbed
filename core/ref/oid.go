package ref

import (
	"encoding/binary"
	"strconv"
	"unsafe"
)

type OID struct {
	_oid struct{}
	num  uint64
}

func (id OID) Uint64() uint64 {
	return id.num
}

func (id OID) IsZero() bool {
	return id.num == 0
}

func (id OID) Sizeof() int {
	return 8
}

func (id OID) Raw() []byte {
	var b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, id.num)
	return b
}

func (id OID) RawString() string {
	raw := id.Raw()
	return unsafe.String(unsafe.SliceData(raw), len(raw))
}

func (id OID) fromRaw(b []byte) (OID, error) {
	if len(b) != 8 {
		panic("oid too short")
	}
	return OID{num: binary.BigEndian.Uint64(b)}, nil
}

func (id OID) FmtHumanize() string {
	return id.FmtString()
}

func (id OID) FmtString() string {
	return strconv.FormatUint(id.num, 10)
}

func NewOID(num uint64) OID {
	return OID{num: num}
}
