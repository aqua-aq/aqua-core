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
	case ast.BinExpression:
		return BinExpression(val)
	case ast.PrefixExpression:
		return PrefixExpression(val)
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
	case ast.SubroutineDec:
		return SubroutineDec(val)
	case ast.SignalExpression:
		return SignalExpression(val)
	default:
		return NullDec{}
	}
}

type (
	ObjectDec        ast.ObjectDec
	IntDec           ast.IntDec
	NumDec           ast.NumDec
	StringDec        ast.StringDec
	NullDec          ast.NullDec
	ErrorDec         ast.ErrorDec
	ArrayDec         ast.ArrayDec
	BinExpression    ast.BinExpression
	PrefixExpression ast.PrefixExpression
	LetExpression    ast.LetExpression
	BlockExpression  ast.BlockExpression
	IfExpression     ast.IfExpression
	ForExpression    ast.ForExpression
	WhileExpression  ast.WhileExpression
	SubroutineDec    ast.SubroutineDec
	SignalExpression ast.SignalExpression
)

func (e ErrorDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	return object.ExpressionResult{
		Value: &object.Value{InnerValue: object.Error(e)},
	}
}

func (s SignalExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	res := IntoEval(s.SigVal).Eval(vm, scope)
	if res.Signal.Has() {
		return res
	}
	return object.ExpressionResult{
		Signal:    s.Signal,
		SignalVal: res.Value,
	}
}

func (s SubroutineDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	panic("unimplemented")
}

func (w WhileExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	res := &object.Value{InnerValue: object.Null{}}
	var hasBreak bool
	if w.IsWhile {
		for {
			cond := IntoEval(w.Condition).Eval(vm, scope)
			if cond.Signal.Has() {
				return cond
			}
			ok, sig := IntoBool(cond.Value.Normalize())
			if sig.Signal.Has() {
				return sig
			}
			if ok {
				break
			}
			block := BlockExpression(w.Block).Eval(vm, scope)
			if block.Signal == signal.SignalContinue {
				res = block.SignalVal.Normalize()
				continue
			}
			if block.Signal == signal.SignalBreak {
				res = block.SignalVal.Normalize()
				hasBreak = true
				break
			}
			if block.Signal.Has() {
				return block
			}
			res = block.Value.Normalize()
		}
	} else {
		for {
			block := BlockExpression(w.Block).Eval(vm, scope)
			if block.Signal == signal.SignalContinue {
				res = block.SignalVal.Normalize()
				continue
			}
			if block.Signal == signal.SignalBreak {
				res = block.SignalVal.Normalize()
				hasBreak = true
				break
			}
			if block.Signal.Has() {
				return block
			}
			res = block.Value.Normalize()
			cond := IntoEval(w.Condition).Eval(vm, scope)
			if cond.Signal.Has() {
				return cond
			}
			ok, sig := IntoBool(cond.Value.Normalize())
			if sig.Signal.Has() {
				return sig
			}
			if !ok {
				break
			}
		}

	}
	if !hasBreak && w.Else != nil {
		return BlockExpression(*w.Else).Eval(vm, scope)
	}
	return object.ExpressionResult{
		Value: res.Normalize(),
	}
}

func (f ForExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
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
	obj := make(map[string]*object.Value)
	for _, val := range o.Vals {
		if val.IsContinuos {
			res := IntoEval(val.Value).Eval(vm, scope)
			if res.Signal.Has() {
				return res
			}
			inner, ok := res.Value.InnerValue.(object.Object)
			if !ok {
				return object.ExpressionResult{
					Signal: signal.SignalRaise,
					SignalVal: &object.Value{InnerValue: object.Error{
						Code:    errors.TypeError,
						Message: fmt.Sprintf("expected object, got %v", res.Value.Normalize().Type()),
					}},
				}
			}
			for k, v := range inner.Map {
				obj[k] = v
			}
		}
		res := IntoEval(val.Value).Eval(vm, scope)
		if res.Signal.Has() {
			return res
		}
		obj[val.Name] = res.Value.Normalize()
	}
	return object.ExpressionResult{
		Value: &object.Value{InnerValue: object.Object{
			Map: obj,
		}},
	}
}

func (i IntDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	return object.ExpressionResult{
		Value: &object.Value{InnerValue: object.Int{Value: int(i)}},
	}
}

func (n NumDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	return object.ExpressionResult{
		Value: &object.Value{InnerValue: object.Number{Value: float64(n)}},
	}
}
func (s StringDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	return object.ExpressionResult{
		Value: &object.Value{InnerValue: object.String{Value: string(s)}},
	}
}
func (n NullDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	return object.ExpressionResult{
		Value: &object.Value{InnerValue: object.Null{}},
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
			inner, ok := res.Value.InnerValue.(*object.Array)
			if !ok {
				return object.ExpressionResult{
					Signal: signal.SignalRaise,
					SignalVal: &object.Value{InnerValue: object.Error{
						Code:    errors.TypeError,
						Message: fmt.Sprintf("expected array, got %v", res.Value.Normalize().Type()),
					}},
				}
			}
			arr = append(arr, inner.Elements...)
		}
		res := IntoEval(val.Value).Eval(vm, scope)
		if res.Signal.Has() {
			return res
		}
		arr = append(arr, res.Value.Normalize())
	}
	return object.ExpressionResult{
		Value: &object.Value{InnerValue: object.Array{Elements: arr}},
	}
}

func (l LetExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	val := &object.Value{InnerValue: object.Null{}}
	scope.Set(l.Name, val)
	return object.ExpressionResult{
		Value: val,
	}
}

func (b BlockExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value]) object.ExpressionResult {
	scope = scope.Push()
	for i, expr := range b.Expressions {
		res := IntoEval(expr).Eval(vm, scope)
		if res.Signal == signal.SignalRaise && b.Catch != nil {
			scope = scope.Rebase()
			scope.Set(b.Catch.Name, res.SignalVal)
			res = BlockExpression(b.Catch.Expressions).Eval(vm, scope)
			return res
		}
		if res.Signal.Has() {
			return res
		}
		if i == len(b.Expressions)-1 {
			return res
		}
	}
	return object.ExpressionResult{
		Value: &object.Value{InnerValue: object.Null{}},
	}
}
