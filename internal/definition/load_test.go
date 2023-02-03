//nolint:gochecknoglobals
package definition_test

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microavia/go-messgen/internal/config"
	"github.com/microavia/go-messgen/internal/definition"
)

//go:embed testdata/*
//go:embed testdata/*/*
//go:embed testdata/*/*/*/*
var testdata embed.FS

var expectedModule = definition.Definition{
	Proto: definition.Proto{ID: 10},
	Enums: []definition.Enum{
		{
			Name:     "Bool",
			BaseType: "uint8",
			Values: []definition.EnumValue{
				{Name: "false", Value: "0"},
				{Name: "true", Value: "1"},
			},
		},
		{
			Name:     "Language",
			BaseType: "string",
			Values: []definition.EnumValue{
				{Name: "go", Value: "go", Description: "go lang"},
				{Name: "js", Value: "js", Description: "js lang"},
				{Name: "cpp", Value: "cpp", Description: "cpp lang"},
				{Name: "md", Value: "md", Description: "md lang"},
			},
		},
	},
	Messages: []definition.Message{
		{
			Name: "message1",
			ID:   1,
			Fields: []definition.MessageField{
				{Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
			},
		},
		{
			Name: "message2",
			ID:   2,
			Fields: []definition.MessageField{
				{Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
				{Name: "field2", Type: fieldType("string", false, 0), Description: "field 2"},
			},
		},
		{
			Name: "message3",
			ID:   3,
			Fields: []definition.MessageField{
				{Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
				{Name: "field2", Type: fieldType("string", false, 0), Description: "field 2"},
				{Name: "field3", Type: fieldType("int8", false, 0), Description: "field 3"},
			},
		},
		{
			Name: "message4",
			ID:   4,
			Fields: []definition.MessageField{
				{Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
				{Name: "field2", Type: fieldType("string", false, 0), Description: "field 2"},
				{Name: "field3", Type: fieldType("int8", false, 0), Description: "field 3"},
				{Name: "field4", Type: fieldType("float", false, 0), Description: "field 4"},
			},
		},
		{
			Name: "message5",
			ID:   5,
			Fields: []definition.MessageField{
				{Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
				{Name: "field2", Type: fieldType("string", false, 0), Description: "field 2"},
				{Name: "field3", Type: fieldType("int8", false, 0), Description: "field 3"},
				{Name: "field4", Type: fieldType("float", true, 10), Description: "field 4"},
				{Name: "field5", Type: fieldType("float", true, 0), Description: "field 5"},
			},
		},
		{
			Name: "message6",
			ID:   6,
			Fields: []definition.MessageField{
				{Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
				{Name: "field2", Type: fieldType("string", false, 0), Description: "field 2"},
				{Name: "field3", Type: fieldType("int8", false, 0), Description: "field 3"},
				{Name: "field4", Type: fieldType("float", true, 10), Description: "field 4"},
				{Name: "field5", Type: fieldType("float", true, 0), Description: "field 5"},
				{Name: "field6", Type: fieldType("Language", false, 0), Description: "field 6"},
			},
		},
	},
	Service: definition.Service{
		Serving: []definition.ServicePair{
			{Request: "message1"},
			{Request: "message2"},
			{Request: "message3", Response: "Bool"},
		},
		Sending: []definition.ServicePair{
			{Request: "message4"},
			{Request: "message5"},
			{Request: "message6", Response: "Bool"},
		},
	},
}

var expected = []definition.Definition{
	setModuleID(expectedModule, config.Module{Vendor: "vendor1", Protocol: "protocol1"}),
	setModuleID(expectedModule, config.Module{Vendor: "vendor1", Protocol: "protocol2"}),
	setModuleID(expectedModule, config.Module{Vendor: "vendor2", Protocol: "protocol1"}),
}

func setModuleID(def definition.Definition, module config.Module) definition.Definition {
	def.Module = module

	return def
}

func TestLoadSingle(t *testing.T) {
	t.Parallel()

	d, err := definition.LoadModules(
		testdata,
		[]string{"testdata/base1"},
		[]config.Module{{Vendor: "vendor1", Protocol: "protocol1"}},
	)
	require.NoError(t, err, "load definition")
	require.Len(t, d, 1, "load definition")
	require.Equal(t, expected[0], d[0], "load definition")
}

func TestLoadMultiple(t *testing.T) {
	t.Parallel()

	d, err := definition.LoadModules(
		testdata,
		[]string{"testdata/base1", "testdata/base2"},
		[]config.Module{
			{Vendor: "vendor1", Protocol: "protocol1"},
			{Vendor: "vendor1", Protocol: "protocol2"},
			{Vendor: "vendor2", Protocol: "protocol1"},
		},
	)
	require.NoError(t, err, "load definition")
	require.Equal(t, expected, d, "load definition")
}

func TestLoadBad(t *testing.T) {
	t.Parallel()

	type testRow struct {
		name     string
		baseDirs []string
		modules  []config.Module
		err      error
	}

	testRows := []testRow{
		{
			name:     "load not existing vendor",
			baseDirs: []string{"testdata/base1"},
			modules:  []config.Module{{Vendor: "vendor3", Protocol: "protocol1"}},
			err:      definition.ErrNotExist,
		},
		{
			name:     "load not existing base dir",
			baseDirs: []string{"testdata/base3"},
			modules:  []config.Module{{Vendor: "vendor1", Protocol: "protocol1"}},
			err:      definition.ErrNotExist,
		},
		{
			name:     "load not a dir protocol",
			baseDirs: []string{"testdata/baseBad"},
			modules:  []config.Module{{Vendor: "vendor1", Protocol: "protocol2"}},
			err:      definition.ErrNotExist,
		},
		{
			name:     "load not a dir vendor",
			baseDirs: []string{"testdata/baseBad"},
			modules:  []config.Module{{Vendor: "vendor2", Protocol: "protocol1"}},
			err:      definition.ErrNotExist,
		},
		{
			name:     "load bad definition",
			baseDirs: []string{"testdata/baseBad"},
			modules:  []config.Module{{Vendor: "vendor1", Protocol: "protocol1"}},
			err:      definition.ErrNotExist,
		},
	}

	for i, row := range testRows {
		d, err := definition.LoadModules(testdata, row.baseDirs, row.modules)
		require.ErrorIs(t, err, row.err, "test %d of %d: %s: %v", i+1, len(testRows), row.name, err)
		require.Nil(t, d, "test %d of %d: %s", i+1, len(testRows), row.name)
	}
}

func fieldType(name string, isArray bool, arraySize int) definition.FieldType {
	return definition.FieldType{
		Name:      name,
		Array:     isArray,
		ArraySize: arraySize,
	}
}
