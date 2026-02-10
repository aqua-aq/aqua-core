package eval

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/keywords"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
	"github.com/vandi37/aqua/source/vm"
)

func Call(vm *vm.VM[*object.Value], sub *object.Value, args []*object.Value, clone bool, pos pos.Pos, export map[string]*object.Value) object.ExpressionResult {
	method, ok := sub.Normalize().InnerValue.(object.Method)
	if !ok {
		subroutine, ok := sub.Normalize().InnerValue.(*object.Subroutine)
		if !ok {
			return object.ExpressionResult{Trace: stacktrace.New(pos),
				Signal: signal.SignalRaise,
				SignalVal: &object.Value{InnerValue: object.Error{
					Code:    errors.TypeError,
					Message: fmt.Sprintf("expected subroutine, got %v", sub.Normalize().Type()),
				}},
			}
		}
		method = object.Method{
			Subroutine: subroutine,
			It:         subroutine.Prototype.Normalize().Clone(),
		}
	}

	subScope := method.Subroutine.Scope.Push()
	subScope.Set(keywords.It, method.It.Normalize())
	object.ParseArgs(method.Subroutine.Arguments, args, subScope)
	if method.Subroutine.BuildIn != nil {
		res := method.Subroutine.BuildIn(vm, subScope).AsExpressionResult()
		res.Trace = res.Trace.Add(method.Subroutine.Name, pos)
		return res
	}
	res := RunBlock(BlockExpression(method.Subroutine.Code), vm, subScope, clone, export)
	res.Trace = res.Trace.Add(method.Subroutine.Name, pos)
	subRes, ok := res.IntoSubroutineResult()
	if !ok {
		return object.ExpressionResult{Trace: stacktrace.New(pos),
			Signal: signal.SignalRaise,
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.InvalidSignal,
				Message: fmt.Sprintf("expected none/return/raise, got %v", res.Signal),
			}},
		}
	}
	if subRes.SignalVal.Normalize().IsNull() {
		return object.ExpressionResult{Signal: subRes.Signal.IntoSignal(), SignalVal: method.It, Trace: subRes.Trace}
	}
	return subRes.AsExpressionResult()
}

func Bind(sub, it *object.Value, pos pos.Pos) object.ExpressionResult {
	switch v := sub.Normalize().InnerValue.(type) {
	case object.Method:
		return object.ExpressionResult{Trace: stacktrace.New(pos),
			SignalVal: &object.Value{InnerValue: object.Method{
				Subroutine: v.Subroutine,
				It:         it,
			}},
		}
	case *object.Subroutine:
		return object.ExpressionResult{Trace: stacktrace.New(pos),
			SignalVal: &object.Value{InnerValue: object.Method{
				Subroutine: v,
				It:         it,
			}},
		}
	default:
		return object.ExpressionResult{Trace: stacktrace.New(pos),
			Signal: signal.SignalRaise,
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected subroutine, got %v", sub.Type()),
			}},
		}
	}
}
