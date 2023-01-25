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
	Proto: definition.Proto{ProtoID: 10},
	Enums: map[string]definition.Enum{
		"Bool": {
			Name:     "Bool",
			BaseType: "uint8",
			Values: map[string]definition.EnumValue{
				"false": {Name: "false", Value: "0"},
				"true":  {Name: "true", Value: "1"},
			},
		},
		"Language": {
			Name:     "Language",
			BaseType: "string",
			Values: map[string]definition.EnumValue{
				"go":  {Name: "go", Value: "go", Description: "go lang"},
				"js":  {Name: "js", Value: "js", Description: "js lang"},
				"cpp": {Name: "cpp", Value: "cpp", Description: "cpp lang"},
				"md":  {Name: "md", Value: "md", Description: "md lang"},
			},
		},
	},
	Messages: map[string]definition.Message{
		"message1": {
			Name: "message1",
			ID:   1,
			Fields: map[string]definition.MessageField{
				"field1": {Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
			},
		},
		"message2": {
			Name: "message2",
			ID:   2,
			Fields: map[string]definition.MessageField{
				"field1": {Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
				"field2": {Name: "field2", Type: fieldType("string", false, 0), Description: "field 2"},
			},
		},
		"message3": {
			Name: "message3",
			ID:   3,
			Fields: map[string]definition.MessageField{
				"field1": {Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
				"field2": {Name: "field2", Type: fieldType("string", false, 0), Description: "field 2"},
				"field3": {Name: "field3", Type: fieldType("int8", false, 0), Description: "field 3"},
			},
		},
		"message4": {
			Name: "message4",
			ID:   4,
			Fields: map[string]definition.MessageField{
				"field1": {Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
				"field2": {Name: "field2", Type: fieldType("string", false, 0), Description: "field 2"},
				"field3": {Name: "field3", Type: fieldType("int8", false, 0), Description: "field 3"},
				"field4": {Name: "field4", Type: fieldType("float", false, 0), Description: "field 4"},
			},
		},
		"message5": {
			Name: "message5",
			ID:   5,
			Fields: map[string]definition.MessageField{
				"field1": {Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
				"field2": {Name: "field2", Type: fieldType("string", false, 0), Description: "field 2"},
				"field3": {Name: "field3", Type: fieldType("int8", false, 0), Description: "field 3"},
				"field4": {Name: "field4", Type: fieldType("float", true, 10), Description: "field 4"},
				"field5": {Name: "field5", Type: fieldType("float", true, 0), Description: "field 5"},
			},
		},
		"message6": {
			Name: "message6",
			ID:   6,
			Fields: map[string]definition.MessageField{
				"field1": {Name: "field1", Type: fieldType("Bool", false, 0), Description: "field 1"},
				"field2": {Name: "field2", Type: fieldType("string", false, 0), Description: "field 2"},
				"field3": {Name: "field3", Type: fieldType("int8", false, 0), Description: "field 3"},
				"field4": {Name: "field4", Type: fieldType("float", true, 10), Description: "field 4"},
				"field5": {Name: "field5", Type: fieldType("float", true, 0), Description: "field 5"},
				"field6": {Name: "field6", Type: fieldType("Language", false, 0), Description: "field 6"},
			},
		},
	},
	Service: definition.Service{
		Serving: map[string]definition.ServicePair{
			"message1": {Request: "message1"},
			"message2": {Request: "message2"},
			"message3": {Request: "message3", Response: "Bool"},
		},
		Sending: map[string]definition.ServicePair{
			"message4": {Request: "message4"},
			"message5": {Request: "message5"},
			"message6": {Request: "message6", Response: "Bool"},
		},
	},
}

var expected = []*definition.Definition{
	setModuleID(&expectedModule, config.Module{Vendor: "vendor1", Protocol: "protocol1"}),
	setModuleID(&expectedModule, config.Module{Vendor: "vendor1", Protocol: "protocol2"}),
	setModuleID(&expectedModule, config.Module{Vendor: "vendor2", Protocol: "protocol1"}),
}

func setModuleID(def *definition.Definition, module config.Module) *definition.Definition {
	out := *def
	out.Module = module
	return &out
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
