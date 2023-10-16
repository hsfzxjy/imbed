package util

import (
	"bytes"
	"encoding/hex"
)

// HalfBytes represents a byte slice that might hold a 4-bit item at the rear.
// E.g., a hex string `deadb` with odd length could be parsed into HalfBytes.
type HalfBytes struct {
	slice []byte
}

func NewHalfBytes(repr string) (HalfBytes, error) {
	n := len(repr)
	var reprSlice []byte
	if n%2 == 0 {
		m := n + 2
		reprSlice = make([]byte, m)
		copy(reprSlice[:n], repr[:])
		reprSlice[n] = '0'
		reprSlice[n+1] = '0'
	} else {
		m := n + 3
		reprSlice = make([]byte, m)
		copy(reprSlice[:n], repr[:])
		reprSlice[n] = '0'
		reprSlice[n+1] = '0'
		reprSlice[n+2] = '1'
	}
	buf := make([]byte, len(reprSlice)/2)
	_, err := hex.Decode(buf, reprSlice)
	if err != nil {
		return HalfBytes{}, err
	}
	return HalfBytes{buf}, nil
}

func (i HalfBytes) PrefixMatch(target []byte) bool {
	n := len(i.slice)
	last := i.slice[n-1]
	if last == 0 {
		// even hex digits
		return bytes.HasPrefix(target, i.slice[:n-1])
	} else {
		// odd hex digits
		if len(target) <= n-2 {
			return false
		}
		if !bytes.HasPrefix(target, i.slice[:n-2]) {
			return false
		}
		return target[n-2]>>4 == i.slice[n-2]>>4
	}
}

func (i HalfBytes) FullMatch(target []byte) bool {
	n := len(i.slice)
	last := i.slice[n-1]
	if last == 0 {
		return bytes.Equal(target, i.slice[:n-1])
	}
	return false
}

func (i HalfBytes) Bytes() []byte {
	return i.slice[:len(i.slice)-1]
}
