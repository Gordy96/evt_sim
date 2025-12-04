package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"unicode"
)

type TemplateData struct {
	PackageName  string
	StructName   string
	InvokerName  string
	MapParameter string
	Types        []TypeName
	EnableWith   bool
	EnableGet    bool
	EnableSet    bool
	UsePointer   bool
	Builder      Struct
}

var specialTypeNames = map[TypeName]string{
	"[]byte":   "Bytes",
	"[]string": "Strings",
}

type TypeName string

func (t TypeName) Name() string {
	if s, ok := specialTypeNames[t]; ok {
		return s
	}

	var name = string(t)
	if strings.Contains(name, "[]") {
		name = strings.Replace(string(t), "[]", "", 1) + "List"
	}
	name = strings.ToUpper(name[0:1]) + name[1:]

	if strings.Contains(name, ".") {
		parts := strings.Split(name, ".")
		name = parts[len(parts)-1]
	}

	return name
}

var primitives = []TypeName{
	"bool",
	"byte",
	"complex64",
	"complex128",
	"float32",
	"float64",
	"string",
	"uint8",
	"uint16",
	"uint32",
	"uint64",
	"int8",
	"int16",
	"int32",
	"int64",
	"time.Time",
	"time.Duration",
}

func makeTypes(types string) []TypeName {
	var producedTypes []TypeName

	if types != "" {
		for _, t := range strings.Split(types, ",") {
			producedTypes = append(producedTypes, TypeName(t))
		}
	} else {
		producedTypes = make([]TypeName, len(primitives))
		copy(producedTypes, primitives)
		for _, p := range primitives {
			producedTypes = append(producedTypes, "[]"+p)
		}
	}

	return producedTypes
}

type Field struct {
	Name        string
	Type        string
	Passthrough bool
}

type Struct struct {
	Name    string
	Fields  []Field
	Package string
}

func parseType(tn string) Struct {
	fset := token.NewFileSet()

	var s Struct

	s.Package = os.Getenv("GOPACKAGE")
	s.Name = tn

	fpath := filepath.Join(os.Getenv("PWD"), os.Getenv("GOFILE"))

	content, err := os.ReadFile(fpath)

	// Parse file for AST
	f, err := parser.ParseFile(fset, fpath, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	ast.Inspect(f, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// Check if it's a struct declaration
		if st, ok := ts.Type.(*ast.StructType); ok {
			if ts.Name.Name != tn {
				return true
			}
			for i, field := range st.Fields.List {
				var names []string
				var tag reflect.StructTag
				for _, name := range field.Names {
					names = append(names, name.Name)
				}

				if field.Tag != nil {
					tag = reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
				}

				ff := Field{
					Name:        field.Names[0].Name,
					Type:        string(content[field.Type.Pos()-1 : field.Type.End()-1]),
					Passthrough: strings.Contains(tag.Get("builder"), "passthrough"),
				}

				s.Fields = append(s.Fields, ff)
				log.Printf("%d %+v", i, ff)
			}
		}

		return true
	})

	return s
}

func main() {
	chain := flag.Bool("chain", false, "allow chaining setters `With{type}()`")
	get := flag.Bool("get", true, "enable getters `Get{type}()`")
	set := flag.Bool("set", true, "enable setters `Set{type()}`")
	pointer := flag.Bool("pointer", true, "use pointer accessor")
	structName := flag.String("struct", "", "struct name")
	invokerName := flag.String("invoker", "", "caller/invoker name")
	mapParameter := flag.String("property", "params", "property name inside struct")
	types := flag.String("types", "", "comma-separated types")
	templatePath := flag.String("templates", ".", "path to templates")
	builder := flag.Bool("builder", false, "create builder")

	flag.Parse()

	if *structName == "" {
		panic("struct name is required")
	}

	if *invokerName == "" {
		*invokerName = string(strings.ToLower(*structName)[0])
	}

	_ = builder

	s := parseType(*structName)

	data := TemplateData{
		PackageName:  os.Getenv("GOPACKAGE"),
		StructName:   *structName,
		InvokerName:  *invokerName,
		MapParameter: *mapParameter,
		EnableWith:   *chain,
		EnableSet:    *set,
		EnableGet:    *get,
		Types:        makeTypes(*types),
		UsePointer:   *pointer,
		Builder:      s,
	}

	execTemplate(
		filepath.Join(*templatePath, "getter_setter.tmpl"),
		fmt.Sprintf("%s_getter_setter.go", strings.Replace(snakeCase(s.Name), "*", "", -1)),
		data,
	)
	execTemplate(
		filepath.Join(*templatePath, "common.tmpl"),
		fmt.Sprintf("%s_getter_setter_common.go", s.Package),
		data,
	)
}

func execTemplate(tpl string, output string, data TemplateData) {
	tmpl, err := template.ParseFiles(tpl)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	err = tmpl.Execute(file, data)
	if err != nil {
		panic(err)
	}
}

func snakeCase(s string) string {
	var builder strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				builder.WriteRune('_')
			}
			builder.WriteRune(unicode.ToLower(r))
		} else {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}
