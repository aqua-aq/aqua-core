package eval

import (
	"github.com/aqua-aq/aqua-core/pkg/scope"
	"github.com/aqua-aq/aqua-core/pkg/stacktrace"
	"github.com/aqua-aq/aqua-core/source/ast"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/vm"
)

func DeclareSubroutine(vm *vm.VM[*object.Value], scope scope.Scope[string, *object.Value], clone bool, name string, s ast.SubroutineDec) object.ExpressionResult {
	arguments := object.Arguments{Last: s.Arguments.Last}
	for _, arg := range s.Arguments.Elements {
		res := IntoEval(arg.Default).Eval(vm, scope, true)
		if res.Signal.Has() {
			return Clone(clone, vm, res, s.Pos)
		}
		arguments.Elements = append(arguments.Elements, object.Argument{
			Name:    arg.Name.Ident,
			Default: res.SignalVal.Normalize(),
		})
	}
	res := IntoEval(s.Prototype).Eval(vm, scope, true)
	if res.Signal.Has() {
		return Clone(clone, vm, res, s.Pos)
	}
	return object.ExpressionResult{Trace: stacktrace.New(s.Pos),
		SignalVal: object.New(&object.Subroutine{
			Arguments: arguments,
			Scope:     scope,
			Prototype: res.SignalVal.Normalize(),
			Code:      s.Body,
			Name:      name,
		}),
	}
}
