package formatter

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"
	"text/template"
)

type Encoder[T any] interface {
	encodeItem(item T) error
	finalize() error
}

type EncoderBuilder[T any] interface{ build(io.Writer) Encoder[T] }

type FlushWriter interface {
	io.Writer
	Flush() error
}

type flushable struct{ io.Writer }

func (f flushable) Flush() error { return nil }

type tableEncoderBuilder[T any] struct {
	tmpl        *template.Template
	fieldMap    map[string]*Field[T]
	humanized   bool
	printHeader bool
}

func (b *tableEncoderBuilder[T]) build(out io.Writer) Encoder[T] {
	proxy := make(map[string]string, len(b.fieldMap))
	var fw FlushWriter
	if b.printHeader {
		fw = tabwriter.NewWriter(out, 10, 4, 1, ' ', 0)
		for name, field := range b.fieldMap {
			proxy[name] = field.Header
		}
		b.tmpl.Execute(fw, proxy)
	} else {
		fw = flushable{out}
	}
	return &tableEncoder[T]{
		proxy:   proxy,
		tabw:    fw,
		builder: b,
	}
}

type tableEncoder[T any] struct {
	proxy   map[string]string
	tabw    FlushWriter
	builder *tableEncoderBuilder[T]
}

func (e *tableEncoder[T]) encodeItem(item T) error {
	proxy := e.proxy
	humanized := e.builder.humanized
	for name, field := range e.builder.fieldMap {
		proxy[name] = asString(field.Getter(item), humanized)
	}
	return e.builder.tmpl.Execute(e.tabw, proxy)
}

func (e *tableEncoder[T]) finalize() error {
	return e.tabw.Flush()
}

func asString(item any, humanized bool) string {
	if humanized {
		switch v := item.(type) {
		case Humanizer:
			return v.FmtHumanize()
		}
	}
	switch v := item.(type) {
	case Stringer:
		return v.FmtString()
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

type jsonEncoderBuilder[T any] struct {
	allFields []string
	fieldMap  map[string]*Field[T]
	humanize  bool
}

func (b *jsonEncoderBuilder[T]) build(w io.Writer) Encoder[T] {
	return &jsonEncoder[T]{
		builder: b,
		encoder: json.NewEncoder(w),
		proxy:   make(map[string]string, len(b.allFields)),
	}
}

type jsonEncoder[T any] struct {
	builder *jsonEncoderBuilder[T]
	encoder *json.Encoder
	proxy   map[string]string
}

func (e *jsonEncoder[T]) encodeItem(item T) error {
	for _, name := range e.builder.allFields {
		field := e.builder.fieldMap[name]
		e.proxy[name] = asString(field.Getter(item), e.builder.humanize)
	}
	return e.encoder.Encode(e.proxy)
}

func (*jsonEncoder[T]) finalize() error {
	return nil
}
