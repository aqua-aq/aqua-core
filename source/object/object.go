package object

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/source/ast"
	"github.com/vandi37/aqua/source/errors"
)

type InnerValue interface {
	value()
	Clone() *Value
	Type() Type
	Equals(value *Value) bool
	fmt.Stringer
}

type Value struct {
	InnerValue
}

type (
	Object struct {
		// an object should always be initialized
		Map map[string]*Value
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
		Name      string
		Arguments Arguments
		Scope     scope.Scope[*Value]
		Prototype *Value
		// optional
		BuildIn func(scope.Scope[*Value]) SubroutineResult
		// optional code
		Code ast.BlockExpression
	}
	Method struct {
		Subroutine *Subroutine
		It         *Value
	}

	Null   struct{}
	Error  errors.Error
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
	return &Value{Object{oMap}}
}
func (Object) Type() Type { return TypeObject }
func (o Object) String() string {
	return fmt.Sprintf("%v", o.Map)
}
func (o Object) Equals(value *Value) bool {
	if obj, ok := value.Normalize().InnerValue.(Object); ok {
		return maps.Equal(o.Map, obj.Map)
	}
	return false
}
func (Null) value()         {}
func (Null) Clone() *Value  { return &Value{Null{}} }
func (Null) Type() Type     { return TypeNull }
func (Null) String() string { return "null" }
func (Null) Equals(value *Value) bool {
	return value.Type() == TypeNull
}
func (Error) value() {}
func (e Error) Clone() *Value {
	return &Value{e}
}
func (Error) Type() Type { return TypeError }
func (e Error) Equals(value *Value) bool {
	if err, ok := value.Normalize().InnerValue.(Error); ok {
		return e == err
	}
	return false
}
func (e Error) String() string {
	return fmt.Sprintf("%s", errors.Error(e).Error())
}
func (*Subroutine) value() {}
func (s *Subroutine) Clone() *Value {
	return &Value{s}
}
func (*Subroutine) Type() Type { return TypeSubroutine }
func (s *Subroutine) String() string {
	return fmt.Sprintf("%v", s.Arguments)
}
func (s *Subroutine) Equals(value *Value) bool {
	if sub, ok := value.Normalize().InnerValue.(*Subroutine); ok {
		return s == sub
	}
	return false
}
func (Method) value() {}
func (m Method) Clone() *Value {
	return &Value{m}
}
func (Method) Type() Type       { return TypeSubroutine }
func (m Method) String() string { return m.Subroutine.String() }
func (m Method) Equals(value *Value) bool {
	if method, ok := value.Normalize().InnerValue.(Method); ok {
		return m.Subroutine == method.Subroutine
	}
	return false
}
func (Number) value() {}
func (n Number) Clone() *Value {
	return &Value{n}
}
func (Number) Type() Type       { return TypeNumber }
func (n Number) String() string { return fmt.Sprintf("%v", n.Value) }
func (n Number) Equals(value *Value) bool {
	if num, ok := value.Normalize().InnerValue.(Number); ok {
		return n == num
	}
	return false
}
func (String) value() {}
func (s String) Clone() *Value {
	return &Value{s}
}
func (String) Type() Type       { return TypeString }
func (s String) String() string { return fmt.Sprintf("\"%v\"", s.Value) }
func (s String) Equals(value *Value) bool {
	if str, ok := value.Normalize().InnerValue.(String); ok {
		return s == str
	}
	return false
}
func (Bool) value() {}
func (b Bool) Clone() *Value {
	return &Value{b}
}
func (Bool) Type() Type       { return TypeBool }
func (b Bool) String() string { return fmt.Sprintf("%v", b.Value) }
func (b Bool) Equals(value *Value) bool {
	if boolean, ok := value.Normalize().InnerValue.(Bool); ok {
		return b == boolean
	}
	return false
}
func (Array) value() {}
func (a Array) Clone() *Value {
	clone := make([]*Value, 0, len(a.Elements))
	for _, v := range a.Elements {
		clone = append(clone, v.Clone())
	}
	return &Value{Array{clone}}
}
func (Array) Type() Type { return TypeArray }
func (a Array) String() string {
	return fmt.Sprintf("%v", a.Elements)
}
func (a Array) Equals(value *Value) bool {
	if arr, ok := value.Normalize().InnerValue.(Array); ok {
		return slices.Equal(a.Elements, arr.Elements)
	}
	return false
}
func (a Arguments) String() string {
	sb := strings.Builder{}
	sb.WriteByte('(')
	for i, arg := range a.Elements {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", arg))
		if arg.Default == nil || !arg.Default.IsNull() {
			sb.WriteString(fmt.Sprintf("= %v", arg.Default.String()))
		}
	}
	if a.Last != nil {
		sb.WriteString(fmt.Sprintf(", ...%v", *a.Last))
	}
	sb.WriteByte(')')
	return sb.String()
}
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

func (v *Value) String() string {
	return v.Normalize().InnerValue.String()
}
