package core

import (
	"bytes"
	"unsafe"

	"github.com/hsfzxjy/imbed/util"
)

type Needle interface {
	Match(key []byte) bool
	Bytes() []byte
}

type bytesNeedle struct {
	util.ICBytes
	matchFunc func(ic util.ICBytes, target []byte) bool
}

func (n bytesNeedle) Match(key []byte) bool {
	return n.matchFunc(n.ICBytes, key)
}

func BytesPrefix(prefixPretty string) (Needle, error) {
	ic, err := util.NewICBytes(prefixPretty)
	if err != nil {
		return nil, err
	}
	return bytesNeedle{ic, util.ICBytes.PrefixMatch}, nil
}

func BytesFull(pretty string) (Needle, error) {
	ic, err := util.NewICBytes(pretty)
	if err != nil {
		return nil, err
	}
	return bytesNeedle{ic, util.ICBytes.FullMatch}, nil
}

func BytesNeedle(pretty string, prefix bool) (Needle, error) {
	if prefix {
		return BytesPrefix(pretty)
	} else {
		return BytesFull(pretty)
	}
}

type stringNeedle struct {
	string
	matchFunc func(bytes, target []byte) bool
}

func (n stringNeedle) Match(key []byte) bool {
	return n.matchFunc(key, n.Bytes())
}

func (n stringNeedle) Bytes() []byte {
	return unsafe.Slice(unsafe.StringData(n.string), len(n.string))
}

func StringPrefix(str string) Needle {
	return stringNeedle{str, bytes.HasPrefix}
}

func StringFull(str string) Needle {
	return stringNeedle{str, bytes.Equal}
}

func StringNeedle(str string, prefix bool) Needle {
	if prefix {
		return StringPrefix(str)
	} else {
		return StringFull(str)
	}
}
