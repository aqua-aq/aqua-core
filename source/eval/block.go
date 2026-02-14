package eval

import (
	"github.com/aqua-aq/aqua-core/pkg/scope"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/object/signal"
	"github.com/aqua-aq/aqua-core/source/vm"
)

func RunBlock(b BlockExpression, vm *vm.VM[*object.Value], scope scope.Scope[*object.Value], clone bool, export map[string]*object.Value) object.ExpressionResult {
	scope = scope.Push()
	var res object.ExpressionResult
	for _, expr := range b.Expressions {
		res = IntoEval(expr).Eval(vm, scope, true)
		if res.Signal == signal.SignalRaise && b.Catch != nil {
			scope = scope.Rebase()
			scope.Set(b.Catch.Name.Ident, res.SignalVal.Normalize())
			return BlockExpression(b.Catch.Expressions).Eval(vm, scope, true)
		}

		if res.Signal.Has() {
			return res.Clone(clone)
		}

	}
	for k := range export {
		v, ok := scope.Get(k)
		if ok {
			export[k] = v

		}

	}
	return res.Clone(clone)
}
