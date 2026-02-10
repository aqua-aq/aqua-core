package eval

import (
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/ast"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/vm"
)

func DeclareSubroutine(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value], clone bool, name string, s ast.SubroutineDec) object.ExpressionResult {
	arguments := object.Arguments{Last: s.Arguments.Last}
	for _, arg := range s.Arguments.Elements {
		res := IntoEval(arg.Default).Eval(vm, scope, true)
		if res.Signal.Has() {
			return res.Clone(clone)
		}
		arguments.Elements = append(arguments.Elements, object.Argument{
			Name:    arg.Name.Ident,
			Default: res.SignalVal.Normalize(),
		})
	}
	res := IntoEval(s.Prototype).Eval(vm, scope, true)
	if res.Signal.Has() {
		return res.Clone(clone)
	}
	return object.ExpressionResult{Trace: stacktrace.New(s.Pos),
		SignalVal: &object.Value{
			InnerValue: &object.Subroutine{
				Arguments: arguments,
				Scope:     scope,
				Prototype: res.SignalVal.Normalize(),
				Code:      s.Body,
				Name:      name,
			},
		},
	}
}
