package eval

import (
	"fmt"

	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
)

func GetAttr(value *object.Value, name string, clone bool) object.ExpressionResult {
	obj, ok := value.Normalize().InnerValue.(object.Object)
	if !ok {
		return object.ExpressionResult{Signal: signal.SignalRaise, SignalVal: &object.Value{
			InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected an object, got %v", value.Type()),
			}},
		}
	}
	if attr, ok := obj.Map[name]; clone && ok {
		return object.ExpressionResult{Value: attr.Normalize().Clone()}
	} else if ok {
		return object.ExpressionResult{Value: attr.Normalize()}
	}
	return object.ExpressionResult{Signal: signal.SignalRaise, SignalVal: &object.Value{
		InnerValue: object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("attribute %v not found", name),
		}},
	}
}

func GetAttrMethod(value *object.Value, name string) object.ExpressionResult {
	attr := GetAttr(value, name, false)
	if attr.Signal.Has() {
		return attr
	}
	return Bind(attr.Value.Normalize(), value)
}
