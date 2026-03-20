package eval

import (
	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/pkg/stacktrace"
	"github.com/aqua-aq/aqua-core/source/keywords"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/vm"
)

func Clone(vm *vm.VM[*object.Value], val *object.Value, pos pos.Pos) object.ExpressionResult {
	if !AttrExists(val.Normalize(), keywords.Clone) {
		return object.ExpressionResult{Trace: stacktrace.New(pos),
			SignalVal: val.Normalize().DeepClone(),
		}
	}
	method := GetAttrMethod(val.Normalize(), keywords.Iter, pos)
	if method.Signal.Has() {
		return method
	}
	return Call(vm, method.SignalVal.Normalize(), nil, false, pos, nil)
}
