package transform

import (
	"bytes"

	"github.com/hsfzxjy/imbed/core"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

type metadata[C any, P Params[C]] struct {
	name    string
	aliases []string

	configSchema schema.Schema[C]
	paramsSchema schema.Schema[P]

	kind Kind
}

func (m *metadata[C, P]) Parse(cs core.ConfigProvider, paramsR schema.Reader) (Transform, error) {
	cfgR, err := cs.ProvideWorkspaceConfig(m.name)
	if err != nil {
		return nil, err
	}
	var cfgInstance C
	err = m.configSchema.DecodeValue(cfgR, &cfgInstance)
	if err != nil {
		return nil, err
	}
	var paramsInstance P
	err = m.paramsSchema.DecodeValue(paramsR, &paramsInstance)
	if err != nil {
		return nil, paramsR.Error(err)
	}
	applier, err := paramsInstance.BuildTransform(&cfgInstance)
	if err != nil {
		return nil, err
	}
	return newSingleTransform(m, &cfgInstance, &paramsInstance, applier), nil
}

func (m *metadata[C, P]) decodeMsg(cs core.ConfigProvider, paramsR *msgp.Reader) (Transform, error) {
	cfgHash, err := ref.FromReader[ref.Sha256Hash](paramsR)
	if err != nil {
		return nil, err
	}
	cfgEncoded, err := cs.ProvideStockConfig(core.StringFull(ref.AsRawString(cfgHash)))
	if err != nil {
		return nil, err
	}
	var cfgR = msgp.NewReader(bytes.NewReader(cfgEncoded))
	var cfgInstance C
	err = m.configSchema.DecodeMsg(cfgR, &cfgInstance)
	if err != nil {
		return nil, err
	}
	var paramsInstance P
	err = m.paramsSchema.DecodeMsg(paramsR, &paramsInstance)
	if err != nil {
		return nil, err
	}
	applier, err := paramsInstance.BuildTransform(&cfgInstance)
	if err != nil {
		return nil, err
	}
	return newSingleTransform(m, &cfgInstance, &paramsInstance, applier), nil
}
