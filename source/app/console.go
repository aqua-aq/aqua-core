package app

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/aqua-aq/aqua-core/pkg/fatal"
	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/pkg/scope"
	"github.com/aqua-aq/aqua-core/source/eval"
	"github.com/aqua-aq/aqua-core/source/eval/global"
	"github.com/aqua-aq/aqua-core/source/lexer"
	"github.com/aqua-aq/aqua-core/source/lexer/tokens"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/object/signal"
	"github.com/aqua-aq/aqua-core/source/parser"
	"github.com/aqua-aq/aqua-core/source/power"
	"github.com/aqua-aq/aqua-core/source/vm"
)

func RunRepl(r io.Reader, vm *vm.VM[*object.Value], pos pos.Pos, quiet bool) error {
	reader := bufio.NewReader(r)
	var current strings.Builder
	lexer := lexer.NewWithPos("", pos)
	lexer.Init()
	scope := scope.New[string, *object.Value]()
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
