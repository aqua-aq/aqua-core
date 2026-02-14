package object

import (
	"fmt"

	"github.com/aqua-aq/aqua/pkg/pos"
	"github.com/aqua-aq/aqua/pkg/stacktrace"
	"github.com/aqua-aq/aqua/source/errors"
	"github.com/aqua-aq/aqua/source/object/signal"
)

type ExpressionResult struct {
	Signal    signal.Signal
	SignalVal *Value
	Trace     stacktrace.StackTrace
}

type SubroutineResult struct {
	Signal    signal.SubroutineSignal
	SignalVal *Value
	Trace     stacktrace.StackTrace
}

func (s SubroutineResult) AsExpressionResult() ExpressionResult {
	if s.Signal {
		return ExpressionResult{
			Signal:    signal.SignalRaise,
			SignalVal: s.SignalVal.Normalize(),
			Trace:     s.Trace,
		}
	}
	return ExpressionResult{
		SignalVal: s.SignalVal.Normalize(),
		Trace:     s.Trace,
	}
}

func (s ExpressionResult) IntoSubroutineResult() (SubroutineResult, bool) {
	sig, ok := s.Signal.IntoSubroutineSignal()
	if !ok {
		return SubroutineResult{}, false
	}
	return SubroutineResult{
		Signal:    sig,
		SignalVal: s.SignalVal.Normalize(),
		Trace:     s.Trace,
	}, true

}

func (s ExpressionResult) IntoSubroutineResultStrict(pos pos.Pos) SubroutineResult {
	res, ok := s.IntoSubroutineResult()
	if !ok {
		return SubroutineResult{Trace: stacktrace.New(pos),
			Signal: signal.SubroutineSignalRaise,
			SignalVal: &Value{InnerValue: Error{
				Code:    errors.InvalidSignal,
				Message: fmt.Sprintf("expected none/return/raise, got %v", res.Signal),
			}},
		}
	}
	return res
}
func (s ExpressionResult) Clone(need bool) ExpressionResult {
	if !need {
		return s
	}
	return ExpressionResult{
		Signal:    s.Signal,
		SignalVal: s.SignalVal.Normalize().Clone(),
		Trace:     s.Trace,
	}
}

func (s ExpressionResult) String() string {
	if s.Signal.Has() {
		return fmt.Sprintf("%v: %s\n%s", s.Signal, s.SignalVal.Normalize(), s.Trace.String())
	}
	return fmt.Sprintf("%v", s.SignalVal.Normalize())
}

func (s SubroutineResult) String() string {
	return fmt.Sprintf("%v: %v\n%s", s.Signal, s.SignalVal.Normalize(), s.Trace.String())
}

func (s ExpressionResult) Error() string {
	return s.String()
}

func (s SubroutineResult) Error() string {
	return s.String()
}
