package ref

import (
	"crypto/sha256"
	"encoding/hex"
	"unsafe"

	"github.com/hsfzxjy/imbed/util/fastbuf"
)

type Sha256 struct {
	_sha256_hash struct{}
	raw          string
}

func (h Sha256) IsZero() bool {
	return h.raw == ""
}

func (h Sha256) Sizeof() int {
	return sha256.Size
}

func (h Sha256) Raw() []byte {
	return unsafe.Slice(unsafe.StringData(h.raw), len(h.raw))
}

func (h Sha256) RawString() string {
	return h.raw
}

func (h Sha256) fromRaw(p []byte) (Sha256, error) {
	if len(p) != sha256.Size {
		panic("sha256 hash too short")
	}
	h.raw = string(p[:sha256.Size])
	return h, nil
}

func (h Sha256) FmtHumanize() string {
	if h.IsZero() {
		return "<none>"
	}
	return hex.EncodeToString(h.Raw())[:HUMANIZED_WIDTH]
}

func (h Sha256) FmtString() string {
	if h.IsZero() {
		return "<none>"
	}
	return hex.EncodeToString(h.Raw())
}

func Sha256HashSum(p []byte) Sha256 {
	s := sha256.Sum256(p)
	return Sha256{raw: string(s[:])}
}

func HashEncodeFunc2[T any](encodeF func(w *fastbuf.W, source T), value T) (Sha256, []byte) {
	var w fastbuf.W
	encodeF(&w, value)
	encoded := w.Result()
	var h Sha256
	var sum = sha256.Sum256(encoded)
	h.raw = string(sum[:])
	return h, encoded
}

func HashEncodeFunc(encodeF func(w *fastbuf.W)) (Sha256, []byte) {
	var w fastbuf.W
	encodeF(&w)
	encoded := w.Result()
	var h Sha256
	var sum = sha256.Sum256(encoded)
	h.raw = string(sum[:])
	return h, encoded
}
