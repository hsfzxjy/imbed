package ref

import (
	"encoding/hex"
	"io"
	"strconv"
	"unsafe"

	"github.com/spaolacci/murmur3"
)

type Murmur3Hash struct {
	_murmur3_hash struct{}
	raw           string
}

func (h Murmur3Hash) IsZero() bool  { return h.raw == "" }
func (h Murmur3Hash) IsValid() bool { return len(h.raw) == MURMUR3_HASH_LEN }
func (h Murmur3Hash) Len() int      { return len(h.raw) }

func (h Murmur3Hash) Bytes() []byte {
	return unsafe.Slice(unsafe.StringData(h.raw), len(h.raw))
}

func (h Murmur3Hash) Humanize() string {
	return hex.EncodeToString(h.Bytes())
}

func (h Murmur3Hash) FmtHumanize() string {
	return hex.EncodeToString(h.Bytes())[:HUMANIZED_WIDTH]
}

func (h Murmur3Hash) FmtString() string {
	return hex.EncodeToString(h.Bytes())
}

const MURMUR3_HASH_LEN = 128 / 8
const MURMUR3_HASH_HUMANIZE_LEN = 2 * MURMUR3_HASH_LEN

func (h Murmur3Hash) fromBytes(p []byte) (Murmur3Hash, []byte) {
	if len(p) < MURMUR3_HASH_LEN {
		panic("murmur3 hash too short")
	}
	h.raw = string(p[:MURMUR3_HASH_LEN])
	return h, p[MURMUR3_HASH_LEN:]
}

func Murmur3HashFromReader(r io.Reader) (Murmur3Hash, error) {
	hasher := murmur3.New128()
	_, err := io.Copy(hasher, r)
	if err != nil {
		return Murmur3Hash{}, err
	}
	hashVal := make([]byte, 0, MURMUR3_HASH_LEN)
	hashVal = hasher.Sum(hashVal)
	if len(hashVal) != MURMUR3_HASH_LEN {
		panic("bad hash result, len=" + strconv.Itoa(len(hashVal)))
	}
	return Murmur3Hash{
		raw: string(hashVal),
	}, nil
}
