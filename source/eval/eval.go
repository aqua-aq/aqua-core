package eval

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/source/ast"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/signal"
	"github.com/vandi37/aqua/source/vm"
)

type Eval interface {
	Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult
}

func IntoEval(expr ast.Expression) Eval {
	switch val := expr.(type) {
	case ast.ObjectDec:
		return ObjectDec(val)
	case ast.IntDec:
		return IntDec(val)
	case ast.NumDec:
		return NumDec(val)
	case ast.StringDec:
		return StringDec(val)
	case ast.NullDec:
		return NullDec(val)
	case ast.ErrorDec:
		return ErrorDec(val)
	case ast.ArrayDec:
		return ArrayDec(val)
	case ast.SubroutineDec:
		return SubroutineDec(val)
	case ast.BinExpression:
		return BinExpression(val)
	case ast.PrefixExpression:
		return PrefixExpression(val)
	case ast.CalLExpression:
		return CalLExpression(val)
	case ast.LetExpression:
		return LetExpression(val)
	case ast.BlockExpression:
		return BlockExpression(val)
	case ast.IfExpression:
		return IfExpression(val)
	case ast.ForExpression:
		return ForExpression(val)
	case ast.WhileExpression:
		return WhileExpression(val)
	case ast.GlobalSubroutineDec:
		return GlobalSubroutineDec(val)
	case ast.SignalExpression:
		return SignalExpression(val)
	default:
		return NullDec{}
	}
}

type (
	ObjectDec           ast.ObjectDec
	IntDec              ast.IntDec
	NumDec              ast.NumDec
	StringDec           ast.StringDec
	NullDec             ast.NullDec
	ErrorDec            ast.ErrorDec
	ArrayDec            ast.ArrayDec
	SubroutineDec       ast.SubroutineDec
	BinExpression       ast.BinExpression
	PrefixExpression    ast.PrefixExpression
	CalLExpression      ast.CalLExpression
	LetExpression       ast.LetExpression
	BlockExpression     ast.BlockExpression
	IfExpression        ast.IfExpression
	ForExpression       ast.ForExpression
	WhileExpression     ast.WhileExpression
	GlobalSubroutineDec ast.GlobalSubroutineDec
	SignalExpression    ast.SignalExpression
)

// Eval implements Eval.
func (c CalLExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	sub := IntoEval(c.Subroutine).Eval(vm, scope)
	if sub.Signal.Has() {
		return sub
	}
	args := make([]*object.Value, 0, len(c.Args))
	for _, arg := range c.Args {
		argRes := IntoEval(arg).Eval(vm, scope)
		if argRes.Signal.Has() {
			return argRes
		}
		args = append(args, argRes.SignalVal)
	}
	return Call(vm, sub.SignalVal, args)
}

// Eval implements Eval.
func (g GlobalSubroutineDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	sub := SubroutineDec(g.SubroutineDec).Eval(vm, scope)
	if sub.Signal.Has() {
		return sub
	}
	scope.Set(g.Name, sub.SignalVal.Normalize())
	return object.ExpressionResult{}
}

func (e ErrorDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Error(e)},
	}
}

func (s SignalExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	res := IntoEval(s.SigVal).Eval(vm, scope)
	if res.Signal.Has() {
		return res
	}
	return object.ExpressionResult{
		Signal:    s.Signal,
		SignalVal: res.SignalVal,
	}
}

func (w WhileExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	var (
		result   *object.Value
		hasBreak bool
		expr     object.ExpressionResult
		ok       bool
	)

start:
	if w.IsWhile {
		goto condition
	} else {
		goto block
	}
condition:
	expr = IntoEval(w.Condition).Eval(vm, scope)
	if expr.Signal == signal.SignalRaise && w.Catch != nil {
		goto catch
	}
	if expr.Signal.Has() {
		return expr
	}

	ok, expr = IntoBool(vm, expr.SignalVal.Normalize())
	if expr.Signal == signal.SignalRaise && w.Catch != nil {
		goto catch
	}
	if expr.Signal.Has() {
		return expr
	}
	if (w.IsWhile && !ok) || (!w.IsWhile && ok) {
		goto after
	}
block:
	expr = BlockExpression(w.Block).Eval(vm, scope)
	switch {
	case expr.Signal == signal.SignalRaise && w.Catch != nil:
		goto catch
	case expr.Signal == signal.SignalContinue:
		result = expr.SignalVal.Normalize()
		goto start
	case expr.Signal == signal.SignalBreak:
		result = expr.SignalVal.Normalize()
		hasBreak = true
		goto after
	case expr.Signal.Has():
		return expr
	default:
		result = expr.SignalVal.Normalize()
	}

	goto condition

after:
	if !hasBreak && w.Else != nil {
		return BlockExpression(*w.Else).Eval(vm, scope)
	}
	return object.ExpressionResult{SignalVal: result.Normalize()}
catch:
	scope = scope.Rebase()
	scope.Set(w.Catch.Name, expr.SignalVal.Normalize())
	return BlockExpression(w.Catch.Expressions).Eval(vm, scope)
}

func (f ForExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	panic("unimplemented")
}

func (i IfExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	panic("unimplemented")
}

func (p PrefixExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	panic("unimplemented")
}

func (b BinExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	panic("unimplemented")
}

func (o ObjectDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	obj, res := createMap(vm, scope, o)
	if res.Signal.Has() {
		return res
	}
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Object{
			Map: obj,
		}},
	}
}
func createMap(vm *vm.VM, scope scope.Scope[*object.Value], o ObjectDec) (map[string]*object.Value, object.ExpressionResult) {
	obj := make(map[string]*object.Value)
	for _, val := range o.Vals {
		if val.IsContinuos {
			res := IntoEval(val.Value).Eval(vm, scope)
			if res.Signal.Has() {
				return nil, res
			}
			inner, ok := res.SignalVal.InnerValue.(object.Object)
			if !ok {
				return nil, object.ExpressionResult{
					Signal: signal.SignalRaise,
					SignalVal: &object.Value{InnerValue: object.Error{
						Code:    errors.TypeError,
						Message: fmt.Sprintf("expected object, got %v", res.SignalVal.Normalize().Type()),
					}},
				}
			}
			for k, v := range inner.Map {
				obj[k] = v
			}
		}
		res := IntoEval(val.Value).Eval(vm, scope)
		if res.Signal.Has() {
			return nil, res
		}
		obj[val.Name] = res.SignalVal.Normalize()
	}
	return obj, object.ExpressionResult{}
}
func (s SubroutineDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	arguments := object.Arguments{Last: s.Arguments.Last}
	for _, arg := range s.Arguments.Elements {
		res := IntoEval(arg.Default).Eval(vm, scope)
		if res.Signal.Has() {
			return res
		}
		arguments.Elements = append(arguments.Elements, object.Argument{
			Name:    arg.Name,
			Default: res.SignalVal.Normalize(),
		})
	}
	obj, res := createMap(vm, scope, ObjectDec(s.Prototype))
	if res.Signal.Has() {
		return res
	}
	return object.ExpressionResult{
		SignalVal: &object.Value{
			InnerValue: &object.Subroutine{
				Arguments: arguments,
				Scope:     scope,
				Prototype: obj,
				Code:      s.Body,
			},
		},
	}
}

func (i IntDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Int{Value: int(i)}},
	}
}

func (n NumDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Number{Value: float64(n)}},
	}
}
func (s StringDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.String{Value: string(s)}},
	}
}
func (n NullDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Null{}},
	}
}

func (a ArrayDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	var arr []*object.Value
	for _, val := range a.Elements {
		if val.IsContinuos {
			res := IntoEval(val.Value).Eval(vm, scope)
			if res.Signal.Has() {
				return res
			}
			inner, ok := res.SignalVal.InnerValue.(*object.Array)
			if !ok {
				return object.ExpressionResult{
					Signal: signal.SignalRaise,
					SignalVal: &object.Value{InnerValue: object.Error{
						Code:    errors.TypeError,
						Message: fmt.Sprintf("expected array, got %v", res.SignalVal.Normalize().Type()),
					}},
				}
			}
			arr = append(arr, inner.Elements...)
		}
		res := IntoEval(val.Value).Eval(vm, scope)
		if res.Signal.Has() {
			return res
		}
		arr = append(arr, res.SignalVal.Normalize())
	}
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Array{Elements: arr}},
	}
}

func (l LetExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	val := &object.Value{InnerValue: object.Null{}}
	scope.Set(l.Name, val)
	return object.ExpressionResult{
		SignalVal: val,
	}
}

func (b BlockExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	scope = scope.Push()
	for i, expr := range b.Expressions {
		res := IntoEval(expr).Eval(vm, scope)
		if res.Signal == signal.SignalRaise && b.Catch != nil {
			scope = scope.Rebase()
			scope.Set(b.Catch.Name, res.SignalVal.Normalize())
			return BlockExpression(b.Catch.Expressions).Eval(vm, scope)
		}
		if res.Signal.Has() {
			return res
		}
		if i == len(b.Expressions)-1 {
			return res
		}
	}
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Null{}},
	}
}
