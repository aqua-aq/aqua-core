package eval

import (
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
	"github.com/vandi37/aqua/source/vm"
)

func RunBlock(b BlockExpression, vm *vm.VM, scope scope.Scope[*object.Value], clone bool, export map[string]*object.Value) object.ExpressionResult {
	scope = scope.Push()
	for i, expr := range b.Expressions {
		res := IntoEval(expr).Eval(vm, scope, true)
		if res.Signal == signal.SignalRaise && b.Catch != nil {
			scope = scope.Rebase()
			scope.Set(b.Catch.Name.Ident, res.SignalVal.Normalize())
			return BlockExpression(b.Catch.Expressions).Eval(vm, scope, true)
		}
		if res.Signal.Has() {
			return res.Clone(clone)
		}
		if i == len(b.Expressions)-1 {
			return res.Clone(clone)
		}
	}
	for k := range export {
		v, ok := scope.Get(k)
		if ok {
			export[k] = v

		}

	}
	return object.ExpressionResult{Trace: stacktrace.New(b.Pos),
		SignalVal: &object.Value{InnerValue: object.Null{}},
	}
}
