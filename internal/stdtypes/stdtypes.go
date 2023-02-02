//nolint:gochecknoglobals,gomnd
package stdtypes

var StdTypes = map[string]StdType{
	"char":    {MinSize: 1, Align: 1},
	"int8":    {MinSize: 1, Align: 1},
	"uint8":   {MinSize: 1, Align: 1},
	"int16":   {MinSize: 2, Align: 2},
	"uint16":  {MinSize: 2, Align: 2},
	"int32":   {MinSize: 4, Align: 4},
	"uint32":  {MinSize: 4, Align: 4},
	"int64":   {MinSize: 8, Align: 8},
	"uint64":  {MinSize: 8, Align: 8},
	"float":   {MinSize: 4, Align: 4},
	"float32": {MinSize: 4, Align: 4},
	"float64": {MinSize: 8, Align: 8},
	"double":  {MinSize: 8, Align: 8},
	"string":  {MinSize: 4 + 1, Align: 1, DynamicSize: true},
}

var PlainTypes = map[string]StdType{
	"char":    StdTypes["char"],
	"int8":    StdTypes["int8"],
	"uint8":   StdTypes["uint8"],
	"int16":   StdTypes["int16"],
	"uint16":  StdTypes["uint16"],
	"int32":   StdTypes["int32"],
	"uint32":  StdTypes["uint32"],
	"int64":   StdTypes["int64"],
	"uint64":  StdTypes["uint64"],
	"float":   StdTypes["float"],
	"float32": StdTypes["float32"],
	"float64": StdTypes["float64"],
	"double":  StdTypes["double"],
}

type StdType struct {
	MinSize     int
	Align       int
	DynamicSize bool
}
