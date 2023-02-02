//nolint:gochecknoglobals
package sizer_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/sizer"
)

var def = definition.Definition{
	Enums: []definition.Enum{{Name: "Int16Enum", BaseType: "int16"}},
	Messages: []definition.Message{
		{
			Name: "Message1",
			Fields: []definition.MessageField{
				{Name: "field1", Type: definition.FieldType{Name: "Int16Enum", Array: true, ArraySize: 5}},
			},
		},
		{
			Name: "Message2",
			Fields: []definition.MessageField{
				{Name: "field1", Type: definition.FieldType{Name: "Message1", Array: true, ArraySize: 5}},
			},
		},
	},
}

type testMinSizeRow struct {
	name    string
	typeDef definition.FieldType
	def     definition.Definition
	size    sizer.TypeSize
}

var testMinSizeRows = []testMinSizeRow{
	{
		name:    "static standard type array",
		typeDef: definition.FieldType{Name: "double", Array: true, ArraySize: 5},
		size:    sizer.TypeSize{MinSize: 40, Align: 8},
	},
	{
		name:    "dynamic standard type array",
		typeDef: definition.FieldType{Name: "double", Array: true},
		size:    sizer.TypeSize{MinSize: 12, Align: 8, Dynamic: true},
	},
	{
		name:    "static enum array",
		typeDef: definition.FieldType{Name: "Int16Enum", Array: true, ArraySize: 5},
		size:    sizer.TypeSize{MinSize: 10, Align: 2},
		def:     def,
	},
	{
		name:    "dynamic enum array",
		typeDef: definition.FieldType{Name: "Int16Enum", Array: true},
		size:    sizer.TypeSize{MinSize: 6, Align: 2, Dynamic: true},
		def:     def,
	},
	{
		name:    "static message array",
		typeDef: definition.FieldType{Name: "Message2", Array: true, ArraySize: 5},
		size:    sizer.TypeSize{MinSize: 250, Align: 2},
		def:     def,
	},
	{
		name:    "static message array by name",
		typeDef: definition.FieldType{Name: "Message2", Array: true, ArraySize: 5},
		size:    sizer.TypeSize{MinSize: 250, Align: 2},
		def:     def,
	},
}

func TestMinSize(t *testing.T) {
	t.Parallel()

	for loopI, loopRow := range testMinSizeRows {
		var (
			i   = loopI
			row = loopRow
		)

		t.Run(
			row.name,
			func(t *testing.T) {
				t.Parallel()

				size := sizer.MinSize(row.typeDef, row.def)
				require.Equal(t, row.size, size, "%d of %d: %q: %+v", i+i, len(testMinSizeRows), row.name, size)
			},
		)
	}
}

type testMinSizeByNameRow struct {
	name     string
	typeName string
	def      definition.Definition
	size     sizer.TypeSize
	err      error
}

var testMinSizeByNameRows = []testMinSizeByNameRow{
	{
		name:     "static message array by name",
		typeName: "Message2[7]",
		size:     sizer.TypeSize{MinSize: 350, Align: 2},
		def:      def,
	},
	{
		name:     "dynamic message array by name",
		typeName: "Message2[]",
		size:     sizer.TypeSize{MinSize: 54, Align: 2, Dynamic: true},
		def:      def,
	},
	{
		name:     "invalid type by name",
		typeName: "Message2[a]",
		size:     sizer.TypeSize{MinSize: 54},
		def:      def,
		err:      definition.ErrInvalidInput,
	},
}

func TestMinSizeByName(t *testing.T) {
	t.Parallel()

	for loopI, loopRow := range testMinSizeByNameRows {
		var (
			i   = loopI
			row = loopRow
		)

		t.Run(
			row.name,
			func(t *testing.T) {
				t.Parallel()

				size, err := sizer.MinSizeByName(row.typeName, row.def)
				if row.err != nil {
					require.ErrorIs(t, err, row.err, "%d of %d: %q: %+v", i+i, len(testMinSizeByNameRows), row.name, err)
				} else {
					require.NoError(t, err, "%d of %d: %q: %+v", i+i, len(testMinSizeByNameRows), row.name, err)
					require.Equal(t, row.size, size, "%d of %d: %q: %+v", i+i, len(testMinSizeByNameRows), row.name, size)
				}
			},
		)
	}
}
