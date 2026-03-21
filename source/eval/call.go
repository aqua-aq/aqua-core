package eval

import (
	"fmt"

	"github.com/aqua-aq/aqua-core/pkg/errors"
	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/pkg/scope"
	"github.com/aqua-aq/aqua-core/pkg/stacktrace"
	"github.com/aqua-aq/aqua-core/source/keywords"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/object/signal"
	"github.com/aqua-aq/aqua-core/source/vm"
)

func Call(vm *vm.VM[*object.Value], sub *object.Value, args []*object.Value, clone bool, pos pos.Pos, export map[string]*object.Value, new bool) object.ExpressionResult {
	if AttrExists(sub, keywords.Call) {
		method := GetAttrMethod(sub, keywords.Call, pos)
		if method.Signal.Has() {
			return method
		}
		sub = method.SignalVal.Normalize()
	}
	method, ok := sub.Normalize().InnerValue().(object.Method)
	if !ok {
		subroutine, ok := sub.Normalize().InnerValue().(*object.Subroutine)
		if !ok {
			return object.ExpressionResult{Trace: stacktrace.New(pos),
				Signal: signal.SignalRaise,
				SignalVal: object.New(object.Error{
					Code:    errors.TypeError,
					Message: fmt.Sprintf("expected subroutine, got %v", sub.Normalize().Type()),
				}),
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
	var res object.SubroutineResult
	if method.Subroutine.BuiltIn != nil {
		res = method.Subroutine.BuiltIn(vm, subScope)
	} else {
		expr := RunBlock(BlockExpression(method.Subroutine.Code), vm, subScope, clone, export)
		res, ok = expr.IntoSubroutineResult()
		if !ok {
			return object.ExpressionResult{Trace: stacktrace.New(pos),
				Signal: signal.SignalRaise,
				SignalVal: object.New(object.Error{
					Code:    errors.InvalidSignal,
					Message: fmt.Sprintf("expected none/return/raise, got %v", res.Signal),
				}),
			}
		}
	}
	res.Trace = res.Trace.Add(method.Subroutine.Name, pos)
	if res.SignalVal.Normalize().IsNull() && new && !bool(res.Signal) {
		return object.ExpressionResult{Signal: res.Signal.IntoSignal(), SignalVal: method.It, Trace: res.Trace}
	}
	return res.AsExpressionResult()
}

func Bind(sub, it *object.Value, pos pos.Pos) object.ExpressionResult {
	switch v := sub.Normalize().InnerValue().(type) {
	case object.Method:
		return object.ExpressionResult{Trace: stacktrace.New(pos),
			SignalVal: object.New(object.Method{
				Subroutine: v.Subroutine,
				It:         it,
			}),
		}
	case *object.Subroutine:
		return object.ExpressionResult{Trace: stacktrace.New(pos),
			SignalVal: object.New(object.Method{
				Subroutine: v,
				It:         it,
			}),
		}
	default:
		return object.ExpressionResult{Trace: stacktrace.New(pos),
			Signal: signal.SignalRaise,
			SignalVal: object.New(object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected subroutine, got %v", sub.Type()),
			}),
		}
	}
}

func (c CallExpression) EvalCanBeNew(vm *vm.VM[*object.Value], scope scope.Scope[string, *object.Value], clone, new bool) object.ExpressionResult {
	sub := IntoEval(c.Subroutine).Eval(vm, scope, false)
	if sub.Signal.Has() {
		return sub
	}
	args := make([]*object.Value, 0, len(c.Args))
	for _, val := range c.Args {
		res := IntoEval(val.Value).Eval(vm, scope, true)
		if res.Signal.Has() {
			return Clone(clone, vm, res, c.Pos)
		}
		if val.IsContinuos {

			inner, ok := res.SignalVal.InnerValue().(object.Array)
			if !ok {
				return object.ExpressionResult{Trace: stacktrace.New(val.Pos),
					Signal: signal.SignalRaise,
					SignalVal: object.New(object.Error{
						Code:    errors.TypeError,
						Message: fmt.Sprintf("expected array, got %v", res.SignalVal.Normalize().Type()),
					}),
				}
			}
			args = append(args, inner.Elements...)
		} else {
			args = append(args, res.SignalVal.Normalize())

		}
	}
	return Call(vm, sub.SignalVal, args, clone, c.Pos, nil, new)
}
