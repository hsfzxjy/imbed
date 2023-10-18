package transform

import (
	"errors"
	"sync"
	"unsafe"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
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

// Unsafe, but neat.
type EncodeMsgHelper[T any] struct{}

func (h *EncodeMsgHelper[T]) EncodeMsg(r *Registry, w *msgp.Writer) error {
	v := (*T)(unsafe.Pointer(h))
	sch, err := r.schemaStore.Get(v)
	if err != nil {
		if errors.Is(err, schema.ErrNotRegistered) {
			sch, err = schema.Register[T](r.schemaStore)
		}
		if err != nil {
			return err
		}
	}
	return sch.EncodeMsgAny(w, v)
}
