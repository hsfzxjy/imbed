package core

import (
	"github.com/hsfzxjy/imbed/core/ref"
	"github.com/hsfzxjy/imbed/schema"
)

type ConfigProvider interface {
	ProvideWorkspaceConfig(key string) (schema.Reader, error)
	ProvideStockConfig(ref.Sha256Hash) ([]byte, error)
}

type App interface {
	DBDir() string
	WorkDir() string
	ConfigFilePath() string
	FilePath(filename string) string
	Mode() Mode
	BuildMode() BuildMode

	Config() ConfigProvider
}

type Mode int

const (
	ModeNone Mode = iota
	ModeReadonly
	ModeReadWrite
)

type BuildMode int

const (
	BuildUseCache BuildMode = iota
	BuildRedeploy
	BuildAll
)
