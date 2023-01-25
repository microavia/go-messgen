//nolint:gochecknoglobals,gomnd
package stdtypes

var (
	intStringer  = `strconv.FormatInt(int64(v), 10)`
	uintStringer = `strconv.FormatUint(uint64(v), 10)`
)

var StdTypes = map[string]StdType{
	"char":    {Size: 1, Align: 1, BaseType: "byte", Stringer: uintStringer},
	"int8":    {Size: 1, Align: 1, Stringer: intStringer},
	"uint8":   {Size: 1, Align: 1, Stringer: uintStringer},
	"int16":   {Size: 2, Align: 2, Stringer: intStringer},
	"uint16":  {Size: 2, Align: 2, Stringer: uintStringer},
	"int32":   {Size: 4, Align: 4, Stringer: intStringer},
	"uint32":  {Size: 4, Align: 4, Stringer: uintStringer},
	"int64":   {Size: 8, Align: 8, Stringer: intStringer},
	"uint64":  {Size: 8, Align: 8, Stringer: uintStringer},
	"float":   {Size: 4, Align: 4, BaseType: "float32"},
	"float32": {Size: 4, Align: 4},
	"float64": {Size: 8, Align: 8},
	"double":  {Size: 8, Align: 8, BaseType: "float64"},
	"string":  {Size: 1, Align: 1, DynamicSize: true, Stringer: `fmt.Sprintf("%q", v)`},
}

type StdType struct {
	Size        int
	Align       int
	BaseType    string
	Stringer    string
	DynamicSize bool
}
