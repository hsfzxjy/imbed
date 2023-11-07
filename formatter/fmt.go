package formatter

import (
	"io"
	"strings"
	"text/template"

	"github.com/hsfzxjy/imbed/util/iter"
)

type Stringer interface {
	FmtString() string
}

type Humanizer interface {
	FmtHumanize() string
}

type Field[T any] struct {
	Name   string
	Header string
	// whether to show by default
	Show   bool
	Getter func(T) any
}

type Provider[T any] interface {
	ProvideFmtMetadata() []Field[T]
}

type kind int

const (
	kTable kind = 1 + iota
	kJson
	kCustom
)

type Formatter[T any] struct {
	humanized bool
	builder   EncoderBuilder[T]
}

func New[T any](fields []*Field[T], tmpl string, humanized bool) *Formatter[T] {
	f := new(Formatter[T])
	fieldMap := make(map[string]*Field[T], len(fields))
	defaultFields := make([]string, 0, len(fields))
	allFields := make([]string, 0, len(fields))
	for _, field := range fields {
		if field.Show {
			defaultFields = append(defaultFields, "{{."+field.Name+"}}")
		}
		allFields = append(allFields, field.Name)
		fieldMap[field.Name] = field
	}

	var kind kind
	if tmpl == "json" {
		kind = kJson
	} else if tmpl == "table" {
		kind = kTable
		tmpl = ""
	} else {
		trimmed := strings.TrimPrefix(tmpl, "table ")
		if len(trimmed) < len(tmpl) {
			kind = kTable
			tmpl = trimmed
		}
		kind = kCustom
	}

	if kind == kTable && tmpl == "" {
		tmpl = strings.Join(defaultFields, "\t")
	}
	if len(tmpl) == 0 || tmpl[len(tmpl)-1] != '\n' {
		tmpl = tmpl + "\n"
	}
	switch kind {
	case kCustom, kTable:
		f.builder = &tableEncoderBuilder[T]{
			tmpl:        template.Must(template.New("").Parse(tmpl)),
			fieldMap:    fieldMap,
			humanized:   humanized,
			printHeader: kind == kTable,
		}
	case kJson:
		f.builder = &jsonEncoderBuilder[T]{
			allFields: allFields,
			fieldMap:  fieldMap,
			humanize:  humanized,
		}
	}

	return f
}

func (f *Formatter[T]) ExecIter(out io.Writer, it iter.Ator[T]) error {
	encoder := f.builder.build(out)

	for it.HasNext() {
		item := it.Next()
		if item.IsErr() {
			return item.UnwrapErr()
		}
		err := encoder.encodeItem(item.Unwrap())
		if err != nil {
			return err
		}
	}
	return encoder.finalize()
}

func (f *Formatter[T]) Exec(out io.Writer, data []T) error {
	encoder := f.builder.build(out)
	for _, item := range data {
		err := encoder.encodeItem(item)
		if err != nil {
			return err
		}
	}
	return encoder.finalize()
}
