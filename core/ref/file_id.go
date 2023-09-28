package ref

import (
	"encoding/hex"

	"github.com/hsfzxjy/imbed/util"
)

// Layout:
// basename (var bytes) + murmur3hash (16 bytes)
type FID struct {
	raw  string
	_fid struct{}
}

func (id FID) IsZero() bool {
	return id.raw == ""
}

func (id FID) Len() int { return len(id.raw) }

func (id FID) Basename() string {
	if id.IsZero() {
		return ""
	}
	return string(id.raw[:len(id.raw)-MURMUR3_HASH_LEN])
}

func (id FID) Hash() Murmur3Hash {
	if id.IsZero() {
		return Murmur3Hash{}
	}
	return Murmur3Hash{raw: id.raw[len(id.raw)-MURMUR3_HASH_LEN:]}
}

func (id FID) Humanize() string {
	if id.IsZero() {
		return ""
	}
	return id.Hash().Humanize() + "-" + id.Basename()
}

func (id FID) fromBytes(p []byte) (FID, []byte) {
	if len(p) < MURMUR3_HASH_LEN+1 {
		panic("fid too short")
	}
	id.raw = string(p)
	return id, nil
}

func FIDFromParts(basename string, hash Murmur3Hash) FID {
	if !hash.IsValid() {
		panic("invalid hash")
	}
	return FID{raw: basename + hash.raw}
}

func FIDFromPretty(input string) FID {
	if len(input) <= 1+MURMUR3_HASH_HUMANIZE_LEN {
		panic("input too short")
	}
	if input[MURMUR3_HASH_HUMANIZE_LEN] != '-' {
		panic("malformed input")
	}
	hash, err := hex.DecodeString(input[:MURMUR3_HASH_HUMANIZE_LEN])
	util.Check(err)
	return FID{raw: input[MURMUR3_HASH_HUMANIZE_LEN+1:] + string(hash)}
}
