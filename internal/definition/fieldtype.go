package definition

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var ErrInvalidInput = errors.New("invalid input")

type TypeName string

type FieldType struct {
	Name      TypeName
	Array     bool
	ArraySize int
}

var parseType = regexp.MustCompile(`^\"([^\[]+)(\[(\d*)\])?\"$`)

func (t *FieldType) UnmarshalJSON(b []byte) error {
	typeAndSize := parseType.FindSubmatch(b)
	if typeAndSize == nil {
		return fmt.Errorf("invalid input: %q: %w", b, ErrInvalidInput)
	}

	*t = FieldType{
		Name:      TypeName(typeAndSize[1]),
		Array:     len(typeAndSize[2]) > 0,
		ArraySize: mustA2I(typeAndSize[3]),
	}

	return nil
}

func (t FieldType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", t.String())), nil
}

func (t FieldType) String() string {
	if t.Array {
		if t.ArraySize == 0 {
			return string(t.Name) + "[]"
		}

		return string(t.Name) + "[" + strconv.Itoa(t.ArraySize) + "]"
	}

	return string(t.Name)
}

func mustA2I(b []byte) int {
	if len(b) == 0 {
		return 0
	}

	n, err := strconv.Atoi(string(b))
	if err != nil {
		panic(err)
	}

	return n
}
