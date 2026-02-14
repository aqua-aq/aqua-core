package operators

import (
	"github.com/aqua-aq/aqua-core/source/keywords"
	"github.com/aqua-aq/aqua-core/source/power"
)

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

	Question
	In
	Index
	Bind

	Dot
	QuestionDot
	Method
	QuestionMethod
	Delete
	QuestionDelete
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
	case Question:
		return "??"
	case In:
		return "in"
	case Index:
		return "[]"
	case Dot:
		return "."
	case QuestionDot:
		return "?."
	case Method:
		return ".>"
	case QuestionMethod:
		return "?.>"
	case Delete:
		return ".~"
	case QuestionDelete:
		return "?.~"
	case Bind:
		return "->"
	default:
		return "unknown"
	}
}

func (o Operator) Method() (string, bool) {
	switch o {
	case Plus:
		return keywords.Add, true
	case Minus:
		return keywords.Sub, true
	case Multiply:
		return keywords.Mul, true
	case Divide:
		return keywords.Div, true
	case Modulo:
		return keywords.Mod, true
	case StrongDivide:
		return keywords.IDiv, true
	case And:
		return keywords.Add, true
	case Or:
		return keywords.Or, true
	case Xor:
		return keywords.Xor, true
	case Shr:
		return keywords.Shr, true
	case Shl:
		return keywords.Shl, true
	case Equal:
		return keywords.Eq, true
	case NotEqual:
		return keywords.Ne, true
	case Greater:
		return keywords.Gt, true
	case Less:
		return keywords.Lt, true
	case GreaterEqual:
		return keywords.Ge, true
	case LessEqual:
		return keywords.Le, true
	case In:
		return keywords.In, true
	case Index:
		return keywords.Index, true
	case Bind:
		return keywords.Bind, true
	default:
		return "", false
	}
}

func (o Operator) IsValidInAssign() bool {
	switch o {
	case None, Plus, Minus, Multiply, Divide, Modulo, StrongDivide, And, Or, Xor, Shr, Shl, Question, In, Bind:
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
	Typeof
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
	case Typeof:
		return "typeof"
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

func (o Operator) Power() power.BindingPower {
	switch o {
	case Question:
		return power.PowerQuestion
	case Or:
		return power.PowerOr
	case And:
		return power.PowerAnd
	case Xor:
		return power.PowerXor
	case Plus, Minus:
		return power.PowerAdditive
	case Multiply, Divide, Modulo, StrongDivide:
		return power.PowerMultiplicative
	case Equal, NotEqual:
		return power.PowerEquality
	case Greater, Less, GreaterEqual, LessEqual, In:
		return power.PowerComparison
	case Shr, Shl:
		return power.PowerShift
	case Bind:
		return power.PowerBind
	case Dot, Method, QuestionDot, QuestionMethod, Delete, QuestionDelete, Index:
		return power.PowerPostfix
	default:
		return power.PowerLowest
	}
}

func (o Operator) IsRight() bool {
	return o == Question || o == Bind
}
