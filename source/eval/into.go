package eval

import (
	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/keywords"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
	"github.com/vandi37/aqua/source/vm"
)

func IntoBool(vm *vm.VM, val *object.Value, pos pos.Pos) (bool, object.ExpressionResult) {
	if b, ok := val.Normalize().InnerValue.(object.Bool); ok {
		return b.Value, object.ExpressionResult{Trace: stacktrace.New(pos)}
	}
	method := GetAttrMethod(val.Normalize(), keywords.Bool, pos)
	if method.Signal.Has() {
		return false, method
	}
	res := Call(vm, method.SignalVal.Normalize(), nil, false, pos)
	if res.Signal.Has() {
		return false, res
	}
	if b, ok := res.SignalVal.Normalize().InnerValue.(object.Bool); ok {
		return b.Value, object.ExpressionResult{Trace: stacktrace.New(pos)}
	}
	return false, object.ExpressionResult{
		Trace:  stacktrace.New(pos),
		Signal: signal.SignalRaise,
		SignalVal: &object.Value{InnerValue: object.Error{
			Code:    errors.TypeError,
			Message: "can't convert value to boolean",
		}},
	}
}

func IntoIter(val *object.Value, vm *vm.VM, pos pos.Pos) object.ExpressionResult {
	if !AttrExists(val.Normalize(), keywords.Iter) {
		return object.ExpressionResult{Trace: stacktrace.New(pos),
			SignalVal: val.Normalize(),
		}
	}
	method := GetAttrMethod(val.Normalize(), keywords.Iter, pos)
	if method.Signal.Has() {
		return method
	}
	return Call(vm, method.SignalVal.Normalize(), nil, false, pos)
}
func IntoString(vm *vm.VM, val *object.Value, pos pos.Pos) (string, object.ExpressionResult) {
	if !AttrExists(val.Normalize(), keywords.Display) {
		return val.Normalize().String(), object.ExpressionResult{Trace: stacktrace.New(pos)}
	}
	method := GetAttrMethod(val.Normalize(), keywords.Display, pos)
	if method.Signal.Has() {
		return "", method
	}
	res := Call(vm, method.SignalVal.Normalize(), nil, false, pos)
	if res.Signal.Has() {
		return "", res
	}
	return res.SignalVal.Normalize().String(), object.ExpressionResult{Trace: stacktrace.New(pos)}
}
