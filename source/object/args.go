package object

import (
	"github.com/vandi37/aqua/pkg/scope"
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
		left := vals[len(args.Elements)-1:]
		scope.Set(*args.Last, &Value{Array{Elements: left}})
	}
}
