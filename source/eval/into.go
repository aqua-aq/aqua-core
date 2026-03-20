package eval

import (
	"fmt"
	"math"
	"strconv"

	"github.com/aqua-aq/aqua-core/pkg/errors"
	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/pkg/stacktrace"
	"github.com/aqua-aq/aqua-core/source/keywords"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/object/signal"
	"github.com/aqua-aq/aqua-core/source/vm"
)

func IntoBool(vm *vm.VM[*object.Value], val *object.Value, pos pos.Pos) (bool, object.ExpressionResult) {
	switch v := val.Normalize().InnerValue().(type) {
	case object.Bool:
		return v.Value, object.ExpressionResult{Trace: stacktrace.New(pos)}
	case object.Null:
		return false, object.ExpressionResult{Trace: stacktrace.New(pos)}
	case object.Number:
		return v.Value != 0, object.ExpressionResult{Trace: stacktrace.New(pos)}
	case object.String:
		return v.Value != "", object.ExpressionResult{Trace: stacktrace.New(pos)}
	case object.Array:
		return len(v.Elements) > 0, object.ExpressionResult{Trace: stacktrace.New(pos)}
	case object.Object:
		if !AttrExists(val.Normalize(), keywords.Bool) {
			return len(v.Map) > 0, object.ExpressionResult{Trace: stacktrace.New(pos)}
		}
		method := GetAttrMethod(val.Normalize(), keywords.Bool, pos)
		if method.Signal.Has() {
			return false, method
		}
		res := Call(vm, method.SignalVal.Normalize(), nil, false, pos, nil)
		if res.Signal.Has() {
			return false, res
		}
		return IntoBool(vm, res.SignalVal.Normalize(), pos)
	}
	return true, object.ExpressionResult{
		Trace: stacktrace.New(pos),
	}
}

func IntoIter(val *object.Value, vm *vm.VM[*object.Value], pos pos.Pos) object.ExpressionResult {
	if !AttrExists(val.Normalize(), keywords.Iter) {
		return object.ExpressionResult{Trace: stacktrace.New(pos),
			SignalVal: val.Normalize(),
		}
	}
	method := GetAttrMethod(val.Normalize(), keywords.Iter, pos)
	if method.Signal.Has() {
		return method
	}
	return Call(vm, method.SignalVal.Normalize(), nil, false, pos, nil)
}
func IntoString(vm *vm.VM[*object.Value], val *object.Value, pos pos.Pos) (string, object.ExpressionResult) {
	if !AttrExists(val.Normalize(), keywords.Display) {
		return val.Normalize().String(), object.ExpressionResult{Trace: stacktrace.New(pos)}
	}
	method := GetAttrMethod(val.Normalize(), keywords.Display, pos)
	if method.Signal.Has() {
		return "", method
	}
	res := Call(vm, method.SignalVal.Normalize(), nil, false, pos, nil)
	if res.Signal.Has() {
		return "", res
	}
	return IntoString(vm, res.SignalVal.Normalize(), pos)
}

func IntoNum(vm *vm.VM[*object.Value], val *object.Value, pos pos.Pos) (float64, object.ExpressionResult) {
	switch v := val.Normalize().InnerValue().(type) {
	case object.Bool:
		if v.Value {
			return 1, object.ExpressionResult{Trace: stacktrace.New(pos)}
		}
		return 0, object.ExpressionResult{Trace: stacktrace.New(pos)}
	case object.Null:
		return 0, object.ExpressionResult{Trace: stacktrace.New(pos)}
	case object.Number:
		return v.Value, object.ExpressionResult{Trace: stacktrace.New(pos)}
	case object.String:
		n, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			return 0, object.ExpressionResult{Trace: stacktrace.New(pos), Signal: signal.SignalRaise, SignalVal: object.New(object.Error{
				Code:    errors.ValueError,
				Message: err.Error(),
			},
			)}
		}
		return n, object.ExpressionResult{Trace: stacktrace.New(pos)}
	case object.Object:
		if !AttrExists(val.Normalize(), keywords.Number) {
			break
		}
		method := GetAttrMethod(val.Normalize(), keywords.Number, pos)
		if method.Signal.Has() {
			return 0, method
		}
		res := Call(vm, method.SignalVal.Normalize(), nil, false, pos, nil)
		if res.Signal.Has() {
			return 0, res
		}
		return IntoNum(vm, res.SignalVal.Normalize(), pos)
	}
	return 0, object.ExpressionResult{
		Trace:  stacktrace.New(pos),
		Signal: signal.SignalRaise,
		SignalVal: object.New(object.Error{
			Code:    errors.TypeError,
			Message: fmt.Sprintf("can't convert value with type %s to number", val.Normalize().Type()),
		}),
	}
}

func IntoInt(vm *vm.VM[*object.Value], val *object.Value, pos pos.Pos) (int, object.ExpressionResult) {
	f, err := IntoNum(vm, val, pos)
	if err.Signal.Has() {
		return 0, err
	}
	if math.Floor(f) != f {
		return 0, object.ExpressionResult{
			Signal: signal.SignalRaise,
			SignalVal: object.New(object.Error{
					Code:    errors.TypeError,
					Message: "expected an integer, got number",
				},
			), Trace: stacktrace.New(pos),
		}
	}
	return int(f), object.ExpressionResult{Trace: stacktrace.New(pos)}
}
