package ref

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strconv"
	"unsafe"
)

type FileHash struct {
	_filehash struct{}
	raw       string
}

func (h FileHash) IsZero() bool { return h.raw == "" }

func (h FileHash) FmtHumanize() string {
	return hex.EncodeToString(h.Raw())[:HUMANIZED_WIDTH]
}

func (h FileHash) FmtString() string {
	return hex.EncodeToString(h.Raw())
}

func (h FileHash) Sizeof() int {
	return sha256.Size
}

func (h FileHash) WithName(name string) string {
	return h.FmtString() + "-" + name
}

func (h FileHash) Raw() []byte {
	return unsafe.Slice(unsafe.StringData(h.raw), len(h.raw))
}

func (h FileHash) RawString() string {
	return h.raw
}

func (h FileHash) fromRaw(p []byte) (FileHash, error) {
	sz := h.Sizeof()
	if sz != len(p) {
		panic("file hash too short")
	}
	return FileHash{raw: unsafe.String(unsafe.SliceData(p), len(p))}, nil
}

func FileHashFromReader(r io.Reader) (FileHash, error) {
	hasher := sha256.New()
	_, err := io.Copy(hasher, r)
	if err != nil {
		return FileHash{}, err
	}
	hashVal := make([]byte, 0, FileHash{}.Sizeof())
	hashVal = hasher.Sum(hashVal)
	if len(hashVal) != (FileHash{}).Sizeof() {
		panic("bad hash result, len=" + strconv.Itoa(len(hashVal)))
	}
	return FileHash{raw: string(hashVal)}, nil
}
