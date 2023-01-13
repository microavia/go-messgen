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
	Proto: definition.Proto{ProtoID: 777},
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
	Messages: map[string]definition.Message{
		"message1": {
			ID: 1,
			Fields: []definition.MessageField{
				{Name: "field1", Type: "Bool", Description: "field 1"},
			},
		},
		"message2": {
			ID: 2,
			Fields: []definition.MessageField{
				{Name: "field1", Type: "Bool", Description: "field 1"},
				{Name: "field2", Type: "string", Description: "field 2"},
			},
		},
		"message3": {
			ID: 3,
			Fields: []definition.MessageField{
				{Name: "field1", Type: "Bool", Description: "field 1"},
				{Name: "field2", Type: "string", Description: "field 2"},
				{Name: "field3", Type: "int8", Description: "field 3"},
			},
		},
		"message4": {
			ID: 4,
			Fields: []definition.MessageField{
				{Name: "field1", Type: "Bool", Description: "field 1"},
				{Name: "field2", Type: "string", Description: "field 2"},
				{Name: "field3", Type: "int8", Description: "field 3"},
				{Name: "field4", Type: "float", Description: "field 4"},
			},
		},
		"message5": {
			ID: 5,
			Fields: []definition.MessageField{
				{Name: "field1", Type: "Bool", Description: "field 1"},
				{Name: "field2", Type: "string", Description: "field 2"},
				{Name: "field3", Type: "int8", Description: "field 3"},
				{Name: "field4", Type: "float", Description: "field 4"},
				{Name: "field5", Type: "float", Description: "field 5"},
			},
		},
		"message6": {
			ID: 6,
			Fields: []definition.MessageField{
				{Name: "field1", Type: "Bool", Description: "field 1"},
				{Name: "field2", Type: "string", Description: "field 2"},
				{Name: "field3", Type: "int8", Description: "field 3"},
				{Name: "field4", Type: "float", Description: "field 4"},
				{Name: "field5", Type: "float", Description: "field 5"},
				{Name: "field6", Type: "Language", Description: "field 6"},
			},
		},
	},
	Service: definition.Service{
		Serving: map[string]string{
			"message1": "",
			"message2": "",
			"message3": "Bool",
		},
		Sending: map[string]string{
			"message5": "",
			"message6": "Language",
		},
	},
}

var expected = map[config.Module]*definition.Definition{
	{Vendor: "vendor1", Protocol: "protocol1"}: &expectedModule,
	{Vendor: "vendor1", Protocol: "protocol2"}: &expectedModule,
	{Vendor: "vendor2", Protocol: "protocol1"}: &expectedModule,
}

func TestLoadSingle(t *testing.T) {
	t.Parallel()

	module1 := config.Module{Vendor: "vendor1", Protocol: "protocol1"}

	d, err := definition.LoadModules(
		testdata,
		[]string{"testdata/base1"},
		[]config.Module{module1},
	)
	require.NoError(t, err, "load definition")
	require.Equal(t, expected[module1], d[module1], "load definition")
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
