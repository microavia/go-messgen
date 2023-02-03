package definition

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var ErrInvalidInput = errors.New("invalid input")

type FieldType struct {
	Name      string
	Array     bool
	ArraySize int
}

var parseType = regexp.MustCompile(`^"([^\[]+)(\[(\d*)])?"$`)

func ParseFieldType(input string) (FieldType, error) {
	typeAndSize := parseType.FindStringSubmatch(input)
	if typeAndSize == nil {
		return FieldType{}, fmt.Errorf("%q: %w", input, ErrInvalidInput)
	}

	return FieldType{
		Name:      typeAndSize[1],
		Array:     len(typeAndSize[2]) > 0,
		ArraySize: mustA2I(typeAndSize[3]),
	}, nil
}

func (f *FieldType) UnmarshalJSON(b []byte) error {
	var err error
	*f, err = ParseFieldType(string(b))

	return err
}

func (f *FieldType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", f.String())), nil
}

func (f *FieldType) String() string {
	if f.Array {
		if f.ArraySize == 0 {
			return f.Name + "[]"
		}

		return f.Name + "[" + strconv.Itoa(f.ArraySize) + "]"
	}

	return f.Name
}

func mustA2I(s string) int {
	if len(s) == 0 {
		return 0
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}

	return n
}
