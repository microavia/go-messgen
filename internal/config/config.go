package config

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct { //nolint:musttag
	App    *kingpin.Application
	Parsed string

	BaseDirs *[]string
	Modules  *[]Module
	OutDir   *string
	Lang     *string
	Defines  *map[string]string
	Verbose  *bool
}

func (c Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	return string(b)
}

//nolint:lll
func initKingpin() Config {
	config := Config{}

	config.App = kingpin.New("messgen", "Lightweight and fast message serialization library code generator")

	config.BaseDirs = config.App.Flag("basedir", "Message definition base directories").Short('b').Required().Strings()
	config.Modules = ModulesFlag(config.App.Flag("module", "Modules").Short('m').Required().PlaceHolder("VENDOR/PROTOCOL"))
	config.OutDir = config.App.Flag("outdir", "Output directory").Short('o').Default(".").String()
	config.Lang = config.App.Flag("lang", "Output language (cpp=C++, go=Golang, js=JavaScript, md=Markdown)").Short('l').Required().Enum("cpp", "go", "js", "md")
	config.Defines = config.App.Flag("define", "Defines variables in 'key=value' format").Short('D').PlaceHolder("key=value").StringMap()
	config.Verbose = config.App.Flag("verbose", "Verbose output").Short('v').Bool()

	return config
}

// Parse exported func should have comment or be unexported.
func Parse(args []string) (Config, error) {
	config := initKingpin()

	_, err := config.App.Parse(args)
	if err != nil {
		return config, fmt.Errorf("parsing config: %w", err)
	}

	return config, nil
}

func ModulesFlag(s kingpin.Settings) *[]Module {
	m := make([]Module, 0, 1)

	s.SetValue((*modulesList)(&m))

	return &m
}

var _ kingpin.Value = (*modulesList)(nil)

type modulesList []Module

func (m *modulesList) IsCumulative() bool { return true }

func (m *modulesList) String() string {
	modules := make([]string, 0, len(*m))

	for _, module := range *m {
		modules = append(modules, module.String())
	}

	return strings.Join(modules, ",")
}

var ErrInvalidArgument = fmt.Errorf("invalid argument")

func (m *modulesList) Set(s string) error {
	var module Module

	if err := module.Set(s); err != nil {
		return err
	}

	*m = append(*m, module)

	return nil
}

type Module struct {
	Vendor   string
	Protocol string
}

func (m *Module) Set(s string) error {
	fields := strings.Split(s, string(filepath.Separator))

	if len(fields) != 2 || fields[0] == "" || fields[1] == "" {
		return fmt.Errorf("invalid module format: %q: %w", s, ErrInvalidArgument)
	}

	*m = Module{Vendor: fields[0], Protocol: fields[1]}

	return nil
}

func (m *Module) String() string {
	return filepath.Join(m.Vendor, m.Protocol)
}
