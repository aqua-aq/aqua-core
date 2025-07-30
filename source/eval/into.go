package eval

import (
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
)

func IntoBool(val *object.Value) (bool, object.ExpressionResult) {
	if b, ok := val.Normalize().InnerValue.(object.Bool); ok {
		return b.Value, object.ExpressionResult{}
	}
	if _, ok := val.Normalize().InnerValue.(object.Object); ok {
		intoBool := GetAttrMethod(val, object.TypeBool.String())
		if intoBool.Signal.Has() {
			return false, intoBool
		}
		if b, ok := intoBool.Value.Normalize().InnerValue.(object.Bool); ok {
			return b.Value, object.ExpressionResult{}
		}
	}

	return false, object.ExpressionResult{
		Signal: signal.SignalRaise,
		SignalVal: &object.Value{InnerValue: object.Error{
			Code:    errors.TypeError,
			Message: "can't convert value to boolean",
		}},
	}
}

func IntoIter()
