package ref

import (
	"encoding/hex"
	"io"
	"strconv"
	"unsafe"

	"github.com/spaolacci/murmur3"
)

type Murmur3 struct {
	_murmur3 struct{}
	raw      string
}

func (h Murmur3) IsZero() bool { return h.raw == "" }

func (h Murmur3) FmtHumanize() string {
	return hex.EncodeToString(h.Raw())[:HUMANIZED_WIDTH]
}

func (h Murmur3) FmtString() string {
	return hex.EncodeToString(h.Raw())
}

func (h Murmur3) Sizeof() int {
	return 128 / 8
}

func (h Murmur3) WithName(name string) string {
	return h.FmtString() + "-" + name
}

func (h Murmur3) Raw() []byte {
	return unsafe.Slice(unsafe.StringData(h.raw), len(h.raw))
}

func (h Murmur3) RawString() string {
	return h.raw
}

func (h Murmur3) fromRaw(p []byte) (Murmur3, error) {
	sz := h.Sizeof()
	if sz != len(p) {
		panic("murmur3 hash too short")
	}
	return Murmur3{raw: unsafe.String(unsafe.SliceData(p), len(p))}, nil
}

func Murmur3FromReader(r io.Reader) (Murmur3, error) {
	hasher := murmur3.New128()
	_, err := io.Copy(hasher, r)
	if err != nil {
		return Murmur3{}, err
	}
	hashVal := make([]byte, 0, Murmur3{}.Sizeof())
	hashVal = hasher.Sum(hashVal)
	if len(hashVal) != (Murmur3{}).Sizeof() {
		panic("bad hash result, len=" + strconv.Itoa(len(hashVal)))
	}
	return Murmur3{raw: string(hashVal)}, nil
}
