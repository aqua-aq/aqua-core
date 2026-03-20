package object

import (
	"github.com/aqua-aq/aqua-core/pkg/scope"
)

func ParseArgs(args Arguments, vals []*Value, scope scope.Scope[string,*Value]) {
	for i, arg := range args.Elements {
		var val *Value
		if i >= len(vals) {
			val = arg.Default.Normalize()
		} else {
			val = vals[i]
		}
		scope.Set(arg.Name, val)
	}

	if args.Last != nil {
		offset := len(args.Elements)

		var rest []*Value
		if len(vals) > offset {
			rest = vals[offset:]
		} else {
			rest = []*Value{}
		}

		scope.Set(*args.Last, New(Array{Elements: rest}))
	}
}
