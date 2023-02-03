package js

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/sortfields"
)

func GenerateModules(outDir string, definitions []definition.Definition) error {
	for i, module := range definitions {
		if err := GenerateModule(outDir, module); err != nil {
			return fmt.Errorf("generating module %d of %d: %w", i+1, len(definitions), err)
		}
	}

	return nil
}

func GenerateModule(outDir string, module definition.Definition) error {
	return GenerateModuleByTemplates(outDir, module, tmplCompiled)
}

func GenerateModuleByTemplates(
	outDir string,
	module definition.Definition,
	tmpls map[string]*template.Template,
) error {
	sortfields.SortFields(module)

	outDir = filepath.Join(outDir, module.Module.Vendor, module.Module.Protocol)

	for fileName, tmpl := range tmpls {
		err := templateExecute(outDir, fileName, tmpl, module)
		if err != nil {
			return fmt.Errorf("module %+v: generating %q: %w", module.Module, fileName, err)
		}
	}

	return nil
}

func templateExecute(dir, fileName string, tmpl *template.Template, data definition.Definition) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("creating directory %q: %w", dir, err)
	}

	log.Printf("writing %q", filepath.Join(dir, fileName))

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
