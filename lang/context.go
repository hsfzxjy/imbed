package lang

import (
	"github.com/hsfzxjy/imbed/app"
	"github.com/hsfzxjy/imbed/parser"
	"github.com/hsfzxjy/imbed/transform"
)

type Context struct {
	app      *app.App
	parser   *parser.Parser
	registry transform.Registry
}

func NewContext(parser *parser.Parser, app *app.App) *Context {
	return &Context{
		app:      app,
		parser:   parser,
		registry: app.Registry(),
	}
}
