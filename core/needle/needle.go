package ndl

import (
	"bytes"
	"unsafe"

	"github.com/hsfzxjy/imbed/core/pos"
	"github.com/hsfzxjy/imbed/util"
)

type Needle interface {
	Match(key []byte) bool
	Bytes() []byte
	pos.PosGetter
}

type hex struct {
	util.HalfBytes
	matchFunc func(ic util.HalfBytes, target []byte) bool
	pos.P
}

func (n *hex) Match(key []byte) bool {
	return n.matchFunc(n.HalfBytes, key)
}

func HexPrefix(prefixPretty string, pos pos.P) (Needle, error) {
	ic, err := util.NewHalfBytes(prefixPretty)
	if err != nil {
		return nil, err
	}
	return &hex{ic, util.HalfBytes.PrefixMatch, pos}, nil
}

func HexFull(pretty string, pos pos.P) (Needle, error) {
	ic, err := util.NewHalfBytes(pretty)
	if err != nil {
		return nil, err
	}
	return &hex{ic, util.HalfBytes.FullMatch, pos}, nil
}

func Hex(pretty string, prefix bool, pos pos.P) (Needle, error) {
	if prefix {
		return HexPrefix(pretty, pos)
	} else {
		return HexFull(pretty, pos)
	}
}

type raw struct {
	string
	matchFunc func(bytes, target []byte) bool
	pos.P
}

func (n *raw) Match(key []byte) bool {
	return n.matchFunc(key, n.Bytes())
}

func (n *raw) Bytes() []byte {
	return unsafe.Slice(unsafe.StringData(n.string), len(n.string))
}

func RawPrefix(str string, pos pos.P) Needle {
	return &raw{str, bytes.HasPrefix, pos}
}

func RawFull(str string, pos pos.P) Needle {
	return &raw{str, bytes.Equal, pos}
}

func Raw(str string, prefix bool, pos pos.P) (Needle, error) {
	if prefix {
		return RawPrefix(str, pos), nil
	} else {
		return RawFull(str, pos), nil
	}
}
