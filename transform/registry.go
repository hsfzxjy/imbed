package transform

import (
	"bytes"
	"fmt"

	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/tinylib/msgp/msgp"
)

type Registry struct {
	metadataTable map[string]Metadata
}

var defaultRegistry = NewRegistry()

func DefaultRegistry() Registry { return defaultRegistry }

func NewRegistry() Registry {
	return Registry{map[string]Metadata{}}
}

func (r Registry) Lookup(name string) (Metadata, bool) {
	m, ok := r.metadataTable[name]
	return m, ok
}

func (r Registry) DecodeParams(buf []byte) (result []Builder, err error) {
	bufR := bytes.NewReader(buf)
	msgR := msgp.NewReader(bufR)
	for bufR.Len() > 0 {
		cfgHash, err := ref.FromReader[ref.Sha256Hash](bufR)
		if err != nil {
			return nil, err
		}
		name, err := msgR.ReadString()
		if err != nil {
			return nil, err
		}
		metadata, ok := r.Lookup(name)
		if !ok {
			return nil, fmt.Errorf("no transform named %q", name)
		}
		pm, err := metadata.decodeMsg(msgR)
		if err != nil {
			return nil, err
		}
		result = append(result, &builder{
			metadata.ConfigBuilderHash(cfgHash),
			pm,
		})
	}
	return result, nil
}
