//nolint:gochecknoglobals,gomnd
package stdtypes

var Types = func() map[string]StdType {
	out := make(map[string]StdType, len(StdTypes)+len(SpecialTypes))

	for k, v := range StdTypes {
		out[k] = v
	}

	for k, v := range SpecialTypes {
		out[k] = v
	}

	return out
}()

var StdTypes = map[string]StdType{
	"char":      {Size: 1, Align: 1},
	"int8":      {Size: 1, Align: 1},
	"uint8":     {Size: 1, Align: 1},
	"int16":     {Size: 2, Align: 2},
	"uint16":    {Size: 2, Align: 2},
	"int32":     {Size: 4, Align: 4},
	"uint32":    {Size: 4, Align: 4},
	"int64":     {Size: 8, Align: 8},
	"uint64":    {Size: 8, Align: 8},
	"float":     {Size: 4, Align: 4},
	"float32":   {Size: 4, Align: 4},
	"float64":   {Size: 8, Align: 8},
	"double":    {Size: 8, Align: 8},
	"bitmask8":  {Size: 1, Align: 1},
	"bitmask16": {Size: 2, Align: 2},
	"bitmask32": {Size: 4, Align: 4},
	"bitmask64": {Size: 8, Align: 8},
}

var SpecialTypes = map[string]StdType{
	"string": {Size: 1, Align: 1},
}

type StdType struct {
	Size  int
	Align int
}
