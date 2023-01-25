package validator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/microavia/go-messgen/internal/config"
	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/validator"
)

type testRow struct {
	name string
	def  []*definition.Definition
	err  error
}

var testRows = []testRow{
	{
		name: "valid 1",
		def: []*definition.Definition{
			buildDefinition(1),
		},
	},
	{
		name: "valid 2",
		def: []*definition.Definition{
			buildDefinition(1),
			buildDefinition(2),
		},
	},
	{
		name: "proto_id not uniq",
		def: []*definition.Definition{
			buildDefinition(1),
			buildDefinition(1),
		},
		err: validator.ErrDupID,
	},
	{
		name: "proto_id not set",
		def: []*definition.Definition{
			buildDefinition(0),
		},
		err: validator.ErrNoProtoID,
	},
	{
		name: "no messages",
		def: []*definition.Definition{
			func() *definition.Definition {
				def := buildDefinition(1)
				def.Messages = nil
				return def
			}(),
		},
		err: validator.ErrNoMessages,
	},
	{
		name: "standard type redefined by enum",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Enums: definition.Enums{"double": buildDefinition(0).Enums["Bool"]}},
				definition.Definition{Enums: definition.Enums{"double": definition.Enum{Name: "double"}}},
			),
		},
		err: validator.ErrRedefined,
	},
	{
		name: "invalid enum type",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Enums: definition.Enums{"Bool": definition.Enum{BaseType: "invalid"}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "no enum values",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Enums: definition.Enums{"badEnum": definition.Enum{Name: "badEnum", BaseType: "double"}}},
			),
		},
		err: validator.ErrNoValues,
	},
	{
		name: "standard type redefined by message",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: definition.Messages{"double": buildDefinition(0).Messages["message1"]}},
				definition.Definition{Messages: definition.Messages{"double": definition.Message{Name: "double", ID: 10}}},
			),
		},
		err: validator.ErrRedefined,
	},
	{
		name: "enum redefined by message",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: definition.Messages{"Bool": definition.Message{Name: "Bool"}}},
			),
		},
		err: validator.ErrRedefined,
	},
	{
		name: "duplicate message id",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: definition.Messages{"message7": definition.Message{Name: "message7", ID: 1}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "no message id",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: definition.Messages{"message7": definition.Message{Name: "message7"}}},
			),
		},
		err: validator.ErrNoID,
	},
	{
		name: "no message fields",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: definition.Messages{"message7": definition.Message{Name: "message7", ID: 17}}},
			),
		},
		err: validator.ErrNoFields,
	},
	{
		name: "unknown field type",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: definition.Messages{"message1": definition.Message{Fields: definition.MessageFields{"field1": definition.MessageField{Type: definition.FieldType{Name: "invalid"}}}}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "no request in serving",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: definition.ServicePairs{"": {}}}},
			),
		},
		err: validator.ErrEmptyRequest,
	},
	{
		name: "no request in sending",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: definition.ServicePairs{"": {}}}},
			),
		},
		err: validator.ErrEmptyRequest,
	},
	{
		name: "bad request in serving",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: definition.ServicePairs{"invalid": {Request: "invalid"}}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "bad response in serving",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: definition.ServicePairs{"message1": {Request: "message1", Response: "invalid"}}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "equal request and response in serving",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: definition.ServicePairs{"message1": {Request: "message1", Response: "message1"}}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "duplicate response in serving",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: definition.ServicePairs{"message1": {Request: "message1", Response: "message3"}}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "bad request in sending",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: definition.ServicePairs{"invalid": {Request: "invalid"}}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "bad response in sending",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: definition.ServicePairs{"message4": {Request: "message4", Response: "invalid"}}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "equal request and response in sending",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: definition.ServicePairs{"message4": {Request: "message4", Response: "message4"}}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "duplicate response in sending",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: definition.ServicePairs{"message4": {Request: "message4", Response: "message6"}}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "equal serving request and sending response",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: definition.ServicePairs{"message6": {Request: "message6"}}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "equal sending request and serving response",
		def: []*definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: definition.ServicePairs{"message3": {Request: "message3"}}}},
			),
		},
		err: validator.ErrDupID,
	},
}

func TestValidate(t *testing.T) {
	for _, row := range testRows {
		t.Run(row.name, func(innerT *testing.T) { testRun(innerT, row) })
	}
}

func testRun(t *testing.T, row testRow) {
	t.Parallel()

	if err := validator.Validate(row.def); row.err != nil {
		require.ErrorIs(t, err, row.err, "validate %q: %+v/%+v: %+v", row.name, err)
	} else {
		require.NoError(t, err, "validate %q: %+v", row.name, err)
	}
}

func buildDefinition(protoID uint8, defs ...definition.Definition) *definition.Definition {
	out := definition.Definition{
		Module: config.Module{Vendor: "vendor1", Protocol: "protocol1"},
		Proto:  definition.Proto{ProtoID: protoID},
		Enums: definition.Enums{
			"Bool": {
				Name:     "Bool",
				BaseType: "uint8",
				Values: definition.EnumValues{
					"true":  {Name: "true", Value: "1"},
					"false": {Name: "false", Value: "0"},
				},
			},
			"Number": {
				Name:     "Number",
				BaseType: "uint8",
				Values: definition.EnumValues{
					"one":  {Name: "one", Value: "1"},
					"zero": {Name: "zero", Value: "0"},
				},
			},
		},
		Messages: definition.Messages{
			"message1": {
				Name: "message1",
				ID:   1,
				Fields: definition.MessageFields{
					"field1": {Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					"field2": {Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					"field3": {Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
			"message2": {
				Name: "message2",
				ID:   2,
				Fields: definition.MessageFields{
					"field1": {Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					"field2": {Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					"field3": {Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
			"message3": {
				Name: "message3",
				ID:   3,
				Fields: definition.MessageFields{
					"field1": {Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					"field2": {Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					"field3": {Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
			"message4": {
				Name: "message4",
				ID:   4,
				Fields: definition.MessageFields{
					"field1": {Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					"field2": {Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					"field3": {Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
			"message5": {
				Name: "message5",
				ID:   5,
				Fields: definition.MessageFields{
					"field1": {Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					"field2": {Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					"field3": {Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
			"message6": {
				Name: "message6",
				ID:   6,
				Fields: definition.MessageFields{
					"field1": {Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					"field2": {Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					"field3": {Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
		},
		Service: definition.Service{
			Serving: definition.ServicePairs{
				"message1": {Request: "message1"},
				"message2": {Request: "message2", Response: "message3"},
			},
			Sending: definition.ServicePairs{
				"message1": {Request: "message4"},
				"message2": {Request: "message5", Response: "message6"},
			},
		},
	}

	for _, def := range defs {
		out = updateDefinition(out, def)
	}

	return &out
}

func updateDefinition(orig, update definition.Definition) definition.Definition {
	if update.Module.Vendor != "" {
		orig.Module.Vendor = update.Module.Vendor
	}

	if update.Module.Protocol != "" {
		orig.Module.Protocol = update.Module.Protocol
	}

	if update.Proto.ProtoID != 0 {
		orig.Proto.ProtoID = update.Proto.ProtoID
	}

	for name, enum := range update.Enums {
		orig.Enums = appendMap(orig.Enums, name, updateEnum(orig.Enums[name], enum))
	}

	for name, msg := range update.Messages {
		orig.Messages = appendMap(orig.Messages, name, updateMessage(orig.Messages[name], msg))
	}

	for name, pair := range update.Service.Serving {
		orig.Service.Serving = appendMap(orig.Service.Serving, name, updateServicePair(orig.Service.Serving[name], pair))
	}

	for name, pair := range update.Service.Sending {
		orig.Service.Sending = appendMap(orig.Service.Sending, name, updateServicePair(orig.Service.Sending[name], pair))
	}

	return orig
}

func updateEnum(orig, update definition.Enum) definition.Enum {
	if update.Name != "" {
		orig.Name = update.Name
	}

	if update.BaseType != "" {
		orig.BaseType = update.BaseType
	}

	for name, value := range update.Values {
		orig.Values = appendMap(orig.Values, name, updateEnumValue(orig.Values[name], value))
	}

	return orig
}

func updateEnumValue(orig, update definition.EnumValue) definition.EnumValue {
	if update.Name != "" {
		orig.Name = update.Name
	}

	if update.Value != "" {
		orig.Value = update.Value
	}

	return orig
}

func updateMessage(orig, update definition.Message) definition.Message {
	if update.Name != "" {
		orig.Name = update.Name
	}

	if update.ID != 0 {
		orig.ID = update.ID
	}

	for name, v := range update.Fields {
		orig.Fields = appendMap(orig.Fields, name, updateMessageField(orig.Fields[name], v))
	}

	return orig
}

func updateMessageField(orig, update definition.MessageField) definition.MessageField {
	if update.Name != "" {
		orig.Name = update.Name
	}

	if update.Type.Name != "" {
		orig.Type.Name = update.Type.Name
	}

	if update.Type.Array {
		orig.Type.Array = update.Type.Array
	}

	if update.Type.ArraySize != 0 {
		orig.Type.ArraySize = update.Type.ArraySize
	}

	return orig
}

func appendMap[K comparable, V any](m map[K]V, k K, v V) map[K]V {
	if m == nil {
		m = make(map[K]V)
	}

	m[k] = v

	return m
}

func updateServicePair(orig, update definition.ServicePair) definition.ServicePair {
	if update.Request != "" {
		orig.Request = update.Request
	}

	if update.Response != "" {
		orig.Response = update.Response
	}

	return orig
}
