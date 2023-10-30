package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var schemagenDirective = regexp.MustCompile("(?m:^//imbed\\:schemagen(?:\\s+(\\S+))?)")

func matchSchemagen(cg *ast.CommentGroup) (string, bool) {
	if cg == nil {
		return "", false
	}
	for _, c := range cg.List {
		matched := schemagenDirective.FindStringSubmatch(c.Text)
		if len(matched) > 0 {
			return matched[1], true
		}
	}
	return "", false
}

func main() {
	fset := token.NewFileSet()
	filename := os.Getenv("GOFILE")
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	check(err)
	var printer Printer
	var imports = map[string]struct{}{}
	var decls []*ast.GenDecl
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		name, ok := matchSchemagen(gd.Doc)
		if !ok {
			continue
		}
		decls = append(decls, gd)
		handleDecl(&printer, gd, name, imports)
	}
	content := printer.String()
	printer.Reset()
	printer.
		Printfln("// Code generated by \"github.com/hsfzxjy/imbed/schema/gen\"; DO NOT EDIT.\n").
		Printfln("package %s\n", os.Getenv("GOPACKAGE"))

	if len(imports) > 0 {
		printer.Printfln("import (")
		for name := range imports {
			printer.Printfln(strconv.Quote(name))
		}
		printer.Printfln(")")
	}

	output, err := format.Source([]byte(printer.String() + content))
	check(err)

	dir, base := path.Split(filename)
	check(os.WriteFile(path.Join(dir, strings.Replace(base, ".go", "_schemagen.go", 1)), output, 0o644))
}

type Printer struct{ strings.Builder }

func (p *Printer) Printf(tmpl string, args ...any) *Printer {
	fmt.Fprintf(&p.Builder, tmpl, args...)
	return p
}
func (p *Printer) Printfln(tmpl string, args ...any) *Printer {
	fmt.Fprintf(&p.Builder, tmpl+"\n", args...)
	return p
}

type FieldConfig struct {
	TypeHint string
	Rename   string
	Default  string
}

func getConfig(tag *ast.BasicLit) *FieldConfig {
	if tag == nil {
		return nil
	}
	value, err := strconv.Unquote(tag.Value)
	check(err)
	item, ok := reflect.StructTag(value).Lookup("imbed")
	if !ok {
		return nil
	}
	configs := strings.SplitN(item, ",", 2)
	var cfg *FieldConfig
	switch len(configs) {
	case 1:
		cfg = &FieldConfig{Rename: configs[0]}
	case 2:
		cfg = &FieldConfig{Rename: configs[0], Default: configs[1]}
	}
	if cfg != nil {
		parts := strings.SplitN(cfg.Rename, "!", 2)
		if len(parts) == 2 {
			cfg.Rename = parts[0]
			cfg.TypeHint = parts[1]
		}
	}
	return cfg
}

func handleFieldType(ftyp ast.Expr, typeHint string) (typname, cntr string, success bool) {
	var ok bool
	switch ftyp := ftyp.(type) {
	case *ast.Ident:
		typname = ftyp.String()
		if cntr, ok = atomTypeMap[typname]; ok {
			cntr = fmt.Sprintf("schema.%s()", cntr)
		} else if cntr, ok = atomTypeMap[typeHint]; typeHint != "" && ok {
			cntr = fmt.Sprintf("schema.%s()", cntr)
		} else {
			cntr = typname + "Schema"
		}
	case *ast.StarExpr:
		var sel *ast.SelectorExpr
		if sel, ok = ftyp.X.(*ast.SelectorExpr); !ok {
			goto DESCEND
		}
		if ident, ok := sel.X.(*ast.Ident); !ok || ident.String() != "big" {
			goto DESCEND
		}
		if sel.Sel.Name != "Rat" {
			goto DESCEND
		}
		cntr = "schema.Rat()"
	DESCEND:
		tname, c, ok := handleFieldType(ftyp.X, "")
		if !ok {
			goto ERROR
		}
		typname = "*" + tname
		cntr = "schema.Ptr(" + c + ")"
	case *ast.MapType:
		if key, ok := ftyp.Key.(*ast.Ident); !ok || key.String() != "string" {
			goto ERROR
		}
		if _, vcntr, ok := handleFieldType(ftyp.Value, ""); !ok {
			goto ERROR
		} else {
			cntr = "schema.Map(" + vcntr + ")"
		}
	default:
		goto ERROR
	}
	success = true
	return
ERROR:
	success = false
	return
}

func handleDecl(printer *Printer, gd *ast.GenDecl, debugName string, imports map[string]struct{}) {
	var spec *ast.TypeSpec
	var ok bool
	if specs := gd.Specs; len(specs) > 0 {
		if spec, ok = specs[0].(*ast.TypeSpec); !ok {
			return
		}
	} else {
		return
	}
	name := spec.Name.String()
	if name == "<nil>" {
		return
	}
	if debugName == "" {
		debugName = name
	}
	var typ *ast.StructType
	if typ, ok = spec.Type.(*ast.StructType); !ok {
		return
	}
	imports["github.com/hsfzxjy/imbed/schema"] = struct{}{}
	imports["github.com/tinylib/msgp/msgp"] = struct{}{}
	printer.
		Printfln("var %[1]sSchema = schema.StructFunc(func(prototype *%[1]s) *schema.StructBuilder[%[1]s] {", name).
		Printfln("return schema.Struct(prototype,")
	for _, field := range typ.Fields.List {
		cfg := getConfig(field.Tag)

		throw := func() {
			println("cannot handle field")
			spew.Dump(field)
			spew.Dump(cfg)
			panic("")
		}

		if cfg == nil {
			continue
		}
		var fieldName, cntr, typname, def string
		switch len(field.Names) {
		case 1:
			fieldName = field.Names[0].Name
		case 0:
		default:
			throw()
		}

		typname, cntr, ok = handleFieldType(field.Type, cfg.TypeHint)
		if !ok {
			throw()
		}

		if fieldName == "" {
			fieldName = typname
		}

		if cfg.Rename == "" {
			cfg.Rename = fieldName
		}

		if cfg.Default != "" {
			parts := strings.Split(cfg.Default, "!")
			for _, p := range parts[:len(parts)-1] {
				p = strings.Replace(p, "@", "github.com/hsfzxjy/imbed", -1)
				imports[p] = struct{}{}
			}
			def = ".Default(" + parts[len(parts)-1] + ")"
		}

		ptr := "&prototype." + fieldName
		if cfg.TypeHint != "" {
			ptr = fmt.Sprintf("(*%s)(%s)", cfg.TypeHint, ptr)
		}

		printer.Printfln("schema.F(%q, %s, %s%s),", cfg.Rename, ptr, cntr, def)
		continue
	}
	printer.
		Printfln(").DebugName(%q)", debugName).
		Printfln("})\n").
		Printfln("func (x*%s)EncodeMsg(w *msgp.Writer)error{", name).
		Printfln("return %sSchema.Build().EncodeMsg(w, x)", name).
		Printfln("}\n")
}

var atomTypeMap = map[string]string{
	"int64":  "Int",
	"string": "String",
	"bool":   "Bool",
}
