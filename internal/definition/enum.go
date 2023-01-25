package definition

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrTypeUnknown = errors.New("type unknown")

type Enum struct {
	Name     string     `json:"name"`
	BaseType string     `json:"basetype"`
	Values   EnumValues `json:"fields"`
}

type Enums map[string]Enum

func (e *Enums) UnmarshalJSON(b []byte) error {
	return MapUnmarshalJSON(
		b,
		(*map[string]Enum)(e),
		func(v Enum) (Enum, error) { return v, nil },
		func(v Enum) string { return v.Name },
	)
}

func (e *Enums) MarshalJSON() ([]byte, error) {
	return json.Marshal(MapToSlice(*e, func(s1, s2 string) bool { return s1 < s2 }))
}

type EnumValue struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"descr"`
}

type EnumValues map[string]EnumValue

func (e *EnumValues) UnmarshalJSON(b []byte) error {
	type rawValue struct {
		Name        interface{} `json:"name"`
		Value       interface{} `json:"value"`
		Description string      `json:"descr"`
	}

	return MapUnmarshalJSON(
		b,
		(*map[string]EnumValue)(e),
		func(v rawValue) (EnumValue, error) {
			return EnumValue{
				Name:        fmt.Sprintf("%v", v.Name),
				Value:       fmt.Sprintf("%v", v.Value),
				Description: v.Description,
			}, nil
		},
		func(v EnumValue) string { return v.Name },
	)
}

func (e *EnumValues) MarshalJSON() ([]byte, error) {
	return json.Marshal(MapToSlice(*e, func(s1, s2 string) bool { return s1 < s2 }))
}
