package golang

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/microavia/go-messgen/internal/config"
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
	outDir = filepath.Join(outDir, module.Module.Vendor, module.Module.Protocol, "message")

	for fileName, tmpl := range tmplCompiled {
		if fileName == "messages.go" {
			continue
		}

		log.Printf("file %q tmpl %+v", fileName, tmpl)

		if err := templateExecute(outDir, fileName, tmpl, module); err != nil {
			return fmt.Errorf("module %+v: generating %q: %w", module.Module, fileName, err)
		}
	}

	for _, message := range module.Messages {
		var (
			fileName = message.Name + ".go"
			tmpl     = tmplCompiled["messages.go"]
			msg      = messageArgs{
				Module:    module.Module,
				Constants: module.Enums,
				Messages:  module.Messages,
				Message:   []definition.Message{message},
			}
		)

		if err := templateExecute(outDir, fileName, tmpl, msg); err != nil {
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
