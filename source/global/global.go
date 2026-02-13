package global

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/eval"
	"github.com/vandi37/aqua/source/keywords"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
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
		BuiltIn:   f,
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
			if sErr := err.IntoSubroutineResultStrict(pos.BuildInPos(name)); err.Signal.Has() {
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
			Message: fmt.Sprintf("expected argument %s", firstArgument),
		})
	}
	str, err := eval.IntoString(vm, prompt, pos.BuildInPos("input"))
	if sErr := err.IntoSubroutineResultStrict(pos.BuildInPos("input")); err.Signal.Has() {
		return sErr
	}
	fmt.Print(str)

	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	return object.SubroutineResult{SignalVal: &object.Value{InnerValue: object.String{Value: s}}, Trace: stacktrace.New(pos.BuildInPos("input"))}
}

var FirstArguments = object.Arguments{Elements: []object.Argument{{Name: firstArgument}}}

func Bool(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
	first, ok := scope.Get(firstArgument)
	if !ok {
		return Raise("bool", object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("expected argument %s", firstArgument),
		})
	}

	b, err := eval.IntoBool(vm, first.Normalize(), pos.BuildInPos("bool"))
	if sErr := err.IntoSubroutineResultStrict(pos.BuildInPos("bool")); err.Signal.Has() {
		return sErr
	}
	return object.SubroutineResult{SignalVal: &object.Value{InnerValue: object.Bool{Value: b}}, Trace: stacktrace.New(pos.BuildInPos("number"))}
}
func Number(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
	first, ok := scope.Get(firstArgument)
	if !ok {
		return Raise("number", object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("expected argument %s", firstArgument),
		})
	}

	n, err := eval.IntoNum(vm, first.Normalize(), pos.BuildInPos("number"))
	if sErr := err.IntoSubroutineResultStrict(pos.BuildInPos("number")); err.Signal.Has() {
		return sErr
	}
	return object.SubroutineResult{SignalVal: &object.Value{InnerValue: object.Number{Value: n}}, Trace: stacktrace.New(pos.BuildInPos("number"))}
}

func Eval(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
	first, ok := scope.Get(firstArgument)
	if !ok {
		return Raise("eval", object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("expected argument %s", firstArgument),
		})
	}
	str, err := eval.IntoString(vm, first, pos.BuildInPos("eval"))
	if sErr := err.IntoSubroutineResultStrict(pos.BuildInPos("eval")); err.Signal.Has() {
		return sErr
	}
	err = eval.Run(vm, scope, str, pos.BuildInPos("eval"), false)
	return err.IntoSubroutineResultStrict(pos.BuildInPos("eval"))
}

func Len(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
	first, ok := scope.Get(firstArgument)
	if !ok {
		return Raise("eval", object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("expected argument %s", firstArgument),
		})
	}

	var l int
	switch v := first.Normalize().InnerValue.(type) {
	case object.Object:
		if eval.AttrExists(first, keywords.Len) {
			method := eval.GetAttrMethod(first, keywords.Len, pos.BuildInPos("len"))
			if err := method.IntoSubroutineResultStrict(pos.BuildInPos("len")); method.Signal.Has() {
				return err
			}
			return eval.Call(
				vm,
				method.SignalVal.Normalize(),
				[]*object.Value{},
				false,
				pos.BuildInPos("len"),
				nil,
			).IntoSubroutineResultStrict(pos.BuildInPos("len"))
		}
		l = len(v.Map)
	case object.Array:
		l = len(v.Elements)
	case object.String:
		l = len(v.Value)
	default:
		return Raise("len", object.Error{
			Code:    errors.TypeError,
			Message: fmt.Sprintf("can't get length of %s", first.Normalize().Type())})
	}
	return object.SubroutineResult{SignalVal: &object.Value{InnerValue: object.Number{Value: float64(l)}}, Trace: stacktrace.New(pos.BuildInPos("len"))}
}

func String(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
	first, ok := scope.Get(firstArgument)
	if !ok {
		return Raise("string", object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("expected argument %s", firstArgument),
		})
	}
	str, err := eval.IntoString(vm, first, pos.BuildInPos("string"))
	if sErr := err.IntoSubroutineResultStrict(pos.BuildInPos("string")); err.Signal.Has() {
		return sErr
	}
	return object.SubroutineResult{SignalVal: &object.Value{InnerValue: object.String{Value: str}}, Trace: stacktrace.New(pos.BuildInPos("string"))}
}

func Code(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
	first, ok := scope.Get(firstArgument)
	if !ok {
		return Raise("code", object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("expected argument %s", firstArgument),
		})
	}
	if err, ok := first.Normalize().InnerValue.(object.Error); ok {
		return object.SubroutineResult{SignalVal: &object.Value{InnerValue: object.Number{Value: float64(err.Code)}}, Trace: stacktrace.New(pos.BuildInPos("code"))}
	}
	return Raise("code", object.Error{
		Code:    errors.TypeError,
		Message: fmt.Sprintf("can't get code of %s", first.Normalize().Type()),
	})
}

func Message(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
	first, ok := scope.Get(firstArgument)
	if !ok {
		return Raise("message", object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("expected argument %s", firstArgument),
		})
	}
	if err, ok := first.Normalize().InnerValue.(object.Error); ok {
		return object.SubroutineResult{SignalVal: &object.Value{InnerValue: object.String{Value: err.Message}}, Trace: stacktrace.New(pos.BuildInPos("message"))}
	}
	return object.SubroutineResult{
		Signal: signal.SubroutineSignalRaise,
		SignalVal: &object.Value{InnerValue: object.Error{
			Code:    errors.TypeError,
			Message: fmt.Sprintf("can't get message of %s", first.Normalize().Type()),
		}},
		Trace: stacktrace.New(pos.BuildInPos("message")),
	}
}

var code = "code"
var message = "message"
var ErrorArguments = object.Arguments{Elements: []object.Argument{{
	Name: code,
}, {
	Name:    message,
	Default: &object.Value{InnerValue: object.String{Value: ""}},
}}}

func Error(vm *vm.VM[*object.Value], scope scope.Scope[*object.Value]) object.SubroutineResult {
	c, ok := scope.Get(code)
	if !ok {
		return Raise("error", object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("expected argument %s", code),
		})
	}
	intCode, err := eval.IntoInt(vm, c.Normalize(), pos.BuildInPos("error"))
	if sErr := err.IntoSubroutineResultStrict(pos.BuildInPos("error")); err.Signal.Has() {
		return sErr
	}
	m, ok := scope.Get(message)
	if !ok {
		return Raise("error", object.Error{
			Code:    errors.ValueError,
			Message: fmt.Sprintf("expected argument %s", message),
		})
	}
	strMessage, err := eval.IntoString(vm, m.Normalize(), pos.BuildInPos("error"))
	if sErr := err.IntoSubroutineResultStrict(pos.BuildInPos("error")); err.Signal.Has() {
		return sErr
	}
	return object.SubroutineResult{SignalVal: &object.Value{InnerValue: object.Error{Code: errors.Code(intCode), Message: strMessage}}, Trace: stacktrace.New(pos.BuildInPos("error"))}
}

func GenerateBuildIn(scope scope.Scope[*object.Value]) {
	AddGlobalFunction("println", Print(true, "println"), PrintArgs, scope)
	AddGlobalFunction("print", Print(false, "print"), PrintArgs, scope)
	AddGlobalFunction("input", Input, InputArguments, scope)
	AddGlobalFunction("bool", Bool, FirstArguments, scope)
	AddGlobalFunction("number", Number, FirstArguments, scope)
	AddGlobalFunction("eval", Eval, FirstArguments, scope)
	AddGlobalFunction("len", Len, FirstArguments, scope)
	AddGlobalFunction("string", String, FirstArguments, scope)
	AddGlobalFunction("code", Code, FirstArguments, scope)
	AddGlobalFunction("message", Message, FirstArguments, scope)
	AddGlobalFunction("error", Error, ErrorArguments, scope)
	scope.Set("stop", &object.Value{InnerValue: object.Error{Code: errors.IteratorStop}})
}
