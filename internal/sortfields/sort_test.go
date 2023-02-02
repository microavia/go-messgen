//nolint:funlen
package sortfields_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/sortfields"
)

func TestSortFields(t *testing.T) {
	def := definition.Definition{
		Messages: []definition.Message{
			{
				Name: "message6",
				ID:   6,
				Fields: []definition.MessageField{
					{Name: "field2", Type: definition.FieldType{Name: "int32", Array: true, ArraySize: 7}},
					{Name: "field3", Type: definition.FieldType{Name: "int16", Array: true}},
					{Name: "field4", Type: definition.FieldType{Name: "message10", Array: true}},
					{Name: "field5", Type: definition.FieldType{Name: "string"}},
					{Name: "field6", Type: definition.FieldType{Name: "float64"}},
				},
			},
			{
				Name: "message10",
				ID:   10,
				Fields: []definition.MessageField{
					{Name: "field_char_slice", Type: definition.FieldType{Name: "char", Array: true}},
					{Name: "field_float32_slice", Type: definition.FieldType{Name: "float32", Array: true}},
					{Name: "field_float64_slice", Type: definition.FieldType{Name: "float64", Array: true}},
					{Name: "field_int16_slice", Type: definition.FieldType{Name: "int16", Array: true}},
					{Name: "field_int32_slice", Type: definition.FieldType{Name: "int32", Array: true}},
					{Name: "field_int64_slice", Type: definition.FieldType{Name: "int64", Array: true}},
					{Name: "field_uint16_slice", Type: definition.FieldType{Name: "uint16", Array: true}},
					{Name: "field_uint32_slice", Type: definition.FieldType{Name: "uint32", Array: true}},
					{Name: "field_uint64_slice", Type: definition.FieldType{Name: "uint64", Array: true}},
					{Name: "field_uint8_slice", Type: definition.FieldType{Name: "uint8", Array: true}},
				},
			},
		},
	}

	expected := definition.Definition{
		Messages: []definition.Message{
			{
				Name: "message6",
				ID:   6,
				Fields: []definition.MessageField{
					{Name: "field6", Type: definition.FieldType{Name: "float64"}},
					{Name: "field2", Type: definition.FieldType{Name: "int32", Array: true, ArraySize: 7}},
					{Name: "field4", Type: definition.FieldType{Name: "message10", Array: true}},
					{Name: "field3", Type: definition.FieldType{Name: "int16", Array: true}},
					{Name: "field5", Type: definition.FieldType{Name: "string"}},
				},
			},
			{
				Name: "message10",
				ID:   10,
				Fields: []definition.MessageField{
					{Name: "field_float64_slice", Type: definition.FieldType{Name: "float64", Array: true}},
					{Name: "field_int64_slice", Type: definition.FieldType{Name: "int64", Array: true}},
					{Name: "field_uint64_slice", Type: definition.FieldType{Name: "uint64", Array: true}},
					{Name: "field_float32_slice", Type: definition.FieldType{Name: "float32", Array: true}},
					{Name: "field_int32_slice", Type: definition.FieldType{Name: "int32", Array: true}},
					{Name: "field_uint32_slice", Type: definition.FieldType{Name: "uint32", Array: true}},
					{Name: "field_int16_slice", Type: definition.FieldType{Name: "int16", Array: true}},
					{Name: "field_uint16_slice", Type: definition.FieldType{Name: "uint16", Array: true}},
					{Name: "field_char_slice", Type: definition.FieldType{Name: "char", Array: true}},
					{Name: "field_uint8_slice", Type: definition.FieldType{Name: "uint8", Array: true}},
				},
			},
		},
	}

	t.Parallel()

	sortfields.SortFields(def)
	require.Equal(t, expected, def)
}
