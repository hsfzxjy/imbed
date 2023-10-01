package cmds

import (
	"github.com/spf13/pflag"
)

type fmtOption struct {
	Raw    bool
	Format string
}

func getFmtFlagSet() (*pflag.FlagSet, *fmtOption) {
	opt := new(fmtOption)
	fs := pflag.NewFlagSet("fmt", pflag.ContinueOnError)
	fs.BoolVar(&opt.Raw, "raw", false, "Don't use humanized output")
	fs.StringVar(&opt.Format, "format", "table", "Format output using a template")
	return fs, opt
}
