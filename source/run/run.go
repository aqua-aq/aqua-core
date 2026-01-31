package run

// func Run(path string) (map[string]*object.Value, error) {
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	lexer, err := lexer.New(string(data), path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	lexer.Init()

// 	err = lexer.Tokenize()
// 	if err != nil {
// 		return nil, err
// 	}
// 	parser, err := parser.New(lexer.Tokens, path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	ending, block, err := parser.ParseBlockExpression(map[tokens.TokenType]struct{}{
// 		tokens.TokenEof:    {},
// 		tokens.TokenExport: {},
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	vm := vm.New()
// 	s := scope.New[*object.Value]()
// 	// Create build-in
// 	s.Set("print", &object.Value{InnerValue: &object.Subroutine{
// 		Name:      "print",
// 		Arguments: object.Arguments{Elements: []object.Argument{{Name: "v"}}},
// 		Scope:     scope.New[*object.Value](),
// 		BuildIn: func(s scope.Scope[*object.Value]) object.SubroutineResult {
// 			print, _ := s.Get("v")
// 			fmt.Println(print)
// 			return object.SubroutineResult{}
// 		},
// 	}})
// 	res := eval.ModExpression(module).Eval(vm, scope.Scope[*object.Value]{}, false)
// }
