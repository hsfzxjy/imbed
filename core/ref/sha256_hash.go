package ref

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"sync"

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
	return hex.EncodeToString(AsRaw(h))[:7]
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

type selfEncoder interface {
	EncodeSelf(w *msgp.Writer) error
}

type EncodableObject interface {
	GetSha256Hash() (Sha256Hash, error)
	GetRawEncoded() ([]byte, error)
}

type LazyEncodableObject[T selfEncoder] struct {
	Inner   T
	once    sync.Once
	hash    Sha256Hash
	encoded []byte
}

func (v *LazyEncodableObject[T]) compute() error {
	var err error
	v.once.Do(func() {
		v.hash, v.encoded, err = HashEncodeFunc(v.Inner.EncodeSelf)
	})
	return err
}

func (v *LazyEncodableObject[T]) GetSha256Hash() (Sha256Hash, error) {
	if err := v.compute(); err != nil {
		return Sha256Hash{}, err
	}
	return v.hash, nil
}

func (v *LazyEncodableObject[T]) GetRawEncoded() ([]byte, error) {
	if err := v.compute(); err != nil {
		return nil, err
	}
	return v.encoded, nil
}

type ConstantEncodableObject struct {
	once    sync.Once
	hash    Sha256Hash
	Encoded []byte
}

func (v *ConstantEncodableObject) GetSha256Hash() (Sha256Hash, error) {
	v.once.Do(func() {
		var sum = sha256.Sum256(v.Encoded)
		v.hash.raw = string(sum[:])
	})
	return v.hash, nil
}

func (v *ConstantEncodableObject) GetRawEncoded() ([]byte, error) {
	return v.Encoded, nil
}
