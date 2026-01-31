package eval

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/pkg/stacktrace"
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
	case ast.NumDec:
		return NumDec(val)
	case ast.BoolDec:
		return BoolDec(val)
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
	case ast.CallExpression:
		return CallExpression(val)
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
	case ast.SignalExpression:
		return SignalExpression(val)
	case ast.IdentExpression:
		return IdentExpression(val)
	case ast.AssigmentExpression:
		return AssigmentExpression(val)
	case ast.ModExpression:
		return ModExpression(val)
	case ast.ImportExpression:
		return ImportExpression(val)
	default:
		return NullDec{}
	}
}

type (
	ObjectDec           ast.ObjectDec
	NumDec              ast.NumDec
	BoolDec             ast.BoolDec
	StringDec           ast.StringDec
	NullDec             ast.NullDec
	ErrorDec            ast.ErrorDec
	ArrayDec            ast.ArrayDec
	SubroutineDec       ast.SubroutineDec
	BinExpression       ast.BinExpression
	PrefixExpression    ast.PrefixExpression
	CallExpression      ast.CallExpression
	LetExpression       ast.LetExpression
	BlockExpression     ast.BlockExpression
	IfExpression        ast.IfExpression
	ForExpression       ast.ForExpression
	WhileExpression     ast.WhileExpression
	SignalExpression    ast.SignalExpression
	IdentExpression     ast.IdentExpression
	AssigmentExpression ast.AssigmentExpression
	ModExpression       ast.ModExpression
	ImportExpression    ast.ImportExpression
)

func (i ImportExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	panic("unimplemented")
}

func (m ModExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	expr := IntoEval(m.Name).Eval(vm, scope, false)
	if expr.Signal.Has() {
		return expr.Clone(clone)
	}
	str, expr := IntoString(vm, expr.SignalVal.Normalize(), m.Pos)
	if expr.Signal.Has() {
		return expr.Clone(clone)
	}

	expr = DeclareSubroutine(vm, scope, false, fmt.Sprintf("<%s>", str), ast.SubroutineDec{
		Arguments: ast.Arguments{},
		Body:      m.Body,
		Prototype: ast.NullDec{Pos: m.Pos},
		Pos:       m.Pos,
	})
	if expr.Signal.Has() {
		return expr.Clone(clone)
	}
	export := make(map[string]*object.Value, len(m.Export))
	for _, v := range m.Export {
		export[v] = &object.Value{InnerValue: object.Null{}}
	}
	expr = Call(vm, expr.SignalVal.Normalize(), []*object.Value{}, clone, m.Pos, export)
	if expr.Signal.Has() {
		return expr
	}
	return object.ExpressionResult{SignalVal: &object.Value{InnerValue: object.Object{
		Map: export,
	}}}
}

func (a AssigmentExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	if len(a.Left) != len(a.Right) {
		return object.ExpressionResult{
			Trace:  stacktrace.New(a.Pos),
			Signal: signal.SignalRaise,
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.ValueError,
				Message: fmt.Sprintf("expected %d right values, got %d", len(a.Left), len(a.Right)),
			}},
		}
	}
	for i, v := range a.Left {
		left := IntoEval(v).Eval(vm, scope, false)
		if left.Signal.Has() {
			return left.Clone(clone)
		}
		right := IntoEval(a.Right[i]).Eval(vm, scope, true)
		if right.Signal.Has() {
			return right.Clone(clone)
		}
		expr := RunBin(vm, scope, false, left.SignalVal.Normalize(), right.SignalVal.Normalize(), a.Operator, a.Pos)
		if expr.Signal.Has() {
			return expr.Clone(clone)
		}
		*left.SignalVal.Normalize() = *expr.SignalVal.Normalize()
	}
	return object.ExpressionResult{Trace: stacktrace.New(a.Pos)}
}

func (i IdentExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	val, ok := scope.Get(i.Ident)
	if !ok {
		return object.ExpressionResult{Trace: stacktrace.New(i.Pos),
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.ValueError,
				Message: fmt.Sprintf("identifier %s is not defined", i.Ident),
			}},
			Signal: signal.SignalRaise,
		}
	}
	return object.ExpressionResult{Trace: stacktrace.New(i.Pos),
		SignalVal: val,
	}.Clone(clone)
}

func (c CallExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	sub := IntoEval(c.Subroutine).Eval(vm, scope, false)
	if sub.Signal.Has() {
		return sub
	}
	args := make([]*object.Value, 0, len(c.Args))
	for _, val := range c.Args {
		if val.IsContinuos {
			res := IntoEval(val.Value).Eval(vm, scope, true)
			if res.Signal.Has() {
				return res.Clone(clone)
			}
			inner, ok := res.SignalVal.InnerValue.(object.Array)
			if !ok {
				return object.ExpressionResult{Trace: stacktrace.New(val.Pos),
					Signal: signal.SignalRaise,
					SignalVal: &object.Value{InnerValue: object.Error{
						Code:    errors.TypeError,
						Message: fmt.Sprintf("expected array, got %v", res.SignalVal.Normalize().Type()),
					}},
				}
			}
			args = append(args, inner.Elements...)
		}
		res := IntoEval(val.Value).Eval(vm, scope, true)
		if res.Signal.Has() {
			return res.Clone(clone)
		}
		args = append(args, res.SignalVal.Normalize())
	}
	return Call(vm, sub.SignalVal, args, clone, c.Pos, nil)
}

func (e ErrorDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return object.ExpressionResult{Trace: stacktrace.New(e.Pos),
		SignalVal: &object.Value{InnerValue: object.Error(e.Value)},
	}
}

func (s SignalExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	res := IntoEval(s.SigVal).Eval(vm, scope, false)
	if res.Signal.Has() {
		return res.Clone(clone)
	}
	return object.ExpressionResult{Trace: stacktrace.New(s.Pos),
		Signal:    s.Signal,
		SignalVal: res.SignalVal,
	}.Clone(clone)
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
	expr = IntoEval(w.Condition).Eval(vm, scope, false)
	if expr.Signal.Has() {
		return expr.Clone(clone)
	}

	ok, expr = IntoBool(vm, expr.SignalVal.Normalize(), w.Pos)
	if expr.Signal.Has() {
		return expr.Clone(clone)
	}
	if (w.IsWhile && !ok) || (!w.IsWhile && ok) {
		goto after
	}
	expr = IntoEval(w.After).Eval(vm, scope, false)
	if expr.Signal.Has() {
		return expr.Clone(clone)
	}
block:
	expr = BlockExpression(w.Block).Eval(vm, scope, false)
	switch {
	case expr.Signal == signal.SignalContinue:
		result = expr.SignalVal.Normalize()
		goto start
	case expr.Signal == signal.SignalBreak:
		result = expr.SignalVal.Normalize()
		goto after
	case expr.Signal.Has():
		return expr
	default:
		result = expr.SignalVal.Normalize()
	}

	goto condition

after:
	if w.Else != nil && result.Normalize().IsNull() {
		return BlockExpression(*w.Else).Eval(vm, scope, clone)
	}
	return object.ExpressionResult{Trace: stacktrace.New(w.Pos), SignalVal: result.Normalize()}.Clone(clone)
}

func (f ForExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	iter := IntoEval(f.Expression).Eval(vm, scope, false)
	if iter.Signal.Has() {
		return iter.Clone(clone)
	}
	iter = IntoIter(iter.SignalVal.Normalize(), vm, f.Pos)
	if iter.Signal.Has() {
		return iter.Clone(clone)
	}
	arguments := object.Arguments{Last: f.Arguments.Last}
	for _, arg := range f.Arguments.Elements {
		res := IntoEval(arg.Default).Eval(vm, scope, true)
		if res.Signal.Has() {
			return res.Clone(clone)
		}
		arguments.Elements = append(arguments.Elements, object.Argument{
			Name:    arg.Name.Ident,
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
			args = append([]*object.Value{{InnerValue: object.Number{Value: float64(iteration)}}}, args...)
		}
		iteration++
		object.ParseArgs(arguments, args, scope)
	}
	scope = scope.Push()
	switch val := iter.SignalVal.Normalize().InnerValue.(type) {
	case object.Array:
		for _, element := range val.Elements {
			declareVals([]*object.Value{element})
			expr = BlockExpression(f.Block).Eval(vm, scope, false)
			switch {
			case expr.Signal == signal.SignalContinue:
				result = expr.SignalVal.Normalize()
			case expr.Signal == signal.SignalBreak:
				result = expr.SignalVal.Normalize()
				goto after
			case expr.Signal.Has():
				return expr.Clone(clone)
			default:
				result = expr.SignalVal.Normalize()
			}
		}
	case object.Object:
		for k, v := range val.Map {
			declareVals([]*object.Value{{InnerValue: object.String{Value: k}}, v})
			expr = BlockExpression(f.Block).Eval(vm, scope, false)
			switch {
			case expr.Signal == signal.SignalContinue:
				result = expr.SignalVal.Normalize()
			case expr.Signal == signal.SignalBreak:
				result = expr.SignalVal.Normalize()
				goto after
			case expr.Signal.Has():
				return expr.Clone(clone)
			default:
				result = expr.SignalVal.Normalize()
			}
		}
	case object.Method, *object.Subroutine:
		for {
			expr = Call(vm, iter.SignalVal.Normalize(), nil, false, f.Pos, nil)
			if err, ok := expr.SignalVal.Normalize().InnerValue.(object.Error); ok && expr.Signal == signal.SignalRaise && err.Code == errors.IteratorStop {
				break
			}
			if expr.Signal.Has() {
				return expr.Clone(clone)
			}
			var elements []*object.Value
			if val, ok := expr.SignalVal.Normalize().InnerValue.(object.Array); ok {
				elements = val.Elements
			} else {
				elements = []*object.Value{expr.SignalVal.Normalize()}
			}
			declareVals(elements)
			expr = BlockExpression(f.Block).Eval(vm, scope, false)
			switch {
			case expr.Signal == signal.SignalContinue:
				result = expr.SignalVal.Normalize()
			case expr.Signal == signal.SignalBreak:
				result = expr.SignalVal.Normalize()
				goto after
			case expr.Signal.Has():
				return expr.Clone(clone)
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
				goto after
			case expr.Signal.Has():
				return expr.Clone(clone)
			default:
				result = expr.SignalVal.Normalize()
			}
		}
	default:
		return object.ExpressionResult{Trace: stacktrace.New(f.Pos),
			Signal: signal.SignalRaise,
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("expected iterable, got %v", iter.SignalVal.Normalize().Type()),
			}},
		}
	}
	goto after
after:
	if f.Else != nil && result.Normalize().IsNull() {
		return BlockExpression(*f.Else).Eval(vm, scope, clone)
	}
	return object.ExpressionResult{Trace: stacktrace.New(f.Pos), SignalVal: result.Normalize()}.Clone(clone)
}

func (i IfExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	expr := IntoEval(i.Condition).Eval(vm, scope, false)
	if expr.Signal.Has() {
		return expr.Clone(clone)
	}
	ok, expr := IntoBool(vm, expr.SignalVal.Normalize(), i.Pos)
	if expr.Signal.Has() {
		return expr.Clone(clone)
	}
	if ok {
		return BlockExpression(i.If).Eval(vm, scope, clone)
	}
	for _, next := range i.Elifs {
		expr := IntoEval(next.Condition).Eval(vm, scope, false)
		if expr.Signal.Has() {
			return expr.Clone(clone)
		}
		ok, expr := IntoBool(vm, expr.SignalVal.Normalize(), next.Pos)
		if ok {
			return BlockExpression(next.Block).Eval(vm, scope, clone)
		}
	}
	if i.Else != nil {
		return BlockExpression(*i.Else).Eval(vm, scope, clone)
	}
	return object.ExpressionResult{Trace: stacktrace.New(i.Pos)}
}

func (p PrefixExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	expr := IntoEval(p.Value).Eval(vm, scope, false)
	if expr.Signal.Has() {
		return expr.Clone(clone)
	}
	typeRaise := object.ExpressionResult{Trace: stacktrace.New(p.Pos),
		Signal: signal.SignalRaise,
		SignalVal: &object.Value{
			InnerValue: object.Error{
				Code:    errors.TypeError,
				Message: fmt.Sprintf("unsupported type: %s for prefix operator '%s'", expr.SignalVal.Normalize().Type(), p.Operator.String()),
			},
		},
	}
	if p.Operator == operators.Ptr {
		return expr // No clone
	}
	if p.Operator == operators.Typeof {
		return object.ExpressionResult{Trace: stacktrace.New(p.Pos),
			SignalVal: &object.Value{
				InnerValue: object.String{
					Value: expr.SignalVal.Normalize().Type().String(),
				},
			},
		}
	}
	switch val := expr.SignalVal.Normalize().InnerValue.(type) {
	case object.Number:
		switch p.Operator {
		case operators.Neg:
			return object.ExpressionResult{Trace: stacktrace.New(p.Pos), SignalVal: &object.Value{InnerValue: object.Number{Value: -val.Value}}}
		case operators.Inc:
			expr.SignalVal.Normalize().InnerValue = object.Number{Value: val.Value + 1}
			return expr.Clone(clone)
		case operators.Dec:
			expr.SignalVal.Normalize().InnerValue = object.Number{Value: val.Value - 1}
			return expr.Clone(clone)
		case operators.Not:
			if float64(int(val.Value)) == val.Value {
				return object.ExpressionResult{Trace: stacktrace.New(p.Pos), SignalVal: &object.Value{InnerValue: object.Number{Value: float64(^int(val.Value))}}}
			}
			return typeRaise
		default:
			return typeRaise
		}
	case object.Bool:
		if p.Operator == operators.Not {
			return object.ExpressionResult{Trace: stacktrace.New(p.Pos), SignalVal: &object.Value{InnerValue: object.Bool{Value: !val.Value}}}
		}
		return typeRaise
	case object.Array:
		if p.Operator == operators.Dec {
			if len(val.Elements) == 0 {
				return object.ExpressionResult{Trace: stacktrace.New(p.Pos),
					Signal: signal.SignalRaise,
					SignalVal: &object.Value{
						InnerValue: object.Error{
							Code:    errors.ValueError,
							Message: "array is empty",
						},
					},
				}
			}
			last := val.Elements[len(val.Elements)-1]
			*expr.SignalVal.Normalize() = object.Value{InnerValue: object.Array{Elements: val.Elements[:len(val.Elements)-1]}}
			return object.ExpressionResult{Trace: stacktrace.New(p.Pos), SignalVal: last.Normalize()}.Clone(clone)
		}
		return typeRaise
	case object.String:
		if p.Operator == operators.Dec {
			runes := []rune(val.Value)
			if len(runes) == 0 {
				return object.ExpressionResult{Trace: stacktrace.New(p.Pos),
					Signal: signal.SignalRaise,
					SignalVal: &object.Value{
						InnerValue: object.Error{
							Code:    errors.ValueError,
							Message: "string is empty",
						},
					},
				}
			}
			last := runes[len(runes)-1]
			*expr.SignalVal.Normalize() = object.Value{InnerValue: object.String{Value: string(runes[:len(runes)-1])}}
			return object.ExpressionResult{Trace: stacktrace.New(p.Pos), SignalVal: &object.Value{InnerValue: object.String{Value: string(last)}}}
		}
		return typeRaise
	case object.Object:
		method := GetAttrMethod(expr.SignalVal.Normalize(), p.Operator.Method(), p.Pos)
		if method.Signal.Has() {
			return method.Clone(clone)
		}
		return Call(vm, method.SignalVal.Normalize(), nil, clone, p.Pos, nil)
	default:
		return typeRaise
	}
}

func (b BinExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	left := IntoEval(b.Left).Eval(vm, scope, false)
	if left.Signal.Has() {
		return left.Clone(clone)
	}

	if b.Operator == operators.Dot || b.Operator == operators.Method {
		if right, ok := b.Right.(ast.IdentExpression); ok {
			if b.Operator == operators.Dot {
				return GetAttr(left.SignalVal.Normalize(), right.Ident, b.Pos).Clone(clone)
			}
			return GetAttrMethod(left.SignalVal.Normalize(), right.Ident, b.Pos).Clone(clone)
		}
		return object.ExpressionResult{Trace: stacktrace.New(b.Pos),
			Signal: signal.SignalRaise,
			SignalVal: &object.Value{
				InnerValue: object.Error{
					Code:    errors.SyntaxError,
					Message: fmt.Sprintf("expected identifier after '%s'", b.Operator.String()),
				},
			},
		}

	}
	right := IntoEval(b.Right).Eval(vm, scope, false)
	if right.Signal.Has() {
		return right.Clone(clone)
	}
	return RunBin(vm, scope, clone, left.SignalVal.Normalize(), right.SignalVal.Normalize(), b.Operator, b.Pos)
}

func (o ObjectDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	obj := make(map[string]*object.Value)
	for _, val := range o.Vals {
		if val.IsContinuos {
			res := IntoEval(val.Value).Eval(vm, scope, true)
			if res.Signal.Has() {
				return res.Clone(clone)
			}
			inner, ok := res.SignalVal.InnerValue.(object.Object)
			if !ok {
				return object.ExpressionResult{Trace: stacktrace.New(val.Pos),
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
			return res.Clone(clone)
		}
		obj[val.Name.Ident] = res.SignalVal.Normalize()
	}
	return object.ExpressionResult{Trace: stacktrace.New(o.Pos),
		SignalVal: &object.Value{InnerValue: object.Object{
			Map: obj,
		}},
	}
}

func (s SubroutineDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	if s.Name != nil {
		sub := DeclareSubroutine(vm, scope, clone, s.Name.Ident, ast.SubroutineDec(s))
		if sub.Signal.Has() {
			return sub.Clone(clone)
		}
		scope.Set(s.Name.Ident, sub.SignalVal.Normalize())
		return object.ExpressionResult{Trace: stacktrace.New(s.Pos)}
	}
	return DeclareSubroutine(vm, scope, clone, "<anonymous>", ast.SubroutineDec(s))
}

func (n NumDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return object.ExpressionResult{Trace: stacktrace.New(n.Pos),
		SignalVal: &object.Value{InnerValue: object.Number{Value: n.Value}},
	}
}

func (b BoolDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return object.ExpressionResult{Trace: stacktrace.New(b.Pos),
		SignalVal: &object.Value{InnerValue: object.Bool{Value: b.Value}}}
}
func (s StringDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return object.ExpressionResult{Trace: stacktrace.New(s.Pos),
		SignalVal: &object.Value{InnerValue: object.String{Value: s.Value}},
	}
}
func (n NullDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return object.ExpressionResult{Trace: stacktrace.New(n.Pos),
		SignalVal: &object.Value{InnerValue: object.Null{}},
	}
}

func (a ArrayDec) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	var arr []*object.Value
	for _, val := range a.Elements {
		if val.IsContinuos {
			res := IntoEval(val.Value).Eval(vm, scope, true)
			if res.Signal.Has() {
				return res.Clone(clone)
			}
			inner, ok := res.SignalVal.InnerValue.(object.Array)
			if !ok {
				return object.ExpressionResult{Trace: stacktrace.New(val.Pos),
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
			return res.Clone(clone)
		}
		arr = append(arr, res.SignalVal.Normalize())
	}
	return object.ExpressionResult{Trace: stacktrace.New(a.Pos),
		SignalVal: &object.Value{InnerValue: object.Array{Elements: arr}},
	}
}

func (l LetExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	val := &object.Value{InnerValue: object.Null{}}
	scope.Set(l.Ident, val)
	return object.ExpressionResult{Trace: stacktrace.New(l.Pos),
		SignalVal: val,
	}
}

func (b BlockExpression) Eval(vm *vm.VM, scope scope.Scope[*object.Value], clone bool) object.ExpressionResult {
	return RunBlock(b, vm, scope, clone, nil)
}
