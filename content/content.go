package content

import (
	"bytes"
	"errors"
	_ "image/jpeg"
	_ "image/png"
	"sync/atomic"

	"sync"

	ref "github.com/hsfzxjy/imbed/core/ref"
)

type content struct {
	loader Loader
	sizer  Sizer
	buf    []byte
	hash   ref.Murmur3Hash

	computed      atomic.Bool
	computedBuf   []byte
	computedHash  ref.Murmur3Hash
	computedError error

	once sync.Once
}

func (c *content) GetHash() (ref.Murmur3Hash, error) {
	if !c.hash.IsZero() {
		return c.hash, nil
	}
	c.load()
	return c.computedHash, c.computedError
}

func (c *content) Size() (Size, error) {
	switch {
	case c.computed.Load():
		return Size(len(c.computedBuf)), c.computedError
	case c.buf != nil:
		return Size(len(c.buf)), nil
	case c.sizer != nil:
		return c.sizer.Size()
	}
	c.load()
	return Size(len(c.computedBuf)), nil
}

func (c *content) BytesReader() (*bytes.Reader, error) {
	switch {
	case c.buf != nil:
		return bytes.NewReader(c.buf), nil
	}
	c.load()
	if c.computedError != nil {
		return nil, c.computedError
	}
	return bytes.NewReader(c.computedBuf), nil
}

func (c *content) load() {
	c.once.Do(func() {
		var (
			buffer bytes.Buffer
			buf    []byte
			hash   ref.Murmur3Hash
			err    error
		)
		defer func() {
			c.computedError = err
			c.computed.Store(true)
		}()
		if c.loader == nil {
			c.computedBuf = c.buf
			return
		}
		err = c.loader.Load(&buffer)
		if err != nil {
			return
		}
		buf = buffer.Bytes()
		hash, err = ref.Murmur3HashFromReader(bytes.NewReader(buf))
		if err != nil {
			return
		}

		if !c.hash.IsZero() && hash != c.hash {
			err = errors.New("hash from loaded file differs from given one")
			return
		}
		c.computedBuf = buf
		c.computedHash = hash
	})
}
