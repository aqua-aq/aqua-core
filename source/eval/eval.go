package eval

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/source/ast"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/operators"
	"github.com/vandi37/aqua/source/signal"
	"github.com/vandi37/aqua/source/vm"
)

type Eval interface {
	Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult
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
	case ast.IdentExpression:
		return IdentExpression(val)
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
	IdentExpression     ast.IdentExpression
)

// Eval implements Eval.
func (i IdentExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	val, ok := scope.Get(i.Ident)
	if !ok {
		return object.ExpressionResult{
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.ValueError,
				Message: fmt.Sprintf("identifier %s is not defined", i.Ident),
			}},
		}
	}
	return object.ExpressionResult{
		SignalVal: val,
	}
}

// Eval implements Eval.
func (c CalLExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	sub := IntoEval(c.Subroutine).Eval(vm, scope, false)
	if sub.Signal.Has() {
		return sub
	}
	args := make([]*object.Value, 0, len(c.Args))
	for _, arg := range c.Args {
		argRes := IntoEval(arg).Eval(vm, scope, true)
		if argRes.Signal.Has() {
			return argRes
		}
		args = append(args, argRes.SignalVal)
	}
	return Call(vm, sub.SignalVal, args, clone)
}

// Eval implements Eval.
func (g GlobalSubroutineDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	sub := SubroutineDec(g.SubroutineDec).Eval(vm, scope, false)
	if sub.Signal.Has() {
		return sub
	}
	scope.Set(g.Name, sub.SignalVal.Normalize())
	return object.ExpressionResult{}
}

func (e ErrorDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Error(e)},
	}
}

func (s SignalExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	res := IntoEval(s.SigVal).Eval(vm, scope, false)
	if res.Signal.Has() {
		return res
	}
	return object.ExpressionResult{
		Signal:    s.Signal,
		SignalVal: res.SignalVal,
	}
}

func (w WhileExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	var (
		result *object.Value
		expr   object.ExpressionResult
		ok     bool
	)

start:
	if w.IsWhile {
		goto condition
	} else {
		goto block
	}
condition:
	expr = IntoEval(w.Condition).Eval(vm, scope, clone)
	if expr.Signal.Has() {
		return expr
	}

	ok, expr = IntoBool(vm, expr.SignalVal.Normalize())
	if expr.Signal.Has() {
		return expr
	}
	if (w.IsWhile && !ok) || (!w.IsWhile && ok) {
		goto after
	}
	expr = IntoEval(w.After).Eval(vm, scope, false)
	if expr.Signal.Has() {
		return expr
	}
block:
	expr = BlockExpression(w.Block).Eval(vm, scope, clone)
	switch {
	case expr.Signal == signal.SignalContinue:
		result = expr.SignalVal.Normalize()
		goto start
	case expr.Signal == signal.SignalBreak:
		result = expr.SignalVal.Normalize()
		goto after_break
	case expr.Signal.Has():
		return expr
	default:
		result = expr.SignalVal.Normalize()
	}

	goto condition
after_break:
	if w.Else != nil {
		return BlockExpression(*w.Else).Eval(vm, scope, clone)
	}
after:
	return object.ExpressionResult{SignalVal: result.Normalize()}
}

func (f ForExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	iter := IntoEval(f.Expression).Eval(vm, scope, false)
	if iter.Signal.Has() {
		return iter
	}
	iter = IntoIter(iter.SignalVal.Normalize(), vm)
	if iter.Signal.Has() {
		return iter
	}
	arguments := object.Arguments{Last: f.Arguments.Last}
	for _, arg := range f.Arguments.Elements {
		res := IntoEval(arg.Default).Eval(vm, scope, true)
		if res.Signal.Has() {
			return res
		}
		arguments.Elements = append(arguments.Elements, object.Argument{
			Name:    arg.Name,
			Default: res.SignalVal.Normalize(),
		})
	}
	var (
		result    *object.Value
		iteration int
		expr      object.ExpressionResult
	)
	declareVals := func(args []*object.Value) {
		scope = scope.Rebase()
		if f.IsEnum {
			args = append([]*object.Value{{InnerValue: object.Int{Value: iteration}}}, args...)
		}
		iteration++
		object.ParseArgs(arguments, args, scope)
	}
	scope = scope.Push()
	switch val := iter.SignalVal.Normalize().InnerValue.(type) {
	case object.Array:
		for _, element := range val.Elements {
			declareVals([]*object.Value{element})
			expr = BlockExpression(f.Block).Eval(vm, scope, clone)
			switch {
			case expr.Signal == signal.SignalContinue:
				result = expr.SignalVal.Normalize()
			case expr.Signal == signal.SignalBreak:
				result = expr.SignalVal.Normalize()
				goto after_break
			case expr.Signal.Has():
				return expr
			default:
				result = expr.SignalVal.Normalize()
			}
		}
	case object.Object:
		for k, v := range val.Map {
			declareVals([]*object.Value{{InnerValue: object.String{Value: k}}, v})
			expr = BlockExpression(f.Block).Eval(vm, scope, clone)
			switch {
			case expr.Signal == signal.SignalContinue:
				result = expr.SignalVal.Normalize()
			case expr.Signal == signal.SignalBreak:
				result = expr.SignalVal.Normalize()
				goto after_break
			case expr.Signal.Has():
				return expr
			default:
				result = expr.SignalVal.Normalize()
			}
		}
	case object.Method, *object.Subroutine:
		for {
			expr = Call(vm, iter.SignalVal.Normalize(), nil, false)
			if expr.Signal == signal.SignalRaise && (expr.SignalVal.Normalize().InnerValue == object.Error{Code: errors.IteratorStop}) {
				break
			}
			if expr.Signal.Has() {
				return expr
			}
			var elements []*object.Value
			if val, ok := expr.SignalVal.Normalize().InnerValue.(object.Array); ok {
				elements = val.Elements
			} else {
				elements = []*object.Value{expr.SignalVal.Normalize()}
			}
			declareVals(elements)
			expr = BlockExpression(f.Block).Eval(vm, scope, clone)
			switch {
			case expr.Signal == signal.SignalContinue:
				result = expr.SignalVal.Normalize()
			case expr.Signal == signal.SignalBreak:
				result = expr.SignalVal.Normalize()
				goto after_break
			case expr.Signal.Has():
				return expr
			default:
				result = expr.SignalVal.Normalize()
			}
		}
	case object.String:
		for _, r := range val.Value {
			declareVals([]*object.Value{{InnerValue: object.String{Value: string(r)}}})
			expr = BlockExpression(f.Block).Eval(vm, scope, clone)
			switch {
			case expr.Signal == signal.SignalContinue:
				result = expr.SignalVal.Normalize()
			case expr.Signal == signal.SignalBreak:
				result = expr.SignalVal.Normalize()
				goto after_break
			case expr.Signal.Has():
				return expr
			default:
				result = expr.SignalVal.Normalize()
			}
		}
	default:
		return object.ExpressionResult{
			Signal: signal.SignalRaise,
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected iterable, got %v", iter.SignalVal.Normalize().Type()),
			}},
		}
	}
	goto after
after_break:
	if f.Else != nil {
		return BlockExpression(*f.Else).Eval(vm, scope, clone)
	}
after:
	return object.ExpressionResult{SignalVal: result.Normalize()}
}

func (i IfExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	expr := IntoEval(i.Condition).Eval(vm, scope, false)
	if expr.Signal.Has() {
		return expr
	}
	ok, expr := IntoBool(vm, expr.SignalVal.Normalize())
	if expr.Signal.Has() {
		return expr
	}
	if ok {
		return BlockExpression(i.If).Eval(vm, scope, clone)
	}
	for _, next := range i.ElseIfs {
		expr := IntoEval(next.Condition).Eval(vm, scope, false)
		ok, expr := IntoBool(vm, expr.SignalVal.Normalize())
		if ok {
			return BlockExpression(next.Block).Eval(vm, scope, clone)
		}
	}
	if i.Else != nil {
		return BlockExpression(*i.Else).Eval(vm, scope, clone)
	}
	return object.ExpressionResult{}
}

func (p PrefixExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	expr := IntoEval(p.Value).Eval(vm, scope, false)
	if expr.Signal.Has() {
		return expr
	}
	typeRaise := object.ExpressionResult{
		Signal: signal.SignalRaise,
		SignalVal: &object.Value{
			InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("unsupported type %s for prefix operator '%s'", expr.SignalVal.Normalize().Type(), p.Operator.String()),
			},
		},
	}
	switch val := expr.SignalVal.Normalize().InnerValue.(type) {
	case object.Int:
		switch p.Operator {
		case operators.Neg:
			return object.ExpressionResult{SignalVal: &object.Value{InnerValue: object.Int{Value: -val.Value}}}
		case operators.Not:
			return object.ExpressionResult{SignalVal: &object.Value{InnerValue: object.Int{Value: ^val.Value}}}
		case operators.Inc:
			expr.SignalVal.Normalize().InnerValue = object.Int{Value: val.Value + 1}
			return expr.Clone(clone)
		case operators.Dec:
			expr.SignalVal.Normalize().InnerValue = object.Int{Value: val.Value - 1}
			return expr.Clone(clone)
		default:
			return typeRaise
		}
	case object.Number:
		switch p.Operator {
		case operators.Neg:
			return object.ExpressionResult{SignalVal: &object.Value{InnerValue: object.Number{Value: -val.Value}}}
		case operators.Inc:
			expr.SignalVal.Normalize().InnerValue = object.Number{Value: val.Value + 1}
			return expr.Clone(clone)
		case operators.Dec:
			expr.SignalVal.Normalize().InnerValue = object.Number{Value: val.Value - 1}
			return expr.Clone(clone)
		default:
			return typeRaise
		}
	}
	method := GetAttrMethod(expr.SignalVal.Normalize(), p.Operator.Method())
	if method.Signal.Has() {
		return method
	}
	return Call(vm, method.SignalVal.Normalize(), nil, clone)
}

func (b BinExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	panic("unimplemented")
}

func (o ObjectDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	obj := make(map[string]*object.Value)
	for _, val := range o.Vals {
		if val.IsContinuos {
			res := IntoEval(val.Value).Eval(vm, scope, true)
			if res.Signal.Has() {
				return res
			}
			inner, ok := res.SignalVal.InnerValue.(object.Object)
			if !ok {
				return object.ExpressionResult{
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
		res := IntoEval(val.Value).Eval(vm, scope, true)
		if res.Signal.Has() {
			return res
		}
		obj[val.Name] = res.SignalVal.Normalize()
	}
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Object{
			Map: obj,
		}},
	}
}
func (s SubroutineDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	arguments := object.Arguments{Last: s.Arguments.Last}
	for _, arg := range s.Arguments.Elements {
		res := IntoEval(arg.Default).Eval(vm, scope, true)
		if res.Signal.Has() {
			return res
		}
		arguments.Elements = append(arguments.Elements, object.Argument{
			Name:    arg.Name,
			Default: res.SignalVal.Normalize(),
		})
	}
	res := IntoEval(s.Body).Eval(vm, scope, true)
	if res.Signal.Has() {
		return res
	}
	return object.ExpressionResult{
		SignalVal: &object.Value{
			InnerValue: &object.Subroutine{
				Arguments: arguments,
				Scope:     scope,
				Prototype: res.SignalVal.Normalize(),
				Code:      s.Body,
			},
		},
	}
}

func (i IntDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Int{Value: int(i)}},
	}
}

func (n NumDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Number{Value: float64(n)}},
	}
}
func (s StringDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.String{Value: string(s)}},
	}
}
func (n NullDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Null{}},
	}
}

func (a ArrayDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	var arr []*object.Value
	for _, val := range a.Elements {
		if val.IsContinuos {
			res := IntoEval(val.Value).Eval(vm, scope, true)
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
		res := IntoEval(val.Value).Eval(vm, scope, true)
		if res.Signal.Has() {
			return res
		}
		arr = append(arr, res.SignalVal.Normalize())
	}
	return object.ExpressionResult{
		SignalVal: &object.Value{InnerValue: object.Array{Elements: arr}},
	}
}

func (l LetExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	val := &object.Value{InnerValue: object.Null{}}
	scope.Set(l.Name, val)
	return object.ExpressionResult{
		SignalVal: val,
	}
}

func (b BlockExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	scope = scope.Push()
	for i, expr := range b.Expressions {
		res := IntoEval(expr).Eval(vm, scope, true)
		if res.Signal == signal.SignalRaise && b.Catch != nil {
			scope = scope.Rebase()
			scope.Set(b.Catch.Name, res.SignalVal.Normalize())
			return BlockExpression(b.Catch.Expressions).Eval(vm, scope, true)
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
