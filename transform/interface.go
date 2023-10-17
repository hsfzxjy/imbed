package transform

import (
	"github.com/hsfzxjy/imbed/asset"
	"github.com/hsfzxjy/imbed/core"
	ndl "github.com/hsfzxjy/imbed/core/needle"
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
	"github.com/tinylib/msgp/msgp"
)

type IParam[C any, A IApplier] interface {
	BuildTransform(cfg C) (A, error)
}

type IApplier interface {
	asset.Applier
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
