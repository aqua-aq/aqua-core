package eval

import (
	"fmt"

	"github.com/aqua-aq/aqua-core/pkg/errors"
	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/pkg/stacktrace"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/object/signal"
)

func GetAttr(value *object.Value, name string, pos pos.Pos) object.ExpressionResult {
	obj, ok := value.Normalize().InnerValue().(object.Object)
	if !ok {
		return object.ExpressionResult{Signal: signal.SignalRaise, SignalVal: object.New(object.Error{
			Code:    errors.TypeError,
			Message: fmt.Sprintf("expected an object, got %v", value.Type()),
		}),
			Trace: stacktrace.New(pos),
		}
	}

	if _, ok := obj.Map[name]; !ok {
		obj.Map[name] = object.New(object.Null{})

	}
	return object.ExpressionResult{SignalVal: obj.Map[name].Normalize(), Trace: stacktrace.New(pos)}
}

func DeleteAttr(value *object.Value, name string, pos pos.Pos) object.ExpressionResult {
	obj, ok := value.Normalize().InnerValue().(object.Object)
	if !ok {
		return object.ExpressionResult{Signal: signal.SignalRaise, SignalVal: object.New(object.Error{
			Code:    errors.TypeError,
			Message: fmt.Sprintf("expected an object, got %v", value.Type()),
		}),
			Trace: stacktrace.New(pos),
		}
	}
	attr := obj.Map[name].Normalize()
	delete(obj.Map, name)

	return object.ExpressionResult{SignalVal: attr, Trace: stacktrace.New(pos)}
}

func GetAttrMethod(value *object.Value, name string, pos pos.Pos) object.ExpressionResult {
	attr := GetAttr(value, name, pos)
	if attr.Signal.Has() {
		return attr
	}
	return Bind(attr.SignalVal.Normalize(), value, pos)
}

func AttrExists(value *object.Value, name string) bool {
	obj, ok := value.Normalize().InnerValue().(object.Object)
	if !ok {
		return false
	}
	_, ok = obj.Map[name]
	return ok
}
