package golang

import (
	"embed"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"

	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/stdtypes"
)

//go:embed templates/*
var tmplSrc embed.FS

var tmplCompiled = func() map[string]*template.Template {
	out := make(map[string]*template.Template)

	templNames, err := fs.Glob(tmplSrc, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	for _, name := range templNames {
		baseName := filepath.Base(name)
		shortName := strings.TrimSuffix(baseName, ".tmpl")

		out[shortName] = template.Must(template.New(baseName).Funcs(TmplFuncs).ParseFS(tmplSrc, name))

		log.Printf("file %q %q tmpl %+v", name, shortName, out[shortName])

	}

	return out
}()

var TmplFuncs = template.FuncMap{
	"CamelCase":     strcase.UpperCamelCase,
	"TrimPrefix":    strings.TrimPrefix,
	"FixValue":      fixValue,
	"In":            isIn,
	"Iterate":       iterate,
	"IsSizeDynamic": isSizeDynamic,
}

var fixValueRE = regexp.MustCompile(`\s*\(?(\d+)U\s*<<\s*(\d+)U\)?`)

func fixValue(value string) string {
	return fixValueRE.ReplaceAllString(value, `$1 << $2`)
}

func isIn(raw any, key string) bool {
	switch list := raw.(type) {
	case []definition.Enum:
		for _, item := range list {
			if item.Name == key {
				return true
			}
		}
	case []definition.Message:
		for _, item := range list {
			if item.Name == key {
				return true
			}
		}
	}

	return false
}

func iterate(n int) []int {
	out := make([]int, n)

	for i := range out {
		out[i] = i
	}

	return out
}

func isSizeDynamic(typeName string, def definition.Definition) bool {
	if stdType, ok := stdtypes.StdTypes[typeName]; ok {
		return stdType.DynamicSize
	}

	for _, enum := range def.Enums {
		if enum.Name == typeName {
			return stdtypes.StdTypes[enum.BaseType].DynamicSize
		}
	}

	for _, m := range def.Messages {
		if m.Name == typeName {
			for _, f := range m.Fields {
				if (f.Type.Array && f.Type.ArraySize == 0) || isSizeDynamic(f.Type.Name, def) {
					return true
				}
			}
		}
	}

	return false
}
