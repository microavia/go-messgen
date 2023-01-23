package golang

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/microavia/go-messgen/internal/definition"
)

func GenerateModules(outDir string, definitions []*definition.Definition) error {
	for i, module := range definitions {
		if err := GenerateModule(outDir, *module); err != nil {
			return fmt.Errorf("generating module %d of %d: %w", i+1, len(definitions), err)
		}
	}

	return nil
}

func GenerateModule(outDir string, module definition.Definition) error {
	for fileName, tmpl := range tmplCompiled() {
		if fileName == "message.go" {
			continue
		}

		if err := templateExecute(outDir, fileName, tmpl, module); err != nil {
			return fmt.Errorf("module %+v: generating %q: %w", module.Module, fileName, err)
		}
	}

	return nil
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
