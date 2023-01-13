package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microavia/go-messgen/internal/config"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	args := []string{
		"-b", "basedir1",
		"-b", "basedir2",
		"-m", "vendor1/protocol1",
		"-m", "vendor1/protocol2",
		"-D", "key1=value1",
		"-D", "key2=value2",
		"-D", "key1=value3",
		"-l", "go",
		"-o", "outputdir",
		"-v",
	}

	expected := config.Config{
		BaseDirs: &[]string{"basedir1", "basedir2"},
		Modules: &[]config.Module{
			{Vendor: "vendor1", Protocol: "protocol1"},
			{Vendor: "vendor1", Protocol: "protocol2"},
		},
		Defines: &map[string]string{
			"key1": "value3",
			"key2": "value2",
		},
		Lang:    Pointer("go"),
		OutDir:  Pointer("outputdir"),
		Verbose: Pointer(true),
	}

	cfg, err := config.Parse(args)
	require.NoError(t, err, "parse config")

	cfg.App = nil

	require.Equal(t, expected, cfg, "parse config")
}

func TestConfigBadModule(t *testing.T) {
	t.Parallel()

	_, err := config.Parse([]string{"-m", "vendor1/", "-b", "basedir1", "-l", "go"})
	require.ErrorContains(t, err, `invalid module format: "vendor1/": invalid argument`, "parse invalid config")
}

func Pointer[T any](v T) *T {
	return &v
}
