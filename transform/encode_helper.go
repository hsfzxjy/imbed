package transform

import (
	"sync"

	"github.com/hsfzxjy/imbed/core/ref"
)

type encodable struct {
	computeOnce sync.Once
	encodeError error
	encoded     []byte
	hash        ref.Sha256Hash
	Compute     func()
}

func (e *encodable) GetRawEncoded() ([]byte, error) {
	e.computeOnce.Do(e.Compute)
	return e.encoded, e.encodeError
}

func (e *encodable) GetSha256Hash() (ref.Sha256Hash, error) {
	e.computeOnce.Do(e.Compute)
	return e.hash, e.encodeError
}
