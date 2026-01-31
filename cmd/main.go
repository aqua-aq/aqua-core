package main

import (
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/source/eval"
	"github.com/vandi37/aqua/source/lexer"
	"github.com/vandi37/aqua/source/lexer/tokens"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/parser"
	"github.com/vandi37/aqua/source/vm"
)

func main() {
	path := "main.aq"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	lexer, err := lexer.New(string(data), path)
	if err != nil {
		panic(err)
	}
	lexer.Init()

	err = lexer.Tokenize()
	if err != nil {
		panic(err)
	}
	// for _, tok := range lexer.Tokens {
	// 	fmt.Println(tok)
	// }
	parser, err := parser.New(lexer.Tokens, path)
	if err != nil {
		panic(err)
	}
	_, block, err := parser.ParseBlockExpression(map[tokens.TokenType]struct{}{
		tokens.TokenEof: {},
	})
	if err != nil {
		panic(err)
	}
	for _, v := range block.Expressions {
		spew.Dump(v)
	}
	vm := vm.New()
	s := scope.New[*object.Value]()
	s.Set("print", &object.Value{InnerValue: &object.Subroutine{
		Name:      "print",
		Arguments: object.Arguments{Elements: []object.Argument{{Name: "v"}}},
		Scope:     scope.New[*object.Value](),
		BuildIn: func(s scope.Scope[*object.Value]) object.SubroutineResult {
			print, _ := s.Get("v")
			fmt.Println(print)
			return object.SubroutineResult{}
		},
	}})

	res := eval.BlockExpression(block).Eval(vm, s, false)
	fmt.Println(res)
}
