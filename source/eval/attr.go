package eval

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
)

func GetAttr(value *object.Value, name string, pos pos.Pos) object.ExpressionResult {
	obj, ok := value.Normalize().InnerValue.(object.Object)
	if !ok {
		return object.ExpressionResult{Signal: signal.SignalRaise, SignalVal: &object.Value{
			InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected an object, got %v", value.Type()),
			}},
			Trace: stacktrace.New(pos),
		}
	}

	if _, ok := obj.Map[name]; !ok {
		obj.Map[name] = &object.Value{InnerValue: object.Null{}}

	}
	return object.ExpressionResult{SignalVal: obj.Map[name].Normalize(), Trace: stacktrace.New(pos)}
}

func GetAttrMethod(value *object.Value, name string, pos pos.Pos) object.ExpressionResult {
	attr := GetAttr(value, name, pos)
	if attr.Signal.Has() {
		return attr
	}
	return Bind(attr.SignalVal.Normalize(), value, pos)
}

func AttrExists(value *object.Value, name string) bool {
	obj, ok := value.Normalize().InnerValue.(object.Object)
	if !ok {
		return false
	}
	_, ok = obj.Map[name]
	return ok
}
