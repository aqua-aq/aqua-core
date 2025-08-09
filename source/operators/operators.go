package operators

type Operator uint32

const (
	Assign Operator = 1 << iota
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
	Call
	Comma
	Method
)

func (o Operator) String() string {
	switch o {
	case Assign:
		return "="
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Multiply:
		return "*"
	case Divide:
		return "/"
	case Modulo:
		return "%"
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
	case Call:
		return "()"
	case Dot:
		return "."
	case Comma:
		return ","
	case Method:
		return "::"
	case Assign | Plus:
		return "+="
	case Assign | Minus:
		return "-="
	case Assign | Multiply:
		return "*="
	case Assign | Divide:
		return "/="
	case Assign | Modulo:
		return "%="
	case Assign | StrongDivide:
		return "//="
	case Assign | And:
		return "and="
	case Assign | Or:
		return "or="
	case Assign | Xor:
		return "xor="
	case Assign | Shr:
		return ">>="
	case Assign | Shl:
		return "<<="
	case Assign | In:
		return "in="
	default:
		return "unknown"
	}
}

type PrefixOperator byte

const (
	Neg PrefixOperator = iota
	Not
	Inc
	Dec
)

func (o PrefixOperator) String() string {
	switch o {
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
