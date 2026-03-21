package eval

import (
	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/source/keywords"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/vm"
)

func Clone(clone bool, vm *vm.VM[*object.Value], expr object.ExpressionResult, pos pos.Pos) object.ExpressionResult {
	if !clone {
		return expr
	}
	if !AttrExists(expr.SignalVal.Normalize(), keywords.Clone) {
		return object.ExpressionResult{
			SignalVal: expr.SignalVal.Clone(),
			Signal:    expr.Signal,
			Trace:     expr.Trace,
		}
	}
	method := GetAttrMethod(expr.SignalVal.Normalize(), keywords.Clone, pos)
	if method.Signal.Has() {
		return method
	}
	res := Call(vm, method.SignalVal.Normalize(), nil, false, pos, nil, false)
	if res.Signal.Has() {
		return res
	}
	return object.ExpressionResult{
		Trace:     expr.Trace,
		Signal:    expr.Signal,
		SignalVal: res.SignalVal,
	}
}
