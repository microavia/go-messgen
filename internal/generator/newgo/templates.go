package golang

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"
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
		out[strings.TrimSuffix(filepath.Base(name), ".tmpl")] = template.Must(template.ParseFS(tmplSrc, name))
	}

	return out
}
