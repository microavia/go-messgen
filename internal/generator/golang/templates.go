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
	"CamelCase":  strcase.UpperCamelCase,
	"TrimPrefix": strings.TrimPrefix,
	"FixValue":   fixValue,
	"In":         isIn,
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
