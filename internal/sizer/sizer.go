package sizer

import (
	"errors"
	"fmt"

	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/stdtypes"
)

type TypeSize struct {
	MinSize int
	Dynamic bool
	Align   int
	Plain   bool
}

const lenSize = 4

var ErrInvalidAlign = errors.New("invalid alignment")

func MinSizeByName(name string, def definition.Definition) (TypeSize, error) {
	typeDef, err := definition.ParseFieldType(`"` + name + `"`)
	if err != nil {
		return TypeSize{}, fmt.Errorf("sizing type %q: %w", name, err)
	}

	return MinSize(typeDef, def), nil
}

func MinSize(typeDef definition.FieldType, def definition.Definition) TypeSize {
	if typeDef.Array {
		return minSizeOfArray(typeDef.ArraySize, MinSize(definition.FieldType{Name: typeDef.Name}, def))
	}

	if v, ok := stdtypes.StdTypes[typeDef.Name]; ok {
		return TypeSize{MinSize: v.MinSize, Dynamic: v.DynamicSize, Align: v.Align, Plain: !v.DynamicSize}
	}

	if v, ok := search(def.Enums, func(enum definition.Enum) bool { return enum.Name == typeDef.Name }); ok {
		baseType := stdtypes.StdTypes[v.BaseType]

		return TypeSize{
			MinSize: baseType.MinSize,
			Dynamic: baseType.DynamicSize,
			Align:   baseType.Align,
			Plain:   !baseType.DynamicSize,
		}
	}

	out := TypeSize{MinSize: 0, Dynamic: false}

	if v, ok := search(def.Messages, func(m definition.Message) bool { return m.Name == typeDef.Name }); ok {
		for _, field := range v.Fields {
			fieldSize := MinSize(field.Type, def)
			out.MinSize += fieldSize.MinSize
			out.Dynamic = out.Dynamic || fieldSize.Dynamic
			out.Align = maxInt(out.Align, fieldSize.Align)
		}
	}

	if err := checkAlighment(out.Align); err != nil {
		panic(fmt.Errorf("%+v: %+v: %w", typeDef, out, err))
	}

	return out
}

func checkAlighment(align int) error {
	if align > 8 || align < 1 || (align > 1 && align%2 != 0) {
		return ErrInvalidAlign
	}

	return nil
}

func minSizeOfArray(arraySize int, typeSize TypeSize) TypeSize {
	if arraySize == 0 {
		return TypeSize{
			MinSize: lenSize + typeSize.MinSize,
			Dynamic: true,
			Align:   typeSize.Align,
		}
	}

	return TypeSize{
		MinSize: arraySize * typeSize.MinSize,
		Dynamic: typeSize.Dynamic,
		Align:   typeSize.Align,
	}
}

func search[V any](l []V, f func(V) bool) (V, bool) { //nolint:ireturn
	for _, v := range l {
		if f(v) {
			return v, true
		}
	}

	return *new(V), false //nolint:gocritic
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}
