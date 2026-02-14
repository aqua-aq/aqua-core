package app

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/vandi37/aqua/pkg/fatal"
	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/source/eval"
	"github.com/vandi37/aqua/source/eval/global"
	"github.com/vandi37/aqua/source/lexer"
	"github.com/vandi37/aqua/source/lexer/tokens"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/object/signal"
	"github.com/vandi37/aqua/source/parser"
	"github.com/vandi37/aqua/source/power"
	"github.com/vandi37/aqua/source/vm"
)

func RunRepl(r io.Reader, vm *vm.VM[*object.Value], pos pos.Pos, quiet bool) error {
	reader := bufio.NewReader(r)
	var current strings.Builder
	lexer := lexer.NewWithPos("", pos)
	lexer.Init()
	scope := scope.New[*object.Value]()
	global.GenerateBuildIn(scope)
	fmt.Print("\033[35m-> \033[0m")
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF && len(line) <= 0 {
			return nil
		}
		if err != nil && err != io.EOF {
			return err
		}

		line = strings.TrimRight(line, "\r\n")

		if strings.HasSuffix(line, "\\") {
			current.WriteString(line[:len(line)-1] + "\n")
		} else {
			current.WriteString(line + "\n")
			lexer.Reload(current.String(), pos)
			err = lexer.Tokenize()
			if err != nil {
				fatal.Fatal(err)
				current.Reset()
				fmt.Print("\033[35m-> \033[0m")
				continue
			}
			parser := parser.New(lexer.Tokens, pos)
			expr, err := parser.Expression(power.PowerLowest, false)
			if err != nil {
				fatal.Fatal(err)
				current.Reset()
				fmt.Print("\033[35m-> \033[0m")
				continue
			}
			_, err = parser.Expect(tokens.TokenEof)
			if err != nil {
				fatal.Fatal(err)
				current.Reset()
				fmt.Print("\033[35m-> \033[0m")
				continue
			}
			res := eval.IntoEval(expr).Eval(vm, scope, false)
			if res.Signal == signal.SignalBreak {
				return nil
			}
			if res.Signal.Has() {
				fatal.Fatal(res.String())
			} else if !quiet && !res.SignalVal.Normalize().IsNull() {
				fmt.Println("\033[34m" + res.SignalVal.Normalize().String() + "\033[0m")
			}
			current.Reset()
			fmt.Print("\033[35m-> \033[0m")
		}
	}
}
