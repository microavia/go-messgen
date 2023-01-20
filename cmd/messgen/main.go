package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/powerman/structlog"

	"github.com/microavia/go-messgen/internal/config"
	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/validator"
)

func main() {
	structlog.DefaultLogger.SetLogLevel(structlog.INF)

	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		structlog.DefaultLogger.Fatal(err)
	}

	if *cfg.Verbose {
		structlog.DefaultLogger.SetLogLevel(structlog.DBG)
	}

	structlog.DefaultLogger.Debug("started", "config", cfg)

	baseDirs, err := absDirs(*cfg.BaseDirs)
	if err != nil {
		structlog.DefaultLogger.Fatal("loading definitions: ", err)
	}

	def, err := definition.LoadModules(os.DirFS("/"), baseDirs, *cfg.Modules)
	if err != nil {
		structlog.DefaultLogger.Fatal("loading definitions: ", err)
	}

	err = validator.Validate(def)
	if err != nil {
		structlog.DefaultLogger.Fatal("validating definitions: ", err)
	}

	fmt.Printf("%s\n", prettyPrint(def))
}

func prettyPrint(v interface{}) string {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		panic(err)
	}

	return string(b)
}

var errNotFound = errors.New("not found")

func absDirs(in []string) ([]string, error) {
	out := make([]string, 0, len(in))
	for _, dir := range in {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}

		out = append(out, absDir)
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no one of %+v: %w", in, errNotFound)
	}

	return out, nil
}
