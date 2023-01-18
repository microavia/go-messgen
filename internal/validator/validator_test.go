package validator_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microavia/go-messgen/internal/config"
	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/validator"
)

//go:embed testdata/*/*/*/*
var testdata embed.FS

type testRow struct {
	name     string
	basedirs []string
	modules  []config.Module
	err      error
}

var testRows = []testRow{
	{
		name:     "valid",
		basedirs: []string{"testdata/valid"},
		modules:  []config.Module{{Vendor: "vendor1", Protocol: "protocol1"}},
	},
	{
		name:     "no proto id",
		basedirs: []string{"testdata/noprotoid"},
		modules:  []config.Module{{Vendor: "vendor1", Protocol: "protocol1"}},
		err:      validator.ErrNoProtoID,
	},
	{
		name:     "no message id",
		basedirs: []string{"testdata/nomsgid"},
		modules:  []config.Module{{Vendor: "vendor1", Protocol: "protocol1"}},
		err:      validator.ErrNoMsgID,
	},
	{
		name:     "duplicated proto id",
		basedirs: []string{"testdata/dupprotoid"},
		modules: []config.Module{
			{Vendor: "vendor1", Protocol: "protocol1"},
			{Vendor: "vendor1", Protocol: "protocol2"},
		},
		err: validator.ErrDupID,
	},
	{
		name:     "no messages",
		basedirs: []string{"testdata/nomessages"},
		modules:  []config.Module{{Vendor: "vendor1", Protocol: "protocol1"}},
		err:      validator.ErrNoMessages,
	},
	{
		name:     "duplicated message ID",
		basedirs: []string{"testdata/redefined"},
		modules:  []config.Module{{Vendor: "messages", Protocol: "dupmsgid"}},
		err:      validator.ErrDupID,
	},
	{
		name:     "standard type redefined by message",
		basedirs: []string{"testdata/redefined"},
		modules:  []config.Module{{Vendor: "messages", Protocol: "stdtype"}},
		err:      validator.ErrRedefined,
	},
	{
		name:     "redefined constant",
		basedirs: []string{"testdata/redefined"},
		modules:  []config.Module{{Vendor: "messages", Protocol: "constant"}},
		err:      validator.ErrRedefined,
	},
	{
		name:     "duplicate constant name",
		basedirs: []string{"testdata/redefined"},
		modules:  []config.Module{{Vendor: "constants", Protocol: "constname"}},
		err:      validator.ErrDupID,
	},
	{
		name:     "standard type redefined by constant",
		basedirs: []string{"testdata/redefined"},
		modules:  []config.Module{{Vendor: "constants", Protocol: "stdtype"}},
		err:      validator.ErrRedefined,
	},
	{
		name:     "invalid constant base type",
		basedirs: []string{"testdata/badbasetype"},
		modules:  []config.Module{{Vendor: "constants", Protocol: "protocol1"}},
		err:      validator.ErrUnknownType,
	},
	{
		name:     "redefined constant field name",
		basedirs: []string{"testdata/redefined"},
		modules:  []config.Module{{Vendor: "constants", Protocol: "fields"}},
		err:      validator.ErrDupID,
	},
	{
		name:     "redefined message field name",
		basedirs: []string{"testdata/redefined"},
		modules:  []config.Module{{Vendor: "messages", Protocol: "fields"}},
		err:      validator.ErrDupID,
	},
	{
		name:     "invalid message field type",
		basedirs: []string{"testdata/badbasetype"},
		modules:  []config.Module{{Vendor: "messages", Protocol: "protocol1"}},
		err:      validator.ErrUnknownType,
	},
}

func TestValidate(t *testing.T) {
	for _, row := range testRows {
		t.Run(row.name, func(innerT *testing.T) { testRun(innerT, row) })
	}
}

func testRun(t *testing.T, row testRow) {
	t.Parallel()

	if err := validateDef(row.basedirs, row.modules); row.err != nil {
		require.ErrorIs(t, err, row.err, "validate %q: %+v", row.name, err)
	} else {
		require.NoError(t, err, "validate %q: %+v", row.name, err)
	}
}

func validateDef(basedirs []string, modules []config.Module) error {
	def, err := definition.LoadModules(testdata, basedirs, modules)
	if err != nil {
		return err
	}

	return validator.Validate(def)
}
