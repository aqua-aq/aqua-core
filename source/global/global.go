package global

import (
	"fmt"
	"strings"

	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/eval"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/vm"
)

func AddGlobalFunction(
	name string,
	f func(*vm.VM[*object.Value], scope.Scope[*object.Value]) object.SubroutineResult,
	args object.Arguments,
	scope scope.Scope[*object.Value]) {
	sub := object.Subroutine{
		Name:      name,
		Arguments: args,
		BuildIn:   f,
		Scope:     scope,
	}
	scope.Set(name, &object.Value{
		InnerValue: &sub,
	})
}

var last = "args"
var PrintArgs = object.Arguments{Last: &last}

func Print(ln bool, name string) func(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
	return func(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
		v, ok := scope.Get(last)
		if !ok {
			return Raise(name, object.Error{
				Code:    errors.ValueError,
				Message: fmt.Sprintf("expected argument %s", last),
			})
		}
		args, ok := v.Normalize().InnerValue.(object.Array)
		if !ok {
			return Raise(name, object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected %s, got %s", object.TypeArray, v.Normalize().Type()),
			})
		}
		b := strings.Builder{}
		for i, v := range args.Elements {
			if i > 0 {
				b.WriteRune(' ')
			}
			str, err := eval.IntoString(vm, v, pos.BuildInPos(name))
			if sErr, _ := err.IntoSubroutineResult(); err.Signal.Has() {
				return sErr
			}
			b.WriteString(str)
		}
		if ln {
			b.WriteRune('\n')
		}
		fmt.Print(b.String())
		return object.SubroutineResult{Trace: stacktrace.New(pos.BuildInPos(name))}
	}
}

var firstArgument = "first"
var InputArguments = object.Arguments{Elements: []object.Argument{{
	Name:    firstArgument,
	Default: &object.Value{InnerValue: object.String{Value: ""}},
}}}

func Input(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
	prompt, ok := scope.Get(firstArgument)
	if !ok {
		return Raise("input", object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("expected argument %s", prompt),
		})
	}
	str, err := eval.IntoString(vm, prompt, pos.BuildInPos("input"))
	if sErr, _ := err.IntoSubroutineResult(); err.Signal.Has() {
		return sErr
	}
	fmt.Print(str)
	var s string
	fmt.Scanln(&s)
	return object.SubroutineResult{SignalVal: &object.Value{InnerValue: object.String{Value: s}}, Trace: stacktrace.New(pos.BuildInPos("input"))}
}

func GenerateBuildIn(scope scope.Scope[*object.Value]) {
	AddGlobalFunction("println", Print(true, "println"), PrintArgs, scope)
	AddGlobalFunction("print", Print(false, "print"), PrintArgs, scope)
	AddGlobalFunction("input", Input, InputArguments, scope)
}
