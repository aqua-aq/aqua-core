package eval

import (
	"fmt"

	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/keywords"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
	"github.com/vandi37/aqua/source/vm"
)

func Call(vm *vm.VM, sub *object.Value, args []*object.Value, clone bool) object.ExpressionResult {
	method, ok := sub.Normalize().InnerValue.(object.Method)
	if !ok {
		subroutine, ok := sub.Normalize().InnerValue.(*object.Subroutine)
		if !ok {
			return object.ExpressionResult{
				Signal: signal.SignalRaise,
				SignalVal: &object.Value{InnerValue: object.Error{
					Code:    errors.TypeError,
					Message: fmt.Sprintf("expected subroutine, got %v", sub.Normalize().Type()),
				}},
			}
		}
		method = object.Method{
			Subroutine: subroutine,
			It:         subroutine.Prototype.Clone(),
		}
	}

	subScope := method.Subroutine.Scope.Push()
	subScope.Set(keywords.It, method.It.Normalize())
	object.ParseArgs(method.Subroutine.Arguments, args, subScope)
	if method.Subroutine.BuildIn != nil {
		return method.Subroutine.BuildIn(subScope).AsExpressionResult()
	}
	res := BlockExpression(method.Subroutine.Code).Eval(vm, subScope, clone)
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
	if subRes.SignalVal.Normalize().IsNull() {
		return object.ExpressionResult{SignalVal: method.It}
	}
	return subRes.AsExpressionResult()
}

func Bind(sub, it *object.Value) object.ExpressionResult {
	switch v := sub.Normalize().InnerValue.(type) {
	case object.Method:
		return object.ExpressionResult{
			SignalVal: &object.Value{InnerValue: object.Method{
				Subroutine: v.Subroutine,
				It:         it,
			}},
		}
	case *object.Subroutine:
		return object.ExpressionResult{
			SignalVal: &object.Value{InnerValue: object.Method{
				Subroutine: v,
				It:         it,
			}},
		}
	default:
		return object.ExpressionResult{
			Signal: signal.SignalRaise,
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected subroutine, got %v", sub.Type()),
			}},
		}
	}
}
