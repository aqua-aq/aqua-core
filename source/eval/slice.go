package eval

import (
	"fmt"
	"math"

	"github.com/aqua-aq/aqua-core/pkg/errors"
	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/pkg/stacktrace"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/object/signal"
)

func ParseSliceIndex(val *object.Value, length int, pos pos.Pos) (int, object.ExpressionResult) {
	num, ok := val.Normalize().InnerValue().(object.Number)
	if !ok || math.Floor(num.Value) != num.Value {
		return 0, object.ExpressionResult{
			Trace:  stacktrace.New(pos),
			Signal: signal.SignalRaise,
			SignalVal: object.New( object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected integer in slice index, got %s", val.Normalize().Type()),
			}),
		}
	}
	idx := int(num.Value)
	if idx > length {
		idx = length
	}
	if -idx > length {
		idx = 0
	}
	if idx < 0 {
		idx = length + idx
	}
	return idx, object.ExpressionResult{Trace: stacktrace.New(pos)}
}
