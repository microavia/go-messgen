package definition

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/microavia/go-messgen/internal/stdtypes"
)

var ErrTypeUnknown = errors.New("type unknown")

type Enum struct {
	Name     string      `json:"name"`
	BaseType string      `json:"basetype"`
	Values   []EnumValue `json:"fields"`
}

type EnumValue struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"descr"`
}

func (e *Enum) UnmarshalJSON(b []byte) error {
	var rawEnum struct {
		Name     string      `json:"name"`
		BaseType string      `json:"basetype"`
		Values   []EnumValue `json:"fields"`
	}

	if err := json.Unmarshal(b, &rawEnum); err != nil {
		return err
	}

	if _, ok := stdtypes.StdTypes[rawEnum.BaseType]; !ok {
		return fmt.Errorf("%q: %w", rawEnum.BaseType, ErrTypeUnknown)
	}

	*e = Enum{
		Name:     rawEnum.Name,
		BaseType: rawEnum.BaseType,
		Values:   make([]EnumValue, 0),
	}

	return nil
}
