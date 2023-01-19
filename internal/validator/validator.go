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
	ErrUnknown       = fmt.Errorf("unknown: %w", ErrBadDefinition)
	ErrUnknownType   = fmt.Errorf("type: %w", ErrUnknown)
	ErrRedefined     = fmt.Errorf("redefined: %w", ErrBadDefinition)
	ErrEmptyRequest  = fmt.Errorf("empty request: %w", ErrBadDefinition)
)

func Validate(modules []*definition.Definition) error {
	stdTypesSet := buildSet(stdtypes.Types)

	if err := checkUniq(modules, func(v *definition.Definition) uint8 { return v.Proto.ProtoID }); err != nil {
		return fmt.Errorf("checking proto_id uniqueness: %w", err)
	}

	for moduleID, module := range modules {
		if module.Proto.ProtoID == 0 {
			return fmt.Errorf("%+v: %w", moduleID, ErrNoProtoID)
		}

		if len(module.Messages) == 0 {
			return fmt.Errorf("%+v: %w", moduleID, ErrNoMessages)
		}

		if err := validateConstants(module.Constants, stdTypesSet); err != nil {
			return fmt.Errorf("%+v: %w", moduleID, err)
		}

		err := validateMessages(
			module.Messages,
			stdTypesSet,
			buildSet(buildMap(module.Constants, func(v definition.Constant) string { return v.Name })),
		)
		if err != nil {
			return fmt.Errorf("%+v: %w", moduleID, err)
		}

		err = validateService(
			module.Service,
			buildSet(buildMap(module.Messages, func(v definition.Message) string { return v.Name })),
		)
		if err != nil {
			return fmt.Errorf("%+v: %w", moduleID, err)
		}
	}

	return nil
}

func validateConstants(
	constants []definition.Constant,
	stdTypes map[string]struct{},
) error {
	if err := checkUniq(constants, func(v definition.Constant) string { return v.Name }); err != nil {
		return fmt.Errorf("checking constants uniqueness: %w", err)
	}

	for _, c := range constants {
		if err := validateConstant(c, stdTypes); err != nil {
			return fmt.Errorf("constant %q: %w", c.Name, err)
		}
	}

	return nil
}

func validateConstant(c definition.Constant, stdTypes map[string]struct{}) error {
	if checkPresence(stdTypes, c.Name) {
		return fmt.Errorf("standard type redefined: %q: %w", c.Name, ErrRedefined)
	}

	if !checkPresence(stdTypes, c.BaseType) {
		return fmt.Errorf("constant %+v: base type: %w", c, ErrUnknownType)
	}

	if err := checkUniq(c.Fields, func(v definition.ConstantField) string { return v.Name }); err != nil {
		return fmt.Errorf("checking fields uniqueness: %w", err)
	}

	return nil
}

func validateMessages(
	messages []definition.Message,
	stdTypes map[string]struct{},
	constants map[string]struct{},
) error {
	if err := checkUniq(messages, func(v definition.Message) int { return v.ID }); err != nil {
		return fmt.Errorf("checking message_id uniqueness: %w", ErrDupID)
	}

	for _, msg := range messages {
		switch {
		case checkPresence(stdTypes, msg.Name):
			return fmt.Errorf("standard type redefined: %q: %w", msg.Name, ErrRedefined)
		case checkPresence(constants, msg.Name):
			return fmt.Errorf("constant type redefined: %q: %w", msg.Name, ErrRedefined)
		}

		if err := validateMessage(msg, mergeSets(stdTypes, constants)); err != nil {
			return fmt.Errorf("%q: %w", msg.Name, err)
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

	if err := checkUniq(msg.Fields, func(v definition.MessageField) string { return v.Name }); err != nil {
		return fmt.Errorf("checking fields uniqueness: %w", err)
	}

	for _, field := range msg.Fields {
		if !checkPresence(types, string(field.Type.Name)) {
			return fmt.Errorf("field %+v: %w", field, ErrUnknownType)
		}
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

func buildSet[K comparable, V any](m map[K]V) map[K]struct{} {
	set := make(map[K]struct{}, len(m))

	for k := range m {
		set[k] = struct{}{}
	}

	return set
}

func mergeSets(sets ...map[string]struct{}) map[string]struct{} {
	set := make(map[string]struct{})

	for _, s := range sets {
		for k := range s {
			set[k] = struct{}{}
		}
	}

	return set
}

func buildMap[K comparable, V any](in []V, f func(v V) K) map[K]V {
	out := make(map[K]V, len(in))

	for _, v := range in {
		out[f(v)] = v
	}

	return out
}

func checkUniq[V any, I comparable](in []V, f func(V) I) error {
	set := make(map[I]int, len(in))

	for k, v := range in {
		if existing, ok := set[f(v)]; ok {
			return fmt.Errorf(
				"%+v: duplicate id %+v in %+v: %w",
				existing,
				f(v),
				k,
				ErrDupID,
			)
		}

		set[f(v)] = k
	}

	return nil
}

func validateService(
	svc definition.Service,
	types map[string]struct{},
) error {
	if err := checkServicePairs(svc.Serving, types); err != nil {
		return fmt.Errorf("service: serving: %w", err)
	}

	return nil
}

func checkServicePairs(pairs []definition.ServicePair, types map[string]struct{}) error {
	for _, pair := range pairs {
		switch {
		case pair.Request == "":
			return fmt.Errorf("%+v: %w", pair, ErrEmptyRequest)
		case !checkPresence(types, pair.Request):
			return fmt.Errorf("%+v: request: %w", pair, ErrUnknownType)
		case pair.Response != "" && !checkPresence(types, pair.Response):
			return fmt.Errorf("%+v: response: %w", pair, ErrUnknownType)
		case pair.Request == pair.Response:
			return fmt.Errorf("same message as request and response for %+v: %w", pair, ErrDupID)
		}
	}

	if err := checkUniq(pairs, func(v definition.ServicePair) string { return v.Request }); err != nil {
		return fmt.Errorf("checking requests uniqueness: %w", err)
	}

	err := checkUniq(
		pairs,
		func(v definition.ServicePair) string {
			return ternary(v.Response != "", v.Response, v.Request)
		},
	)
	if err != nil {
		return fmt.Errorf("checking response uniqueness: %w", err)
	}

	return nil
}

func ternary[T any](cond bool, a, b T) T {
	if cond {
		return a
	}

	return b
}
