package ndl

import (
	"bytes"
	"unsafe"

	"github.com/hsfzxjy/imbed/util"
)

type Needle interface {
	Match(key []byte) bool
	Bytes() []byte
}

type hex struct {
	util.HalfBytes
	matchFunc func(ic util.HalfBytes, target []byte) bool
}

func (n hex) Match(key []byte) bool {
	return n.matchFunc(n.HalfBytes, key)
}

func HexPrefix(prefixPretty string) (Needle, error) {
	ic, err := util.NewHalfBytes(prefixPretty)
	if err != nil {
		return nil, err
	}
	return hex{ic, util.HalfBytes.PrefixMatch}, nil
}

func HexFull(pretty string) (Needle, error) {
	ic, err := util.NewHalfBytes(pretty)
	if err != nil {
		return nil, err
	}
	return hex{ic, util.HalfBytes.FullMatch}, nil
}

func Hex(pretty string, prefix bool) (Needle, error) {
	if prefix {
		return HexPrefix(pretty)
	} else {
		return HexFull(pretty)
	}
}

type raw struct {
	string
	matchFunc func(bytes, target []byte) bool
}

func (n raw) Match(key []byte) bool {
	return n.matchFunc(key, n.Bytes())
}

func (n raw) Bytes() []byte {
	return unsafe.Slice(unsafe.StringData(n.string), len(n.string))
}

func RawPrefix(str string) Needle {
	return raw{str, bytes.HasPrefix}
}

func RawFull(str string) Needle {
	return raw{str, bytes.Equal}
}

func Raw(str string, prefix bool) (Needle, error) {
	if prefix {
		return RawPrefix(str), nil
	} else {
		return RawFull(str), nil
	}
}
