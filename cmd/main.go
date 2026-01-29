package main

import (
	"fmt"
	"os"

	"github.com/vandi37/aqua/source/lexer"
)

func main() {
	data, err := os.ReadFile("main.aq")
	if err != nil {
		panic(err)
	}
	lexer, err := lexer.New(string(data), "main.aq")
	if err != nil {
		panic(err)
	}
	lexer.Init()

	err = lexer.Tokenize()
	if err != nil {
		panic(err)
	}

	for _, token := range lexer.Tokens {
		fmt.Println(token)
	}

}
