package object

import (
	"fmt"

	"github.com/vandi37/aqua/source/signal"
)

type ExpressionResult struct {
	Value     *Value
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
			SignalVal: s.SignalVal,
		}
	}
	return ExpressionResult{
		Value: s.SignalVal,
	}
}

func (s ExpressionResult) IntoSubroutineResult() (SubroutineResult, bool) {
	if !s.Signal.Has() {
		return SubroutineResult{
			SignalVal: s.Value,
		}, true
	}
	sig, ok := s.Signal.IntoSubroutineSignal()
	if !ok {
		return SubroutineResult{}, false
	}
	return SubroutineResult{
		Signal:    sig,
		SignalVal: s.SignalVal,
	}, true

}
func (s ExpressionResult) Clone(need bool) ExpressionResult {
	if !need {
		return s
	}
	if s.Signal.Has() {
		return ExpressionResult{
			Signal:    s.Signal,
			SignalVal: s.SignalVal.Clone(),
		}
	}
	return ExpressionResult{
		Value: s.Value.Clone(),
	}
}

func (s ExpressionResult) String() string {
	if s.Signal.Has() {
		return fmt.Sprintf("%v: %v", s.Signal, s.SignalVal)
	}
	return fmt.Sprintf("%v", s.Value)
}

func (s SubroutineResult) String() string {
	return fmt.Sprintf("%v: %v", s.Signal, s.SignalVal)
}
