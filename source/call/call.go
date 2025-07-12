package call

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/eval"
	"github.com/vandi37/aqua/source/keywords"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
	"github.com/vandi37/aqua/source/vm"
)

func Call(vm *vm.VM, _ scope.Scope[*object.Value], clone bool) func(sub, args *object.Value) object.ExpressionResult {
	return func(sub, args *object.Value) object.ExpressionResult {
		var it = &object.Value{InnerValue: nil}
		subroutine, ok := sub.Normalize().InnerValue.(*object.Subroutine)
		if !ok {
			method, ok := sub.InnerValue.(*object.Method)
			if !ok {
				return object.ExpressionResult{
					Signal: signal.SignalRaise,
					SignalVal: &object.Value{InnerValue: object.Error{
						Code:    errors.TypeError,
						Message: fmt.Sprintf("expected subroutine, got %v", sub.Type()),
					}},
				}
			}
			it = method.It
		}
		subScope := subroutine.Scope.Push()
		subScope.Set(keywords.It, it.Normalize())
		vals, ok := args.Normalize().InnerValue.(object.Array)
		if !ok {
			return object.ExpressionResult{
				Signal: signal.SignalRaise,
				SignalVal: &object.Value{InnerValue: object.Error{
					Code:    errors.TypeError,
					Message: fmt.Sprintf("expected array, got %v", sub.Type()),
				}},
			}
		}
		object.ParseArgs(subroutine.Arguments, vals.Elements, &subScope)
		if subroutine.BuildIn != nil {
			return subroutine.BuildIn(subScope).AsExpressionResult()
		}
		res := eval.BlockExpression(subroutine.Code).Eval(vm, subScope, false)
		subRes, ok := res.IntoSubroutineResult()
		if !ok {
			return object.ExpressionResult{
				Signal: signal.SignalRaise,
				SignalVal: &object.Value{InnerValue: object.Error{
					Code:    errors.InvalidSignal,
					Message: fmt.Sprintf("expected none/return/raise, got %v", res.Signal),
				}},
			}
		}
		return subRes.AsExpressionResult()
	}
}
