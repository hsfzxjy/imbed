package lang

import (
	"github.com/hsfzxjy/imbed/db"
	"github.com/hsfzxjy/imbed/parser"
	"github.com/hsfzxjy/imbed/transform"
)

type Context struct {
	app      db.App
	parser   *parser.Parser
	registry transform.Registry
}

func NewContext(parser *parser.Parser, app db.App, registry transform.Registry) *Context {
	return &Context{
		app:      app,
		parser:   parser,
		registry: registry,
	}
}
