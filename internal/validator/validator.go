package validator

import (
	"errors"
	"fmt"

	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/stdtypes"
)

var (
	ErrBadDefinition = errors.New("error in definition")
	ErrNoID          = fmt.Errorf("no id: %w", ErrBadDefinition)
	ErrNoProtoID     = fmt.Errorf("protocol: %w", ErrNoID)
	ErrNoMsgID       = fmt.Errorf("message: %w", ErrNoID)
	ErrDupID         = fmt.Errorf("duplicate id: %w", ErrBadDefinition)
	ErrNoMessages    = fmt.Errorf("no messages: %w", ErrBadDefinition)
	ErrNoFields      = fmt.Errorf("no fields: %w", ErrBadDefinition)
	ErrNoValues      = fmt.Errorf("no values: %w", ErrBadDefinition)
	ErrUnknown       = fmt.Errorf("unknown: %w", ErrBadDefinition)
	ErrUnknownType   = fmt.Errorf("type: %w", ErrUnknown)
	ErrRedefined     = fmt.Errorf("redefined: %w", ErrBadDefinition)
	ErrEmptyRequest  = fmt.Errorf("empty request: %w", ErrBadDefinition)
)

func Validate(modules []definition.Definition) error {
	err := checkUniq(modules, func(v definition.Definition) uint8 { return v.Proto.ID })
	if err != nil {
		return fmt.Errorf("proto_id is not unique: %w", err)
	}

	for _, module := range modules {
		if module.Proto.ID == 0 {
			return fmt.Errorf("%+v: %w", module.Module, ErrNoProtoID)
		}

		if len(module.Messages) == 0 {
			return fmt.Errorf("%+v: %w", module.Module, ErrNoMessages)
		}

		if err = validateEnums(module.Enums, stdtypes.StdTypes); err != nil {
			return fmt.Errorf("validating constants: %+v: %w", module.Module, err)
		}

		if err = validateMessages(module.Messages, stdtypes.StdTypes, module.Enums); err != nil {
			return fmt.Errorf("validating messages: %+v: %w", module.Module, err)
		}

		if err = validateService(module.Service, module.Messages); err != nil {
			return fmt.Errorf("validating service: %+v: %w", module.Module, err)
		}
	}

	return nil
}

func validateEnums(
	enums []definition.Enum,
	stdTypes map[string]stdtypes.StdType,
) error {
	if err := checkUniq(enums, func(v definition.Enum) string { return v.Name }); err != nil {
		return fmt.Errorf("enum name is not unique: %w", err)
	}

	for _, c := range enums {
		if err := validateEnum(c, stdTypes); err != nil {
			return fmt.Errorf("enum %q: %w", c.Name, err)
		}
	}

	return nil
}

func validateEnum(c definition.Enum, stdTypes map[string]stdtypes.StdType) error {
	if len(c.Values) == 0 {
		return ErrNoValues
	}

	if err := checkUniq(c.Values, func(v definition.EnumValue) string { return v.Name }); err != nil {
		return fmt.Errorf("enum value name is not unique: %w", err)
	}

	if checkPresence(stdTypes, c.Name) {
		return fmt.Errorf("standard type redefined: %q: %w", c.Name, ErrRedefined)
	}

	if !checkPresence(stdTypes, c.BaseType) {
		return fmt.Errorf("constant %+v: base type: %w", c, ErrUnknownType)
	}

	return nil
}

func validateMessages(
	messages []definition.Message,
	stdTypes map[string]stdtypes.StdType,
	enums []definition.Enum,
) error {
	if err := checkUniq(messages, func(v definition.Message) uint16 { return v.ID }); err != nil {
		return fmt.Errorf("message id is not unique: %w", err)
	}

	var (
		enumsMap    = buildMap(enums, func(v definition.Enum) string { return v.Name })
		messagesMap = buildMap(messages, func(v definition.Message) string { return v.Name })
	)

	for _, msg := range messages {
		switch {
		case checkPresence(stdTypes, msg.Name):
			return fmt.Errorf("standard type redefined: %q: %w", msg.Name, ErrRedefined)
		case checkPresence(enumsMap, msg.Name):
			return fmt.Errorf("enum redefined: %q: %w", msg.Name, ErrRedefined)
		}

		err := validateMessage(msg, mergeSets(buildSet(stdtypes.StdTypes), buildSet(enumsMap), buildSet(messagesMap)))
		if err != nil {
			return fmt.Errorf("validating message: %q: %w", msg.Name, err)
		}
	}

	return nil
}

func validateMessage(
	msg definition.Message,
	types map[string]struct{},
) error {
	if msg.ID == 0 {
		return ErrNoMsgID
	}

	if len(msg.Fields) == 0 {
		return ErrNoFields
	}

	for _, field := range msg.Fields {
		if !checkPresence(types, field.Type.Name) {
			return fmt.Errorf("field %+v: %w", field, ErrUnknownType)
		}
	}

	return nil
}

func validateService(
	svc definition.Service,
	messages []definition.Message,
) error {
	if err := checkServicePairs(svc.Serving, messages); err != nil {
		return fmt.Errorf("service: serving: %w", err)
	}

	if err := checkServicePairs(svc.Sending, messages); err != nil {
		return fmt.Errorf("service: sending: %w", err)
	}

	for _, pair := range svc.Serving {
		if pair.Response != "" {
			requests := buildMap(svc.Sending, func(v definition.ServicePair) string { return v.Request })
			if checkPresence(requests, pair.Response) {
				return fmt.Errorf("service: serving: response %q is used as sending request: %w", pair.Response, ErrDupID)
			}
		}
	}

	for _, pair := range svc.Sending {
		requests := buildMap(svc.Serving, func(v definition.ServicePair) string { return v.Request })
		if pair.Response != "" {
			if checkPresence(requests, pair.Response) {
				return fmt.Errorf("service: sending: response %q is used as serving request: %w", pair.Response, ErrDupID)
			}
		}
	}

	return nil
}

func checkServicePairs(
	pairs []definition.ServicePair,
	messages []definition.Message,
) error {
	messagesMap := buildMap(messages, func(v definition.Message) string { return v.Name })
	responses := make(map[string]definition.ServicePair, len(pairs))

	for _, pair := range pairs {
		switch {
		case pair.Request == "":
			return fmt.Errorf("%+v: %w", pair, ErrEmptyRequest)
		case !checkPresence(messagesMap, pair.Request):
			return fmt.Errorf("%+v: request: %w", pair, ErrUnknownType)
		case pair.Response != "" && !checkPresence(messagesMap, pair.Response):
			return fmt.Errorf("%+v: response: %w", pair, ErrUnknownType)
		case pair.Request == pair.Response:
			return fmt.Errorf("same message as request and response for %+v: %w", pair, ErrDupID)
		case pair.Response != "" && checkPresence(responses, pair.Response):
			return fmt.Errorf("same message as response for %+v and %+v: %w", pair, responses[pair.Response], ErrDupID)
		case pair.Response != "":
			responses[pair.Response] = pair
		}
	}

	return nil
}

func buildSet[K comparable, V any](m map[K]V) map[K]struct{} {
	s := make(map[K]struct{}, len(m))
	for k := range m {
		s[k] = struct{}{}
	}

	return s
}

func buildMap[K comparable, V any](in []V, f func(v V) K) map[K]V {
	out := make(map[K]V, len(in))

	for _, v := range in {
		out[f(v)] = v
	}

	return out
}

func mergeSets[K comparable](sets ...map[K]struct{}) map[K]struct{} {
	s := make(map[K]struct{})

	for _, set := range sets {
		for k := range set {
			s[k] = struct{}{}
		}
	}

	return s
}

func checkPresence[K comparable, V any](m map[K]V, k K) bool {
	_, ok := m[k]

	return ok
}

func checkUniq[V any, I comparable](in []V, f func(V) I) error {
	set := make(map[I]V, len(in))

	for k, v := range in {
		i := f(v)

		if _, ok := set[i]; ok {
			return fmt.Errorf("%+v in %+v of %d: %w", i, k, len(in), ErrDupID)
		}

		set[i] = v
	}

	return nil
}
