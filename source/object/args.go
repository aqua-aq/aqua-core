package object

import (
	"github.com/aqua-aq/aqua-core/pkg/scope"
)

func ParseArgs(args Arguments, vals []*Value, scope scope.Scope[*Value]) {
	for i, arg := range args.Elements {
		var val *Value
		if i >= len(vals) {
			val = arg.Default.Normalize()
		} else {
			val = vals[i]
		}
		scope.Set(arg.Name, val)
	}

	if args.Last != nil && len(vals) > len(args.Elements) {
		l := len(args.Elements) - 1
		if l < 0 {
			l = 0
		}
		left := vals[l:]
		scope.Set(*args.Last, &Value{Array{Elements: left}})
	}
}
