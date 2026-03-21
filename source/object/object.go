package object

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/aqua-aq/aqua-core/pkg/errors"
	"github.com/aqua-aq/aqua-core/pkg/scope"
	"github.com/aqua-aq/aqua-core/source/ast"
	"github.com/aqua-aq/aqua-core/source/vm"
	"github.com/google/uuid"
)

type InnerValue interface {
	value()
	Type() Type
	Equals(value *Value) bool
}

type Value struct {
	innerValue InnerValue
	uuid       uuid.UUID
}

func New(i InnerValue) *Value {
	return &Value{
		innerValue: i,
		uuid:       uuid.New(),
	}
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
		Last     *string
	}
	Subroutine struct {
		Name      string
		Arguments Arguments
		Scope     scope.Scope[string, *Value]
		Prototype *Value
		BuiltIn   func(*vm.VM[*Value], scope.Scope[string, *Value]) SubroutineResult
		Code      ast.BlockExpression
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

func (Object) value()     {}
func (Object) Type() Type { return TypeObject }
func (o Object) Equals(value *Value) bool {
	if obj, ok := value.Normalize().innerValue.(Object); ok {
		return maps.Equal(o.Map, obj.Map)
	}
	return false
}
func (Null) value()                   {}
func (Null) Type() Type               { return TypeNull }
func (Null) Equals(value *Value) bool { return value.Type() == TypeNull }
func (Error) value()                  {}
func (Error) Type() Type              { return TypeError }
func (e Error) Equals(value *Value) bool {
	if err, ok := value.Normalize().innerValue.(Error); ok {
		return e == err
	}
	return false
}
func (*Subroutine) value()     {}
func (*Subroutine) Type() Type { return TypeSubroutine }
func (s *Subroutine) Equals(value *Value) bool {
	if sub, ok := value.Normalize().innerValue.(*Subroutine); ok {
		return s == sub
	}
	if method, ok := value.Normalize().innerValue.(Method); ok {
		return s == method.Subroutine
	}
	return false
}
func (Method) value()     {}
func (Method) Type() Type { return TypeSubroutine }
func (m Method) Equals(value *Value) bool {
	if sub, ok := value.Normalize().innerValue.(*Subroutine); ok {
		return m.Subroutine == sub
	}
	if method, ok := value.Normalize().innerValue.(Method); ok {
		return m.Subroutine == method.Subroutine
	}
	return false
}
func (Number) value()     {}
func (Number) Type() Type { return TypeNumber }
func (n Number) Equals(value *Value) bool {
	if num, ok := value.Normalize().innerValue.(Number); ok {
		return n == num
	}
	return false
}
func (String) value()     {}
func (String) Type() Type { return TypeString }
func (s String) Equals(value *Value) bool {
	if str, ok := value.Normalize().innerValue.(String); ok {
		return s == str
	}
	return false
}
func (Bool) value()     {}
func (Bool) Type() Type { return TypeBool }
func (b Bool) Equals(value *Value) bool {
	if boolean, ok := value.Normalize().innerValue.(Bool); ok {
		return b == boolean
	}
	return false
}
func (Array) value()     {}
func (Array) Type() Type { return TypeArray }
func (a Array) Equals(value *Value) bool {
	if arr, ok := value.Normalize().innerValue.(Array); ok {
		return slices.Equal(a.Elements, arr.Elements)
	}
	return false
}
func (a Arguments) String() string {
	return a.stringify(scope.New[uuid.UUID, struct{}]())
}
func (a Arguments) stringify(visited scope.Scope[uuid.UUID, struct{}]) string {
	sb := strings.Builder{}
	sb.WriteByte('(')
	for i, arg := range a.Elements {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.Name)
		if !arg.Default.Normalize().IsNull() {
			sb.WriteString("=" + arg.Default.stringify(visited))
		}
	}
	if a.Last != nil {
		if len(a.Elements) != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("..." + *a.Last)
	}
	sb.WriteByte(')')
	return sb.String()
}
func (v *Value) IsNull() bool {
	_, ok := v.innerValue.(Null)
	return ok
}

func (v *Value) Normalize() *Value {
	if v == nil {
		return New(Null{})
	}
	if v.innerValue == nil {
		v.innerValue = Null{}
	}
	return v
}
func (v *Value) String() string { return v.Normalize().stringify(scope.New[uuid.UUID, struct{}]()) }
func (v *Value) stringify(visited scope.Scope[uuid.UUID, struct{}]) string {
	visited = visited.Push()
	if visited.Has(v.uuid) {
		return "<cycle>"
	}
	visited.Set(v.uuid, struct{}{})
	switch inner := v.innerValue.(type) {
	case String:
		return inner.Value
	case Number:
		return fmt.Sprint(inner.Value)
	case Bool:
		return fmt.Sprint(inner.Value)
	case Error:
		return errors.Error(inner).Error()
	case *Subroutine:
		return fmt.Sprintf("%s %s", inner.Name, inner.Arguments.stringify(visited))
	case Method:
		return fmt.Sprintf("%s %s", inner.Subroutine.Name, inner.Subroutine.Arguments.stringify(visited))
	case Array:
		var b strings.Builder
		b.WriteString("[")
		for i, v := range inner.Elements {
			b.WriteString(v.Normalize().stringify(visited))
			if i < len(inner.Elements)-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString("]")
		return b.String()
	case Object:
		var b strings.Builder
		b.WriteString("{")

		keys := make([]string, 0, len(inner.Map))
		for k := range inner.Map {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		for i, k := range keys {
			fmt.Fprintf(&b, "%s: %s", k, inner.Map[k].Normalize().stringify(visited))
			if i < len(keys)-1 {
				b.WriteString(", ")
			}
		}
		b.WriteString("}")
		return b.String()
	}
	return "null"
}
func (v *Value) Equals(value *Value) bool { return v.Normalize().innerValue.Equals(value.Normalize()) }
func (v *Value) Type() Type               { return v.Normalize().innerValue.Type() }
func IntoValue(err error) *Value {
	if e, ok := err.(errors.Error); ok {
		return &Value{innerValue: Error(e)}
	}
	return &Value{innerValue: String{Value: err.Error()}}
}

func (v *Value) Clone() *Value {
	return v.Normalize().deepClone(map[uuid.UUID]*Value{})
}

func (v *Value) deepClone(visited map[uuid.UUID]*Value) *Value {
	if cloned, ok := visited[v.uuid]; ok {
		return cloned
	}
	cloned := New(Null{})
	visited[v.uuid] = cloned

	switch inner := v.innerValue.(type) {
	case Object:
		object := make(map[string]*Value, len(inner.Map))
		for k, v := range inner.Map {
			object[k] = v.Normalize().deepClone(visited)
		}
		cloned.innerValue = Object{object}
	case Array:
		array := make([]*Value, len(inner.Elements))
		for i, e := range array {
			array[i] = e.Normalize().deepClone(visited)
		}
		cloned.innerValue = Array{array}
	default:
		cloned.innerValue = inner
	}
	return cloned
}
func (v *Value) Uuid() uuid.UUID        { return v.uuid }
func (v *Value) InnerValue() InnerValue { return v.innerValue }
func (v *Value) Set(i InnerValue)       { v.innerValue = i }
