package eval

import (
	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/lexer"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/parser"
	"github.com/vandi37/aqua/source/power"
	"github.com/vandi37/aqua/source/signal"
	"github.com/vandi37/aqua/source/vm"
)

func Run(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value], expression string, pos pos.Pos, clone bool) object.ExpressionResult {
	lexer := lexer.NewRelative(expression, pos)
	lexer.Init()
	inPos := lexer.Pos

	err := lexer.Tokenize()
	if err != nil {
		return object.ExpressionResult{
			Trace:     stacktrace.New(pos),
			Signal:    signal.SignalRaise,
			SignalVal: object.IntoValue(err),
		}
	}
	p := parser.New(lexer.Tokens, inPos)
	expr, err := p.Expression(power.PowerLowest, false)
	if err != nil {
		return object.ExpressionResult{
			Trace:     stacktrace.New(pos),
			Signal:    signal.SignalRaise,
			SignalVal: object.IntoValue(err),
		}
	}
	return IntoEval(expr).Eval(vm, scope, clone)
}
