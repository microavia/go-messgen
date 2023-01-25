package golang

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/microavia/go-messgen/internal/config"
	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/stdtypes"
)

func GenerateModules(outDir string, definitions []*definition.Definition) error {
	for i, module := range definitions {
		if err := GenerateModule(outDir, *module); err != nil {
			return fmt.Errorf("generating module %d of %d: %w", i+1, len(definitions), err)
		}
	}

	return nil
}

type templateArgs struct {
	Module   definition.Definition
	StdTypes map[string]stdtypes.StdType
}

func GenerateModule(outDir string, module definition.Definition) error {
	outDir = filepath.Join(outDir, module.Module.Vendor, module.Module.Protocol, "message")

	for fileName, tmpl := range tmplCompiled {
		err := templateExecute(outDir, fileName, tmpl, templateArgs{Module: module, StdTypes: stdtypes.StdTypes})
		if err != nil {
			return fmt.Errorf("module %+v: generating %q: %w", module.Module, fileName, err)
		}
	}

	return nil
}

type messageArgs struct {
	Module    config.Module
	Constants []definition.Enum
	Messages  []definition.Message
	Message   []definition.Message
}

func templateExecute(dir, fileName string, tmpl *template.Template, data interface{}) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating directory %q: %w", dir, err)
	}

	f, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return fmt.Errorf("creating file %q/%q: %w", dir, fileName, err)
	}

	defer f.Close()

	if err = tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("executing template %q (%q/%q): %w", tmpl.Name(), dir, fileName, err)
	}

	return nil
}
