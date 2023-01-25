package validator

import (
	"errors"
	"fmt"
	"log"

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
	ErrUnknown       = fmt.Errorf("unknown: %w", ErrBadDefinition)
	ErrUnknownType   = fmt.Errorf("type: %w", ErrUnknown)
	ErrRedefined     = fmt.Errorf("redefined: %w", ErrBadDefinition)
	ErrEmptyRequest  = fmt.Errorf("empty request: %w", ErrBadDefinition)
)

func Validate(modules []*definition.Definition) error {
	err := checkUniq(modules, func(v *definition.Definition) uint8 { return v.Proto.ProtoID })
	if err != nil {
		return fmt.Errorf("checking proto_id uniqueness: %w", err)
	}

	for _, module := range modules {
		if module.Proto.ProtoID == 0 {
			return fmt.Errorf("%+v: %w", module.Module, ErrNoProtoID)
		}

		if len(module.Messages) == 0 {
			return fmt.Errorf("%+v: %w", module.Module, ErrNoMessages)
		}

		if err = validateConstants(module.Enums, stdtypes.StdTypes); err != nil {
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

func validateConstants(
	constants map[string]definition.Enum,
	stdTypes map[string]stdtypes.StdType,
) error {
	for _, c := range constants {
		if err := validateConstant(c, stdTypes); err != nil {
			return fmt.Errorf("constant %q: %w", c.Name, err)
		}
	}

	return nil
}

func validateConstant(c definition.Enum, stdTypes map[string]stdtypes.StdType) error {
	if checkPresence(stdTypes, c.Name) {
		return fmt.Errorf("standard type redefined: %q: %w", c.Name, ErrRedefined)
	}

	if !checkPresence(stdTypes, c.BaseType) {
		return fmt.Errorf("constant %+v: base type: %w", c, ErrUnknownType)
	}

	return nil
}

func validateMessages(
	messages definition.Messages,
	stdTypes map[string]stdtypes.StdType,
	enums definition.Enums,
) error {
	for _, msg := range messages {
		switch {
		case checkPresence(stdTypes, msg.Name):
			return fmt.Errorf("standard type redefined: %q: %w", msg.Name, ErrRedefined)
		case checkPresence(enums, msg.Name):
			return fmt.Errorf("constant type redefined: %q: %w", msg.Name, ErrRedefined)
		}

		err := checkUniqMapValues(messages, func(v definition.Message) uint16 { return v.ID })
		if err != nil {
			return fmt.Errorf("validating message: %q: %w", msg.Name, err)
		}

		err = validateMessage(msg, mergeSets(buildSet(stdtypes.StdTypes), buildSet(enums), buildSet(messages)))
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

	for _, field := range msg.Fields {
		if !checkPresence(types, field.Type.Name) {
			return fmt.Errorf("field %+v: %w", field, ErrUnknownType)
		}
	}

	return nil
}

func validateService(
	svc definition.Service,
	messages definition.Messages,
) error {
	log.Printf("checking service: %+v", svc)

	if err := checkServicePairs(svc.Serving, messages); err != nil {
		return fmt.Errorf("service: serving: %w", err)
	}

	if err := checkServicePairs(svc.Sending, messages); err != nil {
		return fmt.Errorf("service: sending: %w", err)
	}

	for _, pair := range svc.Serving {
		if pair.Response != "" {
			if checkPresence(svc.Sending, pair.Response) {
				return fmt.Errorf("service: serving: response %q is used as sending request: %w", pair.Response, ErrDupID)
			}
		}
	}

	for _, pair := range svc.Sending {
		if pair.Response != "" {
			if checkPresence(svc.Serving, pair.Response) {
				return fmt.Errorf("service: sending: response %q is used as serving request: %w", pair.Response, ErrDupID)
			}
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

func mergeSets[K comparable](sets ...map[K]struct{}) map[K]struct{} {
	s := make(map[K]struct{})
	for _, set := range sets {
		for k := range set {
			s[k] = struct{}{}
		}
	}

	return s
}

func checkServicePairs(
	pairs map[string]definition.ServicePair,
	messages definition.Messages,
) error {
	for _, pair := range pairs {
		log.Printf("checking pair: %+v", pair)
		switch {
		case pair.Request == "":
			return fmt.Errorf("%+v: %w", pair, ErrEmptyRequest)
		case !checkPresence(messages, pair.Request):
			return fmt.Errorf("%+v: request: %w", pair, ErrUnknownType)
		case pair.Response != "" && !checkPresence(messages, pair.Response):
			return fmt.Errorf("%+v: response: %w", pair, ErrUnknownType)
		case pair.Request == pair.Response:
			return fmt.Errorf("same message as request and response for %+v: %w", pair, ErrDupID)
		}
	}

	err := checkUniqMapValues(
		filterMap(
			pairs,
			func(_ string, v1 definition.ServicePair) bool {
				return v1.Response != ""
			},
		),
		func(v definition.ServicePair) string { return v.Response },
	)
	if err != nil {
		return fmt.Errorf("checking response uniqueness: %w", err)
	}

	return nil
}

func checkPresence[K comparable, V any](m map[K]V, k K) bool {
	if len(m) == 0 {
		return false
	}

	_, ok := m[k]

	return ok
}

func filterMap[K comparable, V any](m map[K]V, f func(K, V) bool) map[K]V {
	out := make(map[K]V, len(m))

	for k, v := range m {
		if f(k, v) {
			out[k] = v
		}
	}

	return out
}

func checkUniq[V any, I comparable](in []V, f func(V) I) error {
	set := make(map[I]V, len(in))

	for k, v := range in {
		i := f(v)

		if existing, ok := set[i]; ok {
			return fmt.Errorf("%+v: duplicate id %+v in %+v: %w", existing, i, k, ErrDupID)
		}

		set[i] = v
	}

	return nil
}

func checkUniqMapValues[K comparable, V any, I comparable](in map[K]V, f func(V) I) error {
	set := make(map[I]V, len(in))

	for k, v := range in {
		i := f(v)

		if existing, ok := set[i]; ok {
			return fmt.Errorf("%+v: duplicate id %+v in %+v: %w", existing, i, k, ErrDupID)
		}

		set[i] = v
	}

	return nil
}
