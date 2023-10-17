package ref

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"

	"github.com/tinylib/msgp/msgp"
)

type Sha256Hash struct {
	raw          string
	_sha256_hash struct{}
}

func (h sha256hash) IsZero() bool      { return h.raw == "" }
func (h sha256hash) Len() int          { return len(h.raw) }
func (h sha256hash) expectedSize() int { return sha256.Size }
func (h sha256hash) fromBytes(p []byte) (sha256hash, []byte) {
	if len(p) < sha256.Size {
		panic("sha256 hash too short")
	}
	h.raw = string(p[:sha256.Size])
	return h, p[sha256.Size:]
}
func (h sha256hash) FmtHumanize() string {
	if h.IsZero() {
		return "<none>"
	}
	return hex.EncodeToString(AsRaw(h))[:HUMANIZED_WIDTH]
}
func (h sha256hash) FmtString() string {
	if h.IsZero() {
		return "<none>"
	}
	return hex.EncodeToString(AsRaw(h))
}

func Sha256HashSum(p []byte) Sha256Hash {
	s := sha256.Sum256(p)
	return Sha256Hash{raw: string(s[:])}
}

func HashEncodeFunc2[T any](encodeF func(w *msgp.Writer, source T) error, value T) (Sha256Hash, []byte, error) {
	return HashEncodeFunc(func(w *msgp.Writer) error {
		return encodeF(w, value)
	})
}

func HashEncodeFunc(encodeF func(w *msgp.Writer) error) (Sha256Hash, []byte, error) {
	var buf bytes.Buffer
	w := msgp.NewWriter(&buf)
	err := encodeF(w)
	if err != nil {
		return Sha256Hash{}, nil, err
	}
	err = w.Flush()
	if err != nil {
		return Sha256Hash{}, nil, err
	}
	encoded := buf.Bytes()
	var h Sha256Hash
	var sum = sha256.Sum256(encoded)
	h.raw = string(sum[:])
	return h, encoded, nil
}

type EncodableObject interface {
	GetSha256Hash() (Sha256Hash, error)
	GetRawEncoded() ([]byte, error)
}
