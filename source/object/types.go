package object

type Type byte

const (
	TypeNull Type = iota
	TypeObject
	TypeSubroutine
	TypeArray
	TypeInt
	TypeNumber
	TypeString
	TypeBool
	TypeError
)

func (t Type) String() string {
	switch t {
	case TypeObject:
		return "object"
	case TypeSubroutine:
		return "subroutine"
	case TypeArray:
		return "array"
	case TypeInt:
		return "integer"
	case TypeNumber:
		return "number"
	case TypeString:
		return "string"
	case TypeNull:
		return "null"
	case TypeBool:
		return "boolean"
	case TypeError:
		return "error"
	default:
		return "unknown"
	}
}

func (Type) value() {}
