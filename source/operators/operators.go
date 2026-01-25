package operators

import "github.com/vandi37/aqua/source/keywords"

type Operator byte

const (
	None Operator = iota
	Plus
	Minus
	Multiply
	Divide
	Modulo
	StrongDivide

	And
	Or
	Xor
	Shr
	Shl

	Equal
	NotEqual
	Greater
	Less
	GreaterEqual
	LessEqual

	In
	Index
	Dot
	Method
	Bind
)

func (o Operator) String() string {
	switch o {
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Multiply:
		return "*"
	case Divide:
		return "/"
	case Modulo:
		return "%%"
	case StrongDivide:
		return "//"
	case And:
		return "and"
	case Or:
		return "or"
	case Xor:
		return "xor"
	case Shr:
		return ">>"
	case Shl:
		return "<<"
	case Equal:
		return "=="
	case NotEqual:
		return "~="
	case Greater:
		return ">"
	case Less:
		return "<"
	case GreaterEqual:
		return ">="
	case LessEqual:
		return "<="
	case In:
		return "in"
	case Index:
		return "[]"
	case Dot:
		return "."
	case Method:
		return ".>"
	case Bind:
		return "->"
	default:
		return "unknown"
	}
}

func (o Operator) Method() string {
	switch o {
	case Plus:
		return keywords.Add
	case Minus:
		return keywords.Sub
	case Multiply:
		return keywords.Mul
	case Divide:
		return keywords.Div
	case Modulo:
		return keywords.Mod
	case StrongDivide:
		return keywords.IDiv
	case And:
		return keywords.Add
	case Or:
		return keywords.Or
	case Xor:
		return keywords.Xor
	case Shr:
		return keywords.Shr
	case Shl:
		return keywords.Shl
	case Equal:
		return keywords.Eq
	case NotEqual:
		return keywords.Ne
	case Greater:
		return keywords.Gt
	case Less:
		return keywords.Lt
	case GreaterEqual:
		return keywords.Ge
	case LessEqual:
		return keywords.Le
	case In:
		return keywords.In
	case Index:
		return keywords.Index
	case Bind:
		return keywords.Bind
	default:
		return ""
	}
}

func (o Operator) IsValidInAssign() bool {
	switch o {
	case None, Plus, Minus, Multiply, Divide, Modulo, StrongDivide, And, Or, Xor, Shr, Shl, In, Bind:
		return true
	default:
		return false
	}
}

type PrefixOperator byte

const (
	Ptr PrefixOperator = iota
	Neg
	Not
	Inc
	Dec
)

func (o PrefixOperator) String() string {
	switch o {
	case Ptr:
		return "&"
	case Neg:
		return "-"
	case Not:
		return "not"
	case Inc:
		return "++"
	case Dec:
		return "--"
	default:
		return "unknown"
	}
}
func (o PrefixOperator) Method() string {
	switch o {
	case Neg:
		return keywords.Neg
	case Not:
		return keywords.Not
	case Inc:
		return keywords.Inc
	case Dec:
		return keywords.Dec
	default:
		return ""
	}
}
