package eval

import (
	"fmt"
	"math"

	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/object/signal"
)

func ParseSliceIndex(val *object.Value, length int, pos pos.Pos) (int, object.ExpressionResult) {
	num, ok := val.Normalize().InnerValue.(object.Number)
	if !ok || math.Floor(num.Value) != num.Value {
		return 0, object.ExpressionResult{
			Trace:  stacktrace.New(pos),
			Signal: signal.SignalRaise,
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected integer in slice index, got %s", val.Normalize().Type()),
			}},
		}
	}
	idx := int(num.Value)
	if idx < 0 || idx > length {
		return 0, object.ExpressionResult{
			Trace:  stacktrace.New(pos),
			Signal: signal.SignalRaise,
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.ValueError,
				Message: fmt.Sprintf("index %d is out of range, expected [0; %d]", idx, length),
			}},
		}
	}
	return idx, object.ExpressionResult{Trace: stacktrace.New(pos)}
}
