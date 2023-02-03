//nolint:lll,ireturn,funlen
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
	def  []definition.Definition
	err  error
}

var testRows = []testRow{ //nolint:gochecknoglobals
	{
		name: "valid 1",
		def: []definition.Definition{
			buildDefinition(1),
		},
	},
	{
		name: "valid 2",
		def: []definition.Definition{
			buildDefinition(1),
			buildDefinition(2),
		},
	},
	{
		name: "proto_id not uniq",
		def: []definition.Definition{
			buildDefinition(1),
			buildDefinition(1),
		},
		err: validator.ErrDupID,
	},
	{
		name: "proto_id not set",
		def: []definition.Definition{
			buildDefinition(0),
		},
		err: validator.ErrNoProtoID,
	},
	{
		name: "no messages",
		def: []definition.Definition{
			func() definition.Definition {
				def := buildDefinition(1)
				def.Messages = nil

				return def
			}(),
		},
		err: validator.ErrNoMessages,
	},
	{
		name: "enum name is not unique",
		def: []definition.Definition{
			{
				Module: buildDefinition(1).Module,
				Proto:  buildDefinition(1).Proto,
				Enums: []definition.Enum{
					{Name: "enum1", BaseType: "uint8", Values: []definition.EnumValue{{Name: "v1", Value: "v1"}}},
					{Name: "enum1", BaseType: "uint8", Values: []definition.EnumValue{{Name: "v1", Value: "v1"}}},
				},
				Messages: buildDefinition(1).Messages,
			},
		},
		err: validator.ErrDupID,
	},
	{
		name: "enum value name is not unique",
		def: []definition.Definition{
			{
				Module: buildDefinition(1).Module,
				Proto:  buildDefinition(1).Proto,
				Enums: []definition.Enum{
					{Name: "enum1", BaseType: "uint8", Values: []definition.EnumValue{{Name: "v1", Value: "v1"}, {Name: "v1", Value: "v1"}}},
				},
				Messages: buildDefinition(1).Messages,
			},
		},
		err: validator.ErrDupID,
	},
	{
		name: "standard type redefined by enum",
		def: []definition.Definition{
			buildDefinition(1, definition.Definition{Enums: []definition.Enum{updateEnum(search(buildDefinition(0).Enums, func(enum definition.Enum) bool { return enum.Name == "Bool" }), definition.Enum{Name: "double"})}}),
		},
		err: validator.ErrRedefined,
	},
	{
		name: "invalid enum type",
		def: []definition.Definition{
			buildDefinition(1, definition.Definition{Enums: []definition.Enum{{Name: "Bool", BaseType: "invalid"}}}),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "no enum values",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Enums: []definition.Enum{{Name: "badEnum", BaseType: "double"}}},
			),
		},
		err: validator.ErrNoValues,
	},
	{
		name: "standard type redefined by message",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: []definition.Message{{Name: "double", ID: 100}}},
			),
		},
		err: validator.ErrRedefined,
	},
	{
		name: "enum redefined by message",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: []definition.Message{{Name: "Bool", ID: 100}}},
			),
		},
		err: validator.ErrRedefined,
	},
	{
		name: "duplicate message id",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: []definition.Message{{Name: "message7", ID: 1}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "no message id",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: []definition.Message{{Name: "message7"}}},
			),
		},
		err: validator.ErrNoID,
	},
	{
		name: "no message fields",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: []definition.Message{{Name: "message7", ID: 17}}},
			),
		},
		err: validator.ErrNoFields,
	},
	{
		name: "unknown field type",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Messages: []definition.Message{{Name: "message1", Fields: []definition.MessageField{{Type: definition.FieldType{Name: "invalid"}}}}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "no request in serving",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: []definition.ServicePair{{}}}},
			),
		},
		err: validator.ErrEmptyRequest,
	},
	{
		name: "no request in sending",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: []definition.ServicePair{{}}}},
			),
		},
		err: validator.ErrEmptyRequest,
	},
	{
		name: "bad request in serving",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: []definition.ServicePair{{Request: "invalid"}}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "bad response in serving",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: []definition.ServicePair{{Request: "message1", Response: "invalid"}}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "equal request and response in serving",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: []definition.ServicePair{{Request: "message1", Response: "message1"}}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "duplicate response in serving",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: []definition.ServicePair{{Request: "message1", Response: "message3"}}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "bad request in sending",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: []definition.ServicePair{{Request: "invalid"}}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "bad response in sending",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: []definition.ServicePair{{Request: "message4", Response: "invalid"}}}},
			),
		},
		err: validator.ErrUnknownType,
	},
	{
		name: "equal request and response in sending",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: []definition.ServicePair{{Request: "message4", Response: "message4"}}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "duplicate response in sending",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: []definition.ServicePair{{Request: "message4", Response: "message6"}}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "equal serving request and sending response",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Serving: []definition.ServicePair{{Request: "message6"}}}},
			),
		},
		err: validator.ErrDupID,
	},
	{
		name: "equal sending request and serving response",
		def: []definition.Definition{
			buildDefinition(
				1,
				definition.Definition{Service: definition.Service{Sending: []definition.ServicePair{{Request: "message3"}}}},
			),
		},
		err: validator.ErrDupID,
	},
}

func TestValidate(t *testing.T) {
	t.Parallel()

	for _, loopRow := range testRows { //nolint:paralleltest
		row := loopRow

		t.Run(
			row.name,
			func(innerT *testing.T) { testRun(innerT, row) },
		)
	}
}

func testRun(t *testing.T, row testRow) {
	t.Helper()
	t.Parallel()

	if err := validator.Validate(row.def); row.err != nil {
		require.ErrorIs(t, err, row.err, "validate %q: %+v", row.name, err)
	} else {
		require.NoError(t, err, "validate %q: %+v", row.name, err)
	}
}

func buildDefinition(protoID uint8, defs ...definition.Definition) definition.Definition {
	out := definition.Definition{
		Module: config.Module{Vendor: "vendor1", Protocol: "protocol1"},
		Proto:  definition.Proto{ID: protoID},
		Enums: []definition.Enum{
			{
				Name:     "Bool",
				BaseType: "uint8",
				Values: []definition.EnumValue{
					{Name: "true", Value: "1"},
					{Name: "false", Value: "0"},
				},
			},
			{
				Name:     "Number",
				BaseType: "uint8",
				Values: []definition.EnumValue{
					{Name: "one", Value: "1"},
					{Name: "zero", Value: "0"},
				},
			},
		},
		Messages: []definition.Message{
			{
				Name: "message10",
				ID:   10,
				Fields: []definition.MessageField{
					{Name: "field1", Type: definition.FieldType{Name: "uint8"}},
				},
			},
			{
				Name: "message1",
				ID:   1,
				Fields: []definition.MessageField{
					{Name: "field0", Type: definition.FieldType{Name: "message10"}},
					{Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					{Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					{Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
			{
				Name: "message2",
				ID:   2,
				Fields: []definition.MessageField{
					{Name: "field0", Type: definition.FieldType{Name: "message10"}},
					{Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					{Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					{Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
			{
				Name: "message3",
				ID:   3,
				Fields: []definition.MessageField{
					{Name: "field0", Type: definition.FieldType{Name: "message10"}},
					{Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					{Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					{Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
			{
				Name: "message4",
				ID:   4,
				Fields: []definition.MessageField{
					{Name: "field0", Type: definition.FieldType{Name: "message10"}},
					{Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					{Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					{Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
			{
				Name: "message5",
				ID:   5,
				Fields: []definition.MessageField{
					{Name: "field0", Type: definition.FieldType{Name: "message10"}},
					{Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					{Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					{Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
			{
				Name: "message6",
				ID:   6,
				Fields: []definition.MessageField{
					{Name: "field0", Type: definition.FieldType{Name: "message10"}},
					{Name: "field1", Type: definition.FieldType{Name: "uint8"}},
					{Name: "field2", Type: definition.FieldType{Name: "uint8", Array: true}},
					{Name: "field3", Type: definition.FieldType{Name: "uint8", Array: true, ArraySize: 10}},
				},
			},
		},
		Service: definition.Service{
			Serving: []definition.ServicePair{
				{Request: "message1"},
				{Request: "message2", Response: "message3"},
			},
			Sending: []definition.ServicePair{
				{Request: "message4"},
				{Request: "message5", Response: "message6"},
			},
		},
	}

	for _, def := range defs {
		out = updateDefinition(out, def)
	}

	return out
}

func updateDefinition(orig, update definition.Definition) definition.Definition {
	if update.Module.Vendor != "" {
		orig.Module.Vendor = update.Module.Vendor
	}

	if update.Module.Protocol != "" {
		orig.Module.Protocol = update.Module.Protocol
	}

	if update.Proto.ID != 0 {
		orig.Proto.ID = update.Proto.ID
	}

	for _, enum := range update.Enums {
		orig.Enums = upsert(
			orig.Enums,
			updateEnum(
				search(orig.Enums, func(e definition.Enum) bool { return e.Name == enum.Name }),
				enum,
			),
			func(e definition.Enum) bool { return e.Name == enum.Name },
		)
	}

	for _, msg := range update.Messages {
		orig.Messages = upsert(
			orig.Messages,
			updateMessage(
				search(orig.Messages, func(m definition.Message) bool { return m.Name == msg.Name }),
				msg,
			),
			func(m definition.Message) bool { return m.Name == msg.Name },
		)
	}

	for _, pair := range update.Service.Serving {
		orig.Service.Serving = upsert(
			orig.Service.Serving,
			updateServicePair(
				search(orig.Service.Serving, func(p definition.ServicePair) bool { return p.Request == pair.Request }),
				pair,
			),
			func(p definition.ServicePair) bool { return p.Request == pair.Request },
		)
	}

	for _, pair := range update.Service.Sending {
		orig.Service.Sending = upsert(
			orig.Service.Sending,
			updateServicePair(
				search(orig.Service.Sending, func(p definition.ServicePair) bool { return p.Request == pair.Request }),
				pair,
			),
			func(p definition.ServicePair) bool { return p.Request == pair.Request },
		)
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

	for _, value := range update.Values {
		orig.Values = upsert(
			orig.Values,
			updateEnumValue(
				search(orig.Values, func(v definition.EnumValue) bool { return v.Name == value.Name }),
				value,
			),
			func(v definition.EnumValue) bool { return v.Name == value.Name },
		)
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

	for _, v := range update.Fields {
		orig.Fields = upsert(
			orig.Fields,
			updateMessageField(
				search(orig.Fields, func(f definition.MessageField) bool { return f.Name == v.Name }),
				v,
			),
			func(f definition.MessageField) bool { return f.Name == v.Name },
		)
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

func upsert[V any](l []V, v V, f func(V) bool) []V {
	for i, u := range l {
		if f(u) {
			l[i] = v

			return l
		}
	}

	return append(l, v)
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

func search[V any](l []V, f func(V) bool) V {
	for _, v := range l {
		if f(v) {
			return v
		}
	}

	return *new(V) //nolint:gocritic
}
