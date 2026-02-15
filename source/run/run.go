package run

import (
	"fmt"
	"os"

	"github.com/aqua-aq/aqua-core/pkg/scope"
	"github.com/aqua-aq/aqua-core/source/ast"
	"github.com/aqua-aq/aqua-core/source/eval"
	"github.com/aqua-aq/aqua-core/source/eval/global"
	"github.com/aqua-aq/aqua-core/source/lexer"
	"github.com/aqua-aq/aqua-core/source/lexer/tokens"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/parser"
	"github.com/aqua-aq/aqua-core/source/vm"
)

func Run(path, name string, vm *vm.VM[*object.Value]) (map[string]*object.Value, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lexer, err := lexer.New(string(data), path)
	if err != nil {
		return nil, err
	}
	lexer.Init()
	pos := lexer.Pos

	err = lexer.Tokenize()
	if err != nil {
		return nil, err
	}
	parser := parser.New(lexer.Tokens, pos)
	ending, block, err := parser.ParseBlockExpression(map[tokens.TokenType]struct{}{
		tokens.TokenEof:    {},
		tokens.TokenExport: {},
	})
	if err != nil {
		return nil, err
	}
	export := []string{}
	if ending == tokens.TokenExport {
		firstExport, err := parser.Expect(tokens.TokenIdentifier)
		if err != nil {
			return nil, err
		}
		export = append(export, firstExport.Value)
		for peek, ok := parser.Peek(0); ok && peek.Type == tokens.TokenComma; peek, ok = parser.Peek(0) {
			parser.Move()
			next, err := parser.Expect(tokens.TokenIdentifier)
			if err != nil {
				return nil, err
			}
			export = append(export, next.Value)
		}
	}
	scope := scope.New[*object.Value]()
	global.GenerateBuildIn(scope)
	sub := eval.DeclareSubroutine(vm, scope, false, fmt.Sprintf("<%s>", name), ast.SubroutineDec{
		Arguments: ast.Arguments{},
		Body:      block,
		Prototype: ast.NullDec{Pos: pos},
		Pos:       pos,
	})
	if sub.Signal.Has() {
		return nil, sub
	}
	vals := make(map[string]*object.Value, len(export))
	for _, v := range export {
		vals[v] = &object.Value{InnerValue: object.Null{}}
	}
	call := eval.Call(vm, sub.SignalVal.Normalize(), []*object.Value{}, false, pos, vals)
	if call.Signal.Has() {
		return nil, call
	}
	return vals, nil
}
