package eval

import (
	"fmt"

	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
)

func GetAttr(value *object.Value, name string) object.ExpressionResult {
	obj, ok := value.Normalize().InnerValue.(object.Object)
	if !ok {
		return object.ExpressionResult{Signal: signal.SignalRaise, SignalVal: &object.Value{
			InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected an object, got %v", value.Type()),
			}},
		}
	}
	if attr, ok := obj.Map[name]; ok {
		return object.ExpressionResult{SignalVal: attr.Normalize()}
	}
	return object.ExpressionResult{Signal: signal.SignalRaise, SignalVal: &object.Value{
		InnerValue: object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("attribute %v not found", name),
		}},
	}
}

func GetAttrMethod(value *object.Value, name string) object.ExpressionResult {
	attr := GetAttr(value, name)
	if attr.Signal.Has() {
		return attr
	}
	return Bind(attr.SignalVal.Normalize(), value)
}

func AttrExists(value *object.Value, name string) bool {
	obj, ok := value.Normalize().InnerValue.(object.Object)
	if !ok {
		return false
	}
	_, ok = obj.Map[name]
	return ok
}
