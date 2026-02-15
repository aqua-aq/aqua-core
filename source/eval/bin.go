package eval

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/aqua-aq/aqua-core/pkg/errors"
	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/pkg/scope"
	"github.com/aqua-aq/aqua-core/pkg/stacktrace"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/object/signal"
	"github.com/aqua-aq/aqua-core/source/operators"
	"github.com/aqua-aq/aqua-core/source/vm"
)

func toNumber[T int | float64](b bool) T {
	if b {
		return 1
	}
	return 0
}

func RunBin(
	vm *vm.VM[*object.Value], scope scope.Scope[*object.Value], clone bool,
	left, right *object.Value, operator operators.Operator,
	pos pos.Pos,
) object.ExpressionResult {
	if operator == operators.None {
		return object.ExpressionResult{SignalVal: right, Trace: stacktrace.New(pos)}.Clone(clone)
	}
	typeRaise := object.ExpressionResult{
		Signal: signal.SignalRaise,
		SignalVal: &object.Value{
			InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("unsupported type %s and %s for binary operator '%s'", left.Normalize().Normalize().Type(), right.Normalize().Type(), operator.String()),
			},
		}, Trace: stacktrace.New(pos),
	}

	if operator == operators.Bind {
		if method, ok := operator.Method(); ok && AttrExists(right.Normalize(), method) {
			expr := GetAttrMethod(right.Normalize(), method, pos)
			if expr.Signal.Has() {
				return expr.Clone(clone)
			}
			return Call(vm, expr.SignalVal, []*object.Value{left.Normalize()}, clone, pos, nil).Clone(clone)
		}
		switch right.Normalize().InnerValue.(type) {
		case *object.Subroutine, object.Method:
			return Bind(right.Normalize(), left.Normalize(), pos).Clone(clone)
		}
	}
	if operator == operators.In {
		if method, ok := operator.Method(); ok && AttrExists(right.Normalize(), method) {
			expr := GetAttrMethod(right.Normalize(), method, pos)
			if expr.Signal.Has() {
				return expr.Clone(clone)
			}
			return Call(vm, expr.SignalVal, []*object.Value{left.Normalize()}, clone, pos, nil).Clone(clone)
		}
		switch r := right.Normalize().InnerValue.(type) {
		case object.Array:
			return object.ExpressionResult{SignalVal: &object.Value{InnerValue: object.Bool{Value: slices.ContainsFunc(r.Elements, func(e *object.Value) bool {
				return left.Normalize().Equals(e.Normalize())
			})}}, Trace: stacktrace.New(pos)}
		case object.String:
			str, expr := IntoString(vm, left.Normalize(), pos)
			if expr.Signal.Has() {
				return expr.Clone(clone)
			}
			return object.ExpressionResult{SignalVal: &object.Value{InnerValue: object.Bool{Value: strings.Contains(r.Value, str)}}, Trace: stacktrace.New(pos)}
		case object.Object:
			str, expr := IntoString(vm, left.Normalize(), pos)
			if expr.Signal.Has() {
				return expr.Clone(clone)
			}
			_, ok := r.Map[str]
			return object.ExpressionResult{SignalVal: &object.Value{InnerValue: object.Bool{Value: ok}}, Trace: stacktrace.New(pos)}
		}
	}
	if method, ok := operator.Method(); ok && AttrExists(left.Normalize(), method) {
		expr := GetAttrMethod(left.Normalize(), method, pos)
		if expr.Signal.Has() {
			return expr.Clone(clone)
		}
		return Call(vm, expr.SignalVal, []*object.Value{right.Normalize()}, clone, pos, nil).Clone(clone)
	}
	if operator == operators.Equal {
		return object.ExpressionResult{
			SignalVal: &object.Value{
				InnerValue: object.Bool{Value: left.Equals(right)},
			},
			Trace: stacktrace.New(pos),
		}
	}
	if operator == operators.NotEqual {
		return object.ExpressionResult{
			SignalVal: &object.Value{
				InnerValue: object.Bool{Value: !left.Equals(right)},
			},
			Trace: stacktrace.New(pos),
		}
	}
	switch l := left.Normalize().InnerValue.(type) {
	case object.Number:
		if r, ok := right.Normalize().InnerValue.(object.Number); ok {
			switch operator {
			case operators.Plus:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: l.Value + r.Value}}}
			case operators.Minus:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: l.Value - r.Value}}}
			case operators.Multiply:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: l.Value * r.Value}}}
			case operators.Divide:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: l.Value / r.Value}}}
			case operators.Modulo:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: math.Mod(l.Value, r.Value)}}}
			case operators.StrongDivide:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: math.Floor(l.Value / r.Value)}}}
			case operators.And:
				if math.Floor(r.Value) != r.Value || math.Floor(l.Value) != l.Value {
					return typeRaise
				}
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: float64(int(l.Value) & int(r.Value))}}}
			case operators.Or:
				if math.Floor(r.Value) != r.Value || math.Floor(l.Value) != l.Value {
					return typeRaise
				}
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: float64(int(l.Value) | int(r.Value))}}}
			case operators.Xor:
				if math.Floor(r.Value) != r.Value || math.Floor(l.Value) != l.Value {
					return typeRaise
				}
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: float64(int(l.Value) ^ int(r.Value))}}}
			case operators.Shr:
				if math.Floor(r.Value) != r.Value || math.Floor(l.Value) != l.Value {
					return typeRaise
				}
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: float64(int(l.Value) >> int(r.Value))}}}
			case operators.Shl:
				if math.Floor(r.Value) != r.Value || math.Floor(l.Value) != l.Value {
					return typeRaise
				}
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: float64(int(l.Value) << int(l.Value))}}}
			case operators.Greater:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: l.Value > r.Value}}}
			case operators.Less:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: l.Value < r.Value}}}
			case operators.GreaterEqual:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: l.Value >= r.Value}}}
			case operators.LessEqual:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: l.Value <= r.Value}}}
			}
		}
	case object.Bool:
		if r, ok := right.Normalize().InnerValue.(object.Bool); ok {
			switch operator {
			case operators.Plus:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: toNumber[float64](l.Value) + toNumber[float64](r.Value)}}}
			case operators.Minus:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: toNumber[float64](l.Value) - toNumber[float64](r.Value)}}}
			case operators.Divide, operators.StrongDivide:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: toNumber[float64](l.Value) / toNumber[float64](r.Value)}}}
			case operators.Modulo:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: math.Mod(toNumber[float64](l.Value), toNumber[float64](r.Value))}}}
			case operators.And, operators.Multiply:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: l.Value && r.Value}}}
			case operators.Or:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: l.Value || r.Value}}}
			case operators.Xor:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: toNumber[int](l.Value)^toNumber[int](r.Value) != 0}}}
			case operators.Shr:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: float64(toNumber[int](l.Value) << toNumber[int](r.Value))}}}
			case operators.Shl:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Number{Value: float64(toNumber[int](l.Value) >> toNumber[int](r.Value))}}}
			case operators.Greater:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: l.Value && !r.Value}}}
			case operators.Less:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: !l.Value && r.Value}}}
			case operators.GreaterEqual:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: l.Value || !r.Value}}}
			case operators.LessEqual:
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: !l.Value || r.Value}}}
			}
		}
	case object.Array:
		switch operator {
		case operators.Plus:
			return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.Array{Elements: append(l.Elements, right.Normalize())}}}
		case operators.Index:
			if r, ok := right.Normalize().InnerValue.(object.Number); ok && math.Floor(r.Value) == r.Value {
				rInt := int(r.Value)
				if r.Value < 0 || rInt >= len(l.Elements) {
					return object.ExpressionResult{Trace: stacktrace.New(pos)}
				}
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: l.Elements[rInt]}.Clone(clone)
			}
		}
	case object.String:
		switch operator {
		case operators.Plus:
			str, expr := IntoString(vm, right.Normalize(), pos)
			if expr.Signal.Has() {
				return expr.Clone(clone)
			}
			return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.String{Value: l.Value + str}}}
		case operators.Index:
			runes := []rune(l.Value)
			if r, ok := right.Normalize().InnerValue.(object.Number); ok && math.Floor(r.Value) == r.Value {
				rInt := int(r.Value)
				if r.Value < 0 || rInt >= len(runes) {
					return object.ExpressionResult{Trace: stacktrace.New(pos)}
				}
				return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: &object.Value{InnerValue: object.String{Value: string(runes[rInt])}}}.Clone(clone)
			}
		}
	case object.Object:
		if operator == operators.Index {
			str, expr := IntoString(vm, right.Normalize(), pos)
			if expr.Signal.Has() {
				return expr.Clone(clone)
			}
			res, ok := l.Map[str]
			if !ok {
				return object.ExpressionResult{Trace: stacktrace.New(pos),
					Signal: signal.SignalRaise,
					SignalVal: &object.Value{
						InnerValue: object.Error{
							Code:    errors.ValueError,
							Message: fmt.Sprintf("key '%s' not found", str),
						},
					},
				}
			}
			return object.ExpressionResult{Trace: stacktrace.New(pos), SignalVal: res.Normalize()}
		}
	}
	return typeRaise
}
