package eval

import (
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/keywords"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
	"github.com/vandi37/aqua/source/vm"
)

func IntoBool(vm *vm.VM, val *object.Value) (bool, object.ExpressionResult) {
	if b, ok := val.Normalize().InnerValue.(object.Bool); ok {
		return b.Value, object.ExpressionResult{}
	}
	method := GetAttrMethod(val.Normalize(), object.TypeBool.String())
	if method.Signal.Has() {
		return false, method
	}
	res := Call(vm, method.SignalVal.Normalize(), nil, false)
	if res.Signal.Has() {
		return false, res
	}
	if b, ok := res.SignalVal.Normalize().InnerValue.(object.Bool); ok {
		return b.Value, object.ExpressionResult{}
	}
	return false, object.ExpressionResult{
		Signal: signal.SignalRaise,
		SignalVal: &object.Value{InnerValue: object.Error{
			Code:    errors.TypeError,
			Message: "can't convert value to boolean",
		}},
	}
}

func IntoIter(val *object.Value, vm *vm.VM) object.ExpressionResult {
	if !AttrExists(val.Normalize(), keywords.Iter) {
		return object.ExpressionResult{
			SignalVal: val.Normalize(),
		}
	}
	method := GetAttrMethod(val.Normalize(), keywords.Iter)
	if method.Signal.Has() {
		return method
	}
	return Call(vm, method.SignalVal.Normalize(), nil, false)
}
