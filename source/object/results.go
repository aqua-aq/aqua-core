package object

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/signal"
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
