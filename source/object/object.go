package object

import (
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/source/ast"
	"github.com/vandi37/aqua/source/errors"
)

type InnerValue interface {
	value()
	Clone() *Value
	Type() Type
}

type Value struct {
	InnerValue
}

type (
	Object struct {
		// an object should always be initialized
		Map         map[string]*Value
		Constructor *Subroutine
	}
	Argument struct {
		Name    string
		Default *Value
	}
	Arguments struct {
		Elements []Argument
		// optional
		Last *string
	}
	Subroutine struct {
		Arguments Arguments
		Scope     scope.Scope[*Value]
		// use empty map for no prototype
		Prototype map[string]*Value
		// optional
		BuildIn func(scope.Scope[*Value]) SubroutineResult
		// optional code
		Code ast.BlockExpression
	}
	Method struct {
		*Subroutine
		It *Value
	}

	Null   struct{}
	Error  errors.Error
	Int    struct{ Value int }
	Number struct{ Value float64 }
	String struct{ Value string }
	Bool   struct{ Value bool }
	Array  struct{ Elements []*Value }
)

func (Object) value() {}
func (o Object) Clone() *Value {
	oMap := make(map[string]*Value, len(o.Map))
	for k, v := range o.Map {
		oMap[k] = v.Clone()
	}
	return &Value{Object{oMap, o.Constructor}}
}
func (Object) Type() Type  { return TypeObject }
func (Null) value()        {}
func (Null) Clone() *Value { return &Value{Null{}} }
func (Null) Type() Type    { return TypeNull }
func (Error) value()       {}
func (e Error) Clone() *Value {
	return &Value{e}
}
func (Error) Type() Type   { return TypeError }
func (*Subroutine) value() {}
func (s *Subroutine) Clone() *Value {
	return &Value{s}
}
func (*Subroutine) Type() Type { return TypeSubroutine }
func (m Method) Clone() *Value {
	return &Value{m}
}
func (Int) value() {}
func (i Int) Clone() *Value {
	return &Value{i}
}
func (Int) Type() Type { return TypeInt }
func (Number) value()  {}
func (n Number) Clone() *Value {
	return &Value{n}
}
func (Number) Type() Type { return TypeNumber }
func (String) value()     {}
func (s String) Clone() *Value {
	return &Value{s}
}
func (String) Type() Type { return TypeString }
func (Bool) value()       {}
func (b Bool) Clone() *Value {
	return &Value{b}
}
func (Bool) Type() Type { return TypeBool }
func (Array) value()    {}
func (a Array) Clone() *Value {
	clone := make([]*Value, 0, len(a.Elements))
	for _, v := range a.Elements {
		clone = append(clone, v.Clone())
	}
	return &Value{Array{clone}}
}
func (Array) Type() Type { return TypeArray }

func (v *Value) IsNull() bool {
	_, ok := v.InnerValue.(Null)
	return ok
}

func (v *Value) Normalize() *Value {
	if v == nil {
		return &Value{Null{}}
	}
	if v.InnerValue == nil {
		v.InnerValue = Null{}
	}
	return v
}
