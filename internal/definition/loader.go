package definition

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/microavia/go-messgen/internal/config"
)

var (
	ErrBadSource = errors.New("bad source")
	ErrBadData   = fmt.Errorf("bad data: %w", ErrBadSource)
	ErrNotExist  = fmt.Errorf("not exist: %w", ErrBadSource)
	ErrNotDir    = fmt.Errorf("not a directory: %w", ErrBadSource)
)

func LoadModules(fsys fs.FS, roots []string, modules []config.Module) (map[config.Module]*Definition, error) {
	out := make(map[config.Module]*Definition, len(modules))

	for i, module := range modules {
		def, err := loadModule(fsys, roots, module)
		if err != nil {
			return nil, fmt.Errorf("loading module %d of %d: %+v: %w", i+1, len(modules), module, err)
		}

		out[module] = def
	}

	return out, nil
}

func loadModule(fsys fs.FS, baseDirs []string, module config.Module) (*Definition, error) {
	for _, baseDir := range baseDirs {
		def, err := Load(fsys, baseDir, module)
		if err == nil {
			return def, nil
		}

		if errors.Is(err, ErrBadSource) {
			continue
		}

		return nil, fmt.Errorf("loading %+v: %q: %w", module, err, ErrBadData)
	}

	return nil, fmt.Errorf("loading %+v: %w", module, ErrNotExist)
}

func Load(fsys fs.FS, baseDir string, module config.Module) (*Definition, error) {
	out := &Definition{}

	root := filepath.Join(baseDir, module.Vendor, module.Protocol)

	if err := checkDir(fsys, root); err != nil {
		return nil, fmt.Errorf("loading %+v from %q: %q: %w", module, baseDir, err, ErrBadData)
	}

	err := fs.WalkDir(
		fsys,
		root,
		func(path string, d fs.DirEntry, errInner error) error {
			return checkFile(fsys, root, out, path, d, errInner)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("loading %+v from %q: %q: %w", module, baseDir, err, err)
	}

	return out, nil
}

func checkFile(fsys fs.FS, root string, out *Definition, path string, d fs.DirEntry, err error) error {
	const yamlSuffix = ".yaml"

	switch {
	case err != nil:
		return err
	case d.IsDir():
		return nil
	case path != filepath.Join(root, d.Name()):
		return nil
	case d.Name() == "_protocol.yaml":
		err = loadFile(fsys, path, &out.Proto)
	case d.Name() == "_constants.yaml":
		err = loadFile(fsys, path, &out.Constants)
	case d.Name() == "_service.yaml":
		err = loadFile(fsys, path, &out.Service)
	case filepath.Ext(d.Name()) == yamlSuffix:
		v := Message{}
		if err = loadFile(fsys, path, &v); err == nil {
			out.Messages = appendMap(out.Messages, strings.TrimSuffix(d.Name(), yamlSuffix), v)
		}
	}

	return err
}

func loadFile(fsys fs.FS, path string, o interface{}) error {
	b, err := fs.ReadFile(fsys, path)
	if err != nil {
		return err //nolint:wrapcheck
	}

	return yaml.Unmarshal(b, o) //nolint:wrapcheck
}

func appendMap[T any](m map[string]T, k string, v T) map[string]T {
	if m == nil {
		m = make(map[string]T)
	}

	m[k] = v

	return m
}

func checkDir(fsys fs.FS, name string) error {
	stat, err := fs.Stat(fsys, name)
	if err != nil {
		return fmt.Errorf("%q is not a directory: %w", name, err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("%q: %w", name, ErrNotDir)
	}

	return nil
}
