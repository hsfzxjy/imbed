package transform

import (
	"bytes"
	"sync"

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

func newEncodable(value schema.GenericValue) *encodable {
	x := new(encodable)
	x.Compute = func() {
		var buf bytes.Buffer
		var w = msgp.NewWriter(&buf)
		err := value.EncodeMsg(w)
		if err != nil {
			x.encodeError = err
			return
		}
		w.Flush()
		x.encoded = buf.Bytes()
		x.hash = ref.Sha256HashSum(buf.Bytes())
	}
	return x
}
