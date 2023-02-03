package golang

import (
	"embed"
	"io/fs"
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/stoewer/go-strcase"

	"github.com/microavia/go-messgen/internal/sizer"
	"github.com/microavia/go-messgen/internal/stdtypes"
)

//go:embed templates/*
var tmplSrc embed.FS

var tmplCompiled = func() map[string]*template.Template { //nolint:gochecknoglobals
	out := make(map[string]*template.Template)

	templNames, err := fs.Glob(tmplSrc, "templates/*.tmpl")
	if err != nil {
		panic(err)
	}

	for _, name := range templNames {
		baseName := filepath.Base(name)
		shortName := strings.TrimSuffix(baseName, ".tmpl")

		out[shortName] = template.Must(template.New(baseName).Funcs(TmplFuncs).ParseFS(tmplSrc, name))
	}

	return out
}()

var TmplFuncs = template.FuncMap{ //nolint:gochecknoglobals
	"CamelCase":     camelCase,
	"FixValue":      fixValue,
	"MinSize":       sizer.MinSize,
	"MinSizeByName": sizer.MinSizeByName,
	"ListStrings":   listStrings,
	"RandInt":       rand.Intn,
	"HasSuffix":     strings.HasSuffix,
}

func camelCase(str string) string {
	if _, ok := stdtypes.StdTypes[str]; ok {
		return str
	}

	return strcase.UpperCamelCase(str)
}

var fixValueRE = regexp.MustCompile(`\s*\(?(\d+)U\s*<<\s*(\d+)U\)?`)

func fixValue(value string) string {
	return fixValueRE.ReplaceAllString(value, `$1 << $2`)
}

func listStrings(in ...string) []string {
	return in
}
