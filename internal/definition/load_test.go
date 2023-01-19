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
	Constants: []definition.Constant{
		{
			Name:     "Bool",
			BaseType: "uint8",
			Fields: []definition.ConstantField{
				{Name: "false", Value: "0"},
				{Name: "true", Value: "1"},
			},
		},
		{
			Name:     "Language",
			BaseType: "string",
			Fields: []definition.ConstantField{
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

func TestLoadNonExisting(t *testing.T) {
	t.Parallel()

	d, err := definition.LoadModules(
		testdata,
		[]string{"testdata/base1"},
		[]config.Module{{Vendor: "vendor3", Protocol: "protocol1"}},
	)
	require.ErrorIs(t, err, definition.ErrNotExist, "load not existing: %v", err)
	require.Nil(t, d, "load not existing")

	d, err = definition.LoadModules(
		testdata,
		[]string{"testdata/base3"},
		[]config.Module{{Vendor: "vendor1", Protocol: "protocol1"}},
	)
	require.ErrorIs(t, err, definition.ErrNotExist, "load not existing: %v", err)
	require.Nil(t, d, "load not existing")

	d, err = definition.LoadModules(
		testdata,
		[]string{"testdata/baseBad"},
		[]config.Module{{Vendor: "vendor1", Protocol: "protocol2"}},
	)
	require.ErrorIs(t, err, definition.ErrNotExist, "load not existing: %v", err)
	require.Nil(t, d, "load not existing")

	d, err = definition.LoadModules(
		testdata,
		[]string{"testdata/baseBad"},
		[]config.Module{{Vendor: "vendor2", Protocol: "protocol1"}},
	)
	require.ErrorIs(t, err, definition.ErrNotExist, "load not existing: %v", err)
	require.Nil(t, d, "load not existing")

	d, err = definition.LoadModules(
		testdata,
		[]string{"testdata/baseBad"},
		[]config.Module{{Vendor: "vendor1", Protocol: "protocol1"}},
	)
	require.ErrorIs(t, err, definition.ErrBadSource, "load not existing: %v", err)
	require.Nil(t, d, "load not existing")
}

func fieldType(name string, isArray bool, arraySize int) definition.FieldType {
	return definition.FieldType{
		Name:      definition.TypeName(name),
		Array:     isArray,
		ArraySize: arraySize,
	}
}
