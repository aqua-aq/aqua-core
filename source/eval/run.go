package eval

import (
	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/pkg/scope"
	"github.com/aqua-aq/aqua-core/pkg/stacktrace"
	"github.com/aqua-aq/aqua-core/source/lexer"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/parser"
	"github.com/aqua-aq/aqua-core/source/power"
	"github.com/aqua-aq/aqua-core/source/object/signal"
	"github.com/aqua-aq/aqua-core/source/vm"
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
