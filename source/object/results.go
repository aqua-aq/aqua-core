package object

import (
	"fmt"

	"github.com/vandi37/aqua/source/signal"
)

type ExpressionResult struct {
	Signal    signal.Signal
	SignalVal *Value
}

type SubroutineResult struct {
	Signal    signal.SubroutineSignal
	SignalVal *Value
}

func (s SubroutineResult) AsExpressionResult() ExpressionResult {
	if s.Signal {
		return ExpressionResult{
			Signal:    signal.SignalRaise,
			SignalVal: s.SignalVal.Normalize(),
		}
	}
	return ExpressionResult{
		SignalVal: s.SignalVal.Normalize(),
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
	}, true

}
func (s ExpressionResult) Clone(need bool) ExpressionResult {
	if !need {
		return s
	}
	return ExpressionResult{
		Signal:    s.Signal,
		SignalVal: s.SignalVal.Normalize().Clone(),
	}
}

func (s ExpressionResult) String() string {
	if s.Signal.Has() {
		return fmt.Sprintf("%v: %v", s.Signal, s.SignalVal.Normalize())
	}
	return fmt.Sprintf("%v", s.SignalVal.Normalize())
}

func (s SubroutineResult) String() string {
	return fmt.Sprintf("%v: %v", s.Signal, s.SignalVal.Normalize())
}
