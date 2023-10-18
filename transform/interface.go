package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

type ParamFor[C any] interface {
	BuildTransform(cfg C) (Applier, error)
}

type Applier interface {
	asset.Applier
	EncodeMsg(r *Registry, w *msgp.Writer) error
}

type Metadata interface {
	Name() string
	ScanParams(paramsR schema.Scanner) (ParamsWithMetadata, error)
	decodeMsg(msgR *msgp.Reader) (ParamsWithMetadata, error)
	ConfigBuilderWorkspace() ConfigBuilder
	ConfigBuilderNeedle(ndl.Needle) ConfigBuilder
	ConfigBuilderHash(ref.Sha256Hash) ConfigBuilder
}

type ParamsWithMetadata interface {
	Metadata() Metadata
	VisitParams(v schema.Visitor) error
	BuildWith(cfgBuilder ConfigBuilder, cp core.ConfigProvider) (Transform, error)
}

type Builder interface {
	ParamsWithMetadata
	ConfigHash() ref.Sha256Hash
	Build(cp core.ConfigProvider) (Transform, error)
}

type Transform interface {
	asset.Transform
	Kind() Kind
}
