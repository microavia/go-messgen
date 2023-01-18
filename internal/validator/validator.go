package validator

import (
	"errors"
	"fmt"

	"github.com/microavia/go-messgen/internal/config"
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
)

func Validate(modules map[config.Module]*definition.Definition) error {
	stdTypesSet := buildSet(stdtypes.Types)

	if err := checkUniqInMap(modules, func(v *definition.Definition) uint8 { return v.Proto.ProtoID }); err != nil {
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

		if err := validateMessages(module.Messages, stdTypesSet, buildSet(constantsMap(module.Constants))); err != nil {
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
	messages map[string]definition.Message,
	stdTypes map[string]struct{},
	constants map[string]struct{},
) error {
	if err := checkUniqInMap(messages, func(v definition.Message) int { return v.ID }); err != nil {
		return fmt.Errorf("checking message_id uniqueness: %w", ErrDupID)
	}

	for name, msg := range messages {
		switch {
		case checkPresence(stdTypes, name):
			return fmt.Errorf("standard type redefined: %q: %w", name, ErrRedefined)
		case checkPresence(constants, name):
			return fmt.Errorf("constant type redefined: %q: %w", name, ErrRedefined)
		}

		if err := validateMessage(msg, mergeSets(stdTypes, constants)); err != nil {
			return fmt.Errorf("%q: %w", name, err)
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
		if !checkPresence(types, field.Type) {
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

func buildSet[V any](m map[string]V) map[string]struct{} {
	set := make(map[string]struct{}, len(m))

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

func constantsMap(in []definition.Constant) map[string]definition.Constant {
	out := make(map[string]definition.Constant, len(in))

	for _, c := range in {
		out[c.Name] = c
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

func checkUniqInMap[K comparable, V any, I comparable](in map[K]V, f func(V) I) error {
	set := make(map[I]K, len(in))

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
