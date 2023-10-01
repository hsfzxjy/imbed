package util

import (
	"bytes"
	"encoding/hex"
)

type ICBytes struct {
	slice []byte
}

func NewICBytes(repr string) (ICBytes, error) {
	n := len(repr)
	var reprSlice []byte
	if n%2 == 0 {
		m := n + 2
		reprSlice = make([]byte, m)
		copy(reprSlice[:n], repr[:])
		reprSlice[n] = '0'
		reprSlice[n+1] = '0'
	} else {
		m := n + 1
		reprSlice = make([]byte, m)
		copy(reprSlice[:n], repr[:])
		reprSlice[n] = '1'
	}
	buf := make([]byte, len(reprSlice)/2)
	_, err := hex.Decode(buf, reprSlice)
	if err != nil {
		return ICBytes{}, err
	}
	return ICBytes{buf}, nil
}

func (i ICBytes) PrefixMatch(target []byte) bool {
	n := len(i.slice)
	last := i.slice[n-1]
	if last == 0 {
		return bytes.HasPrefix(target, i.slice[:n-1])
	} else {
		if len(target) < n {
			return false
		}
		if !bytes.HasPrefix(target, i.slice[:n-1]) {
			return false
		}
		return target[n-1]>>4 == last>>4
	}
}

func (i ICBytes) FullMatch(target []byte) bool {
	n := len(i.slice)
	last := i.slice[n-1]
	if last == 0 {
		return bytes.Equal(target, i.slice[:n-1])
	}
	return false
}

func (i ICBytes) Bytes() []byte {
	return i.slice[:len(i.slice)-1]
}
