//nolint:gochecknoglobals,funlen
package golang_test

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"

	"github.com/microavia/go-messgen/internal/config"
	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/generator/golang"
)

//go:embed testdata/templates/compare_test.go.tmpl
var tmplSrc []byte

var tmplCompiled = map[string]*template.Template{
	"compare_test.go": template.Must(template.New("compare_test.go.tmpl").Parse(string(tmplSrc))),
}

//go:embed testdata/message6.old.bin
var message6OldBin []byte

const filePerms = 0600 //nolint:gofumpt

func TestGenerateModule(t *testing.T) {
	t.Parallel()

	outDir := "./test.tmp/generated"
	def := definition.Definition{
		Module: config.Module{Vendor: "vendor1", Protocol: "protocol1"},
		Proto:  definition.Proto{ID: 1},
		Enums: []definition.Enum{
			{
				Name:     "Bool",
				BaseType: "uint8",
				Values: []definition.EnumValue{
					{Name: "False", Value: "0"},
					{Name: "True", Value: "1"},
				},
			},
		},
		Messages: []definition.Message{
			{
				Name: "message1",
				ID:   1,
				Fields: []definition.MessageField{
					{Name: "field1", Type: definition.FieldType{Name: "Bool"}},
				},
			},
			{
				Name: "message2",
				ID:   2,
				Fields: []definition.MessageField{
					{Name: "field1", Type: definition.FieldType{Name: "Bool"}},
					{Name: "field2", Type: definition.FieldType{Name: "int32", Array: true, ArraySize: 7}},
				},
			},
			{
				Name: "message3",
				ID:   3,
				Fields: []definition.MessageField{
					{Name: "field1", Type: definition.FieldType{Name: "Bool"}},
					{Name: "field2", Type: definition.FieldType{Name: "int32", Array: true, ArraySize: 7}},
					{Name: "field3", Type: definition.FieldType{Name: "int16", Array: true}},
				},
			},
			{
				Name: "message4",
				ID:   4,
				Fields: []definition.MessageField{
					{Name: "field1", Type: definition.FieldType{Name: "Bool"}},
					{Name: "field2", Type: definition.FieldType{Name: "int32", Array: true, ArraySize: 7}},
					{Name: "field3", Type: definition.FieldType{Name: "int16", Array: true}},
					{Name: "field4", Type: definition.FieldType{Name: "message10", Array: true}},
				},
			},
			{
				Name: "message5",
				ID:   5,
				Fields: []definition.MessageField{
					{Name: "field1", Type: definition.FieldType{Name: "Bool"}},
					{Name: "field2", Type: definition.FieldType{Name: "int32", Array: true, ArraySize: 7}},
					{Name: "field3", Type: definition.FieldType{Name: "int16", Array: true}},
					{Name: "field4", Type: definition.FieldType{Name: "message10", Array: true}},
					{Name: "field5", Type: definition.FieldType{Name: "string"}},
				},
			},
			{
				Name: "message6",
				ID:   6,
				Fields: []definition.MessageField{
					{Name: "field1", Type: definition.FieldType{Name: "Bool"}},
					{Name: "field2", Type: definition.FieldType{Name: "int32", Array: true, ArraySize: 7}},
					{Name: "field3", Type: definition.FieldType{Name: "int16", Array: true}},
					{Name: "field4", Type: definition.FieldType{Name: "message10", Array: true}},
					{Name: "field5", Type: definition.FieldType{Name: "string"}},
					{Name: "field6", Type: definition.FieldType{Name: "double"}},
				},
			},
			{
				Name: "message10",
				ID:   10,
				Fields: []definition.MessageField{
					{Name: "field_char_slice", Type: definition.FieldType{Name: "char", Array: true}},
					{Name: "field_double_slice", Type: definition.FieldType{Name: "double", Array: true}},
					{Name: "field_float_slice", Type: definition.FieldType{Name: "float", Array: true}},
					{Name: "field_float32_slice", Type: definition.FieldType{Name: "float32", Array: true}},
					{Name: "field_float64_slice", Type: definition.FieldType{Name: "float64", Array: true}},
					{Name: "field_int16_slice", Type: definition.FieldType{Name: "int16", Array: true}},
					{Name: "field_int32_slice", Type: definition.FieldType{Name: "int32", Array: true}},
					{Name: "field_int64_slice", Type: definition.FieldType{Name: "int64", Array: true}},
					{Name: "field_int8_slice", Type: definition.FieldType{Name: "int8", Array: true}},
					{Name: "field_uint16_slice", Type: definition.FieldType{Name: "uint16", Array: true}},
					{Name: "field_uint32_slice", Type: definition.FieldType{Name: "uint32", Array: true}},
					{Name: "field_uint64_slice", Type: definition.FieldType{Name: "uint64", Array: true}},
					{Name: "field_uint8_slice", Type: definition.FieldType{Name: "uint8", Array: true}},
					{Name: "field_string_slice", Type: definition.FieldType{Name: "string", Array: true}},
					{Name: "field1", Type: definition.FieldType{Name: "Bool", Array: true}},
					{Name: "field_float64_slice_empty", Type: definition.FieldType{Name: "float64", Array: true}},
				},
			},
		},
	}

	defer os.RemoveAll(outDir)

	err := golang.GenerateModule(outDir, def)
	require.NoError(t, err, "generating: %v", err)

	b, err := exec.Command("go", "test", "./"+filepath.Join(outDir, "...")).CombinedOutput() //nolint:gosec
	require.NoError(t, err, "testing: %v: %q", err, b)
}

func TestCompare(t *testing.T) { //nolint:paralleltest
	outDir := "./test.tmp/compare"

	def := definition.Definition{
		Module: config.Module{Vendor: "compare", Protocol: "current"},
		Proto:  definition.Proto{ID: 1},
		Messages: []definition.Message{
			{
				Name: "message6",
				ID:   6,
				Fields: []definition.MessageField{
					{Name: "field2", Type: definition.FieldType{Name: "int32", Array: true, ArraySize: 7}},
					{Name: "field3", Type: definition.FieldType{Name: "int16", Array: true}},
					{Name: "field4", Type: definition.FieldType{Name: "message10", Array: true}},
					{Name: "field5", Type: definition.FieldType{Name: "string"}},
					{Name: "field6", Type: definition.FieldType{Name: "double"}},
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

	defer os.RemoveAll(outDir)

	err := golang.GenerateModule(outDir, def)
	require.NoError(t, err, "generating: %v", err)

	err = golang.GenerateModuleByTemplates(outDir, def, tmplCompiled)
	require.NoError(t, err, "generating additional tests: %v", err)

	err = os.WriteFile(
		filepath.Join(outDir, def.Module.Vendor, def.Module.Protocol, "message6.old.bin"),
		message6OldBin,
		filePerms,
	)
	require.NoError(t, err, "writing predefined message: %v", err)

	b, err := exec.Command("go", "test", "./"+filepath.Join(outDir, "...")).CombinedOutput() //nolint:gosec
	require.NoError(t, err, "testing: %v: %q", err, b)
}
