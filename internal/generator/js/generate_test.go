//nolint:funlen
package js_test

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microavia/go-messgen/internal/config"
	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/generator/js"
)

//go:embed testdata/generatedold/vendor1/protocol1/*.json
var testDataJSON embed.FS

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
			{
				Name:     "Number",
				BaseType: "float64",
				Values: []definition.EnumValue{
					{Name: "Zero", Value: "0"},
					{Name: "One", Value: "1"},
				},
			},
		},
		Messages: []definition.Message{
			{
				Name: "message6",
				ID:   6,
				Fields: []definition.MessageField{
					{Name: "field1", Type: definition.FieldType{Name: "Bool"}},
					{Name: "field2", Type: definition.FieldType{Name: "int32", Array: true, ArraySize: 7}},
					{Name: "field3", Type: definition.FieldType{Name: "int16", Array: true}},
					{Name: "field4", Type: definition.FieldType{Name: "message10", Array: true}},
					{Name: "field5", Type: definition.FieldType{Name: "string"}},
					{Name: "field6", Type: definition.FieldType{Name: "float64"}},
					{Name: "field7", Type: definition.FieldType{Name: "Number"}},
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

	os.RemoveAll(outDir)
	defer os.RemoveAll(outDir)

	err := js.GenerateModule(outDir, def)
	require.NoError(t, err, "generating: %v", err)

	expected, err := loadJSONs(testDataJSON, filepath.Join("testdata/generatedold", def.Module.Vendor, def.Module.Protocol))
	require.NoError(t, err, "loading expected: %v", err)

	actual, err := loadJSONs(os.DirFS(outDir), filepath.Join(def.Module.Vendor, def.Module.Protocol))
	require.NoError(t, err, "loading actual: %v", err)
	require.Equal(t, expected, actual, "comparing")
}

func loadJSONs(fsys fs.FS, root string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	err := fs.WalkDir(
		fsys,
		root,
		func(path string, d fs.DirEntry, errInner error) error {
			if errInner != nil {
				return errInner
			}

			if d.IsDir() || !strings.HasSuffix(path, ".json") {
				return nil
			}

			var fileOut map[string]interface{}

			f, err := fsys.Open(path)
			if err != nil {
				return fmt.Errorf("opening %s: %w", path, err)
			}

			defer f.Close()

			de := json.NewDecoder(f)
			de.DisallowUnknownFields()

			err = de.Decode(&fileOut)
			if err != nil {
				return fmt.Errorf("reading %s: %w", path, err)
			}

			for k, v := range fileOut {
				if k == "version" {
					v = "UNKNOWN"
				}
				out[k] = v
			}

			return nil
		},
	)
	if err != nil { //nolint:wsl
		return nil, fmt.Errorf("loading from %q: %w", root, err)
	}

	return out, nil
}
