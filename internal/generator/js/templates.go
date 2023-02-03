package js

import (
	"crypto/md5"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
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
	"CamelCase":       camelCase,
	"VersionProtocol": versionProtocol,
}

func camelCase(str string) string {
	if _, ok := stdtypes.StdTypes[str]; ok {
		return strcase.UpperCamelCase(str)
	}

	return str
}

func versionProtocol(def definition.Definition) string {
	b, err := json.Marshal(def)
	if err != nil {
		panic(fmt.Errorf("marshaling %+v: %w", def.Module, err))
	}

	checksum := md5.Sum(b)

	return hex.EncodeToString(checksum[:])[:6]
}

var fixValueRE = regexp.MustCompile(`\s*\(?(\d+)U\s*<<\s*(\d+)U\)?`)

func fixValue(value string) string {
	return fixValueRE.ReplaceAllString(value, `$1 << $2`)
}

func listStrings(in ...string) []string {
	return in
}
