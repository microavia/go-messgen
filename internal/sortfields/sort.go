package sortfields

import (
	"sort"

	"github.com/microavia/go-messgen/internal/definition"
	"github.com/microavia/go-messgen/internal/sizer"
)

func SortFields(def definition.Definition) {
	for _, msg := range def.Messages {
		sort.SliceStable(
			msg.Fields,
			func(i, j int) bool {
				return fieldsLess(
					sizer.MinSize(msg.Fields[i].Type, def),
					sizer.MinSize(msg.Fields[j].Type, def),
				)
			},
		)
	}
}

func fieldsLess(f1, f2 sizer.TypeSize) bool { //nolint:cyclop
	switch {
	case f1.Plain && !f2.Plain:
		return true
	case !f1.Plain && f2.Plain:
		return false
	case f1.Plain && f2.Plain:
		return f1.Align > f2.Align
	case !f1.Dynamic && f2.Dynamic:
		return true
	case f1.Dynamic && !f2.Dynamic:
		return false
	case f1.Dynamic && f2.Dynamic:
		return f1.Align > f2.Align
	}

	return f1.Align > f2.Align
}
