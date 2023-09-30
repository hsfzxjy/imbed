package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"text/template"
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
	kind      kind
	tmpl      *template.Template
	allFields []string
	fieldMap  map[string]*Field[T]
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
	f.allFields = allFields
	f.fieldMap = fieldMap
	f.humanized = humanized

	if tmpl == "json" {
		f.kind = kJson
	} else if tmpl == "table" {
		f.kind = kTable
		tmpl = ""
	} else {
		trimmed := strings.TrimPrefix(tmpl, "table ")
		if len(trimmed) < len(tmpl) {
			f.kind = kTable
			tmpl = trimmed
		}
		f.kind = kCustom
	}

	if f.kind == kTable && tmpl == "" {
		tmpl = strings.Join(defaultFields, "\t")
	}
	if len(tmpl) == 0 || tmpl[len(tmpl)-1] != '\n' {
		tmpl = tmpl + "\n"
	}
	f.tmpl = template.Must(template.New("").Parse(tmpl))

	return f
}

func (f *Formatter[T]) Exec(out io.Writer, data []T) error {
	if f.kind == kJson {
		return f.execJson(out, data)
	}
	proxy := make(map[string]string, len(f.allFields))
	var tabw *tabwriter.Writer
	if f.kind == kTable {
		tabw = tabwriter.NewWriter(out, 10, 4, 1, ' ', 0)
		out = tabw
		for name, field := range f.fieldMap {
			proxy[name] = field.Header
		}
		err := f.tmpl.Execute(out, proxy)
		if err != nil {
			return err
		}
	}
	for _, item := range data {
		for name, field := range f.fieldMap {
			proxy[name] = f.asString(field.Getter(item))
		}
		err := f.tmpl.Execute(out, proxy)
		if err != nil {
			return err
		}
	}
	if tabw != nil {
		err := tabw.Flush()
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter[T]) execJson(out io.Writer, data []T) error {
	dict := make(map[string]string, len(f.allFields))
	encoder := json.NewEncoder(out)
	for _, item := range data {
		for _, name := range f.allFields {
			field := f.fieldMap[name]
			dict[name] = f.asString(field.Getter(item))
		}
		err := encoder.Encode(dict)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter[T]) asString(value any) string {
	if f.humanized {
		switch v := value.(type) {
		case Humanizer:
			return v.FmtHumanize()
		}
	}
	switch v := value.(type) {
	case Stringer:
		return v.FmtString()
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
