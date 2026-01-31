package ast

import (
	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/operators"
	"github.com/vandi37/aqua/source/signal"
)

type Expression interface {
	expression()
	// String() string
}

type (
	Arguments struct {
		Elements []struct {
			Name    IdentExpression
			Default Expression
		}
		Last *string
	}
	ObjectVal struct {
		Name        IdentExpression
		Value       Expression
		IsContinuos bool
		Pos         pos.Pos
	}
	ObjectDec struct {
		Vals []ObjectVal
		Pos  pos.Pos
	}
	NumDec struct {
		Value float64
		Pos   pos.Pos
	}
	BoolDec struct {
		Value bool
		Pos   pos.Pos
	}
	StringDec struct {
		Value string
		Pos   pos.Pos
	}
	NullDec struct {
		Pos pos.Pos
	}
	ErrorDec struct {
		Value errors.Error
		Pos   pos.Pos
	}
	ArrayElement struct {
		Value       Expression
		IsContinuos bool
		Pos         pos.Pos
	}
	ArrayDec struct {
		Pos      pos.Pos
		Elements []ArrayElement
	}
	SubroutineDec struct {
		Arguments Arguments
		Body      BlockExpression
		Prototype Expression
		Pos       pos.Pos
		Name      *IdentExpression
	}
	BinExpression struct {
		Left     Expression
		Operator operators.Operator
		Pos      pos.Pos
		Right    Expression
	}
	PrefixExpression struct {
		Operator operators.PrefixOperator
		Value    Expression
		Pos      pos.Pos
	}
	CallExpression struct {
		Subroutine Expression
		Args       []ArrayElement
		Pos        pos.Pos
	}

	LetExpression struct {
		IdentExpression
	}

	CatchBlock struct {
		Name        IdentExpression
		Expressions BlockExpression
		Pos         pos.Pos
	}
	BlockExpression struct {
		Pos         pos.Pos
		Expressions []Expression
		Catch       *CatchBlock
	}
	ElifBlock struct {
		Pos       pos.Pos
		Block     BlockExpression
		Condition Expression
	}
	IfExpression struct {
		Pos       pos.Pos
		If        BlockExpression
		Condition Expression
		Elifs     []ElifBlock
		Else      *BlockExpression
	}
	ForExpression struct {
		Pos        pos.Pos
		Arguments  Arguments
		Expression Expression
		IsEnum     bool
		Block      BlockExpression
		Else       *BlockExpression
	}
	WhileExpression struct {
		// a marker is a expression while or repeat-until
		Pos       pos.Pos
		IsWhile   bool
		Condition Expression
		After     Expression
		Block     BlockExpression
		Else      *BlockExpression
	}
	SignalExpression struct {
		Signal signal.Signal
		SigVal Expression
		Pos    pos.Pos
	}
	IdentExpression struct {
		Ident string
		Pos   pos.Pos
	}
	AssigmentExpression struct {
		Left     []Expression
		Right    []Expression
		Operator operators.Operator
		Pos      pos.Pos
	}
	ModExpression struct {
		Name   IdentExpression
		Pos    pos.Pos
		Body   BlockExpression
		Export []string
	}
	ImportExpression struct {
		Path Expression
		Name *IdentExpression
		Pos  pos.Pos
	}
)

func (ObjectDec) expression()           {}
func (NumDec) expression()              {}
func (NullDec) expression()             {}
func (BoolDec) expression()             {}
func (StringDec) expression()           {}
func (ErrorDec) expression()            {}
func (ArrayDec) expression()            {}
func (BinExpression) expression()       {}
func (PrefixExpression) expression()    {}
func (LetExpression) expression()       {}
func (CallExpression) expression()      {}
func (BlockExpression) expression()     {}
func (IfExpression) expression()        {}
func (ForExpression) expression()       {}
func (WhileExpression) expression()     {}
func (SubroutineDec) expression()       {}
func (SignalExpression) expression()    {}
func (IdentExpression) expression()     {}
func (AssigmentExpression) expression() {}
func (ModExpression) expression()       {}
func (ImportExpression) expression()    {}

// func (o ObjectDec) String() string {
// 	sb := strings.Builder{}
// 	sb.WriteString("ObjectDec(\n	Vals: [\n")
// 	for _, val := range o.Vals {
// 		if val.IsContinuos {
// 			fmt.Fprintf(&sb, "ObjectVal(\n	Value: ...%s\nPos: %s\n)\n",
// 				val.Value.String(),
// 				val.Pos.String(),
// 			)
// 		} else {
// 			fmt.Fprintf(&sb, "ObjectVal(\nIdent: %s\nValue: %s\nPos: %s\n)\n",
// 				val.Name.String(),
// 				val.Value.String(),
// 				val.Pos.String(),
// 			)
// 		}
// 	}
// 	fmt.Fprintf(&sb, "]\nPos: %s\n)", o.Pos.String())
// 	return sb.String()
// }
// func (n NumDec) String() string {
// 	return fmt.Sprintf("NumDec(\nValue: %v\nPos: %s\n)", n.Value, n.Pos.String())
// }
// func (n NullDec) String() string {
// 	return fmt.Sprintf("NullDec(Pos: %s)", n.Pos.String())
// }
// func (b BoolDec) String() string {
// 	return fmt.Sprintf("BoolDec(\nValue: %v\nPos: %s\n)", b.Value, b.Pos.String())
// }
// func (s StringDec) String() string {
// 	return fmt.Sprintf("StringDec(\nValue: %v\nPos: %s\n)", s.Value, s.Pos.String())
// }
// func (e ErrorDec) String() string {
// 	return fmt.Sprintf("ErrorDec(\nValue: %v\nPos: %s\n)", e.Value, e.Pos.String())
// }
// func (a ArrayElement) String() string {
// 	if a.IsContinuos {
// 		return fmt.Sprintf("ArrayElement(\nValue: ...%s\nPos: %s\n)",
// 			a.Value.String(),
// 			a.Pos.String(),
// 		)
// 	}
// 	return fmt.Sprintf("ArrayElement(\nValue: %s\nPos: %s\n)",
// 		a.Value.String(),
// 		a.Pos.String(),
// 	)

// }
// func (a ArrayDec) String() string {
// 	sb := strings.Builder{}
// 	sb.WriteString("ArrayDec(\nVals: [\n")
// 	for _, val := range a.Elements {
// 		fmt.Fprintf(&sb, "%s\n", val.String())
// 	}
// 	fmt.Fprintf(&sb, "]\nPos: %s\n)", a.Pos.String())
// 	return sb.String()
// }
// func (b BinExpression) String() string {
// 	return fmt.Sprintf("BinExpression(\nLeft: %s\nOperator: %s\nRight: %s\nPos: %s\n)",
// 		b.Left.String(), b.Operator.String(), b.Right.String(), b.Pos.String())
// }
// func (p PrefixExpression) String() string {
// 	return fmt.Sprintf("PrefixExpression(\nOperator: %s\nValue: %s\nPos: %s\n)",
// 		p.Operator.String(), p.Value.String(), p.Pos.String())
// }
// func (l LetExpression) String() string {
// 	return fmt.Sprintf("LetExpression(%s)", l.IdentExpression.String())
// }
// func (c CallExpression) String() string {
// 	sb := strings.Builder{}
// 	fmt.Fprintf(&sb, "CallExpression(\nSubroutine: %s\nElements: [", c.Subroutine.String())
// 	for _, element := range c.Args {
// 		fmt.Fprintf(&sb, "\n%s", element.String())
// 	}
// 	fmt.Fprintf(&sb, "]\nPos: %s\n)", c.Pos.String())
// 	return sb.String()

// }
// func (b BlockExpression) String() string {
// 	sb := strings.Builder{}
// 	sb.WriteString("BlockExpression(\nVals: [\n")
// 	for _, val := range b.Expressions {
// 		fmt.Fprintf(&sb, "%s\n", val.String())
// 	}
// 	sb.WriteString("]\n")
// 	if b.Catch != nil {
// 		fmt.Fprintf(&sb, "Catch: CatchBlock(\nName: %s\nExpressions: %s\nPos: %s\n)\n", b.Catch.Name.String(), b.Catch.Expressions.String(), b.Catch.Pos.String())
// 	}
// 	fmt.Fprintf(&sb, "Pos: %s\n)", b.Pos.String())
// 	return sb.String()

// }
// func (i IfExpression) String() string {
// 	sb := strings.Builder{}
// 	fmt.Fprintf(&sb, "IfExpression(\nCondition: %s\nIf: %s\nElifs: [\n",
// 		i.Condition.String(), i.If.String())
// 	for _, val := range i.Elifs {
// 		fmt.Fprintf(&sb, "ElifBlock(\nCondition: %s\nBlock: %s\nPos: %s\n)\n", val.Condition.String(), val.Block.String(), val.Pos.String())
// 	}
// 	sb.WriteString("]\n")
// 	if i.Else != nil {
// 		fmt.Fprintf(&sb, "Else: %s\n", i.Else.String())
// 	}
// 	fmt.Fprintf(&sb, "Pos: %s\n)", i.Pos.String())
// 	return sb.String()
// }
// func (f ForExpression) String() string {
// 	sb := strings.Builder{}
// 	fmt.Fprintf(&sb, "ForExpression(\nArguments: %s\nExpression: %s\nIsEnum: %v\nBlock: %s\n",
// 		f.Arguments.String(), f.Expression.String(), f.IsEnum, f.Block.String())
// 	if f.Else != nil {
// 		fmt.Fprintf(&sb, "Else: %s\n", f.Else.String())
// 	}
// 	fmt.Fprintf(&sb, "Pos: %s\n)", f.Pos.String())
// 	return sb.String()
// }
// func (w WhileExpression) String() string {
// 	sb := strings.Builder{}
// 	fmt.Fprintf(&sb, "WhileExpression(\nCondition: %s\nAfter: %s\nIsWhile: %v\nBlock: %s\n",
// 		w.Condition.String(), w.After.String(), w.IsWhile, w.Block.String())
// 	if w.Else != nil {
// 		fmt.Fprintf(&sb, "Else: %s\n", w.Else.String())
// 	}
// 	fmt.Fprintf(&sb, "Pos: %s\n)", w.Pos.String())
// 	return sb.String()
// }
// func (s SubroutineDec) String() string {
// 	sb := strings.Builder{}
// 	fmt.Fprintf(&sb, "SubroutineDec(\nArguments: %s\nBody: %s\n",
// 		s.Arguments.String(), s.Body.String())
// 	if s.Prototype != nil {
// 		fmt.Fprintf(&sb, "Prototype: %s\n", s.Prototype.String())
// 	}
// 	if s.Name != nil {
// 		fmt.Fprintf(&sb, "Name: %s\n", s.Name.String())
// 	}
// 	fmt.Fprintf(&sb, "Pos: %s\n)", s.Pos.String())
// 	return sb.String()
// }
// func (s SignalExpression) String() string {
// 	return fmt.Sprintf("SignalExpression(\nSignal: %s\nSigVal: %s\nPos: %s\n)",
// 		s.Signal.String(), s.SigVal.String(), s.Pos.String())
// }
// func (i IdentExpression) String() string {
// 	return fmt.Sprintf("IdentExpression(\nIdent: %s\nPos: %s\n)", i.Ident, i.Pos)
// }
// func (a AssigmentExpression) String() string {
// 	sb := strings.Builder{}
// 	sb.WriteString("AssigmentExpression(\nLeft: [\n")
// 	for _, v := range a.Left {
// 		fmt.Fprintf(&sb, "%s\n", v.String())
// 	}
// 	fmt.Fprintf(&sb, "]\nOperator: %s\nRight: [\n", a.Operator.String())
// 	for _, v := range a.Left {
// 		fmt.Fprintf(&sb, "%s\n", v.String())
// 	}
// 	fmt.Fprintf(&sb, "Pos: %s\n)", a.Pos.String())
// 	return sb.String()
// }
// func (m ModExpression) String() string {
// 	sb := strings.Builder{}
// 	sb.WriteString("ModExpression(\nName: %s\nBlock: %s\nExport: [\n")
// 	for _, v := range m.Export {
// 		fmt.Fprintf(&sb, "%s\n", v)
// 	}
// 	fmt.Fprintf(&sb, "]\nPos: %s\n)", m.Pos.String())
// 	return sb.String()
// }
// func (i ImportExpression) String() string {
// 	sb := strings.Builder{}
// 	fmt.Fprintf(&sb, "ImportExpression(\nPath: %s\n", i.Path.String())
// 	if i.Name != nil {
// 		fmt.Fprintf(&sb, "Name: %s\n", i.Name.String())
// 	}
// 	fmt.Fprintf(&sb, "Pos: %s\n)", i.Pos.String())
// 	return sb.String()
// }

// func (a Arguments) String() string {
// 	sb := strings.Builder{}
// 	fmt.Fprintf(&sb, "Arguments(\nElements: [\n")
// 	for _, val := range a.Elements {
// 		fmt.Fprintf(&sb, "Argument(\nName: %s\n", val.Name.String())
// 		if val.Default != nil {
// 			fmt.Fprintf(&sb, "Default: %s\n", val.Default.String())
// 		}
// 		sb.WriteString(")\n")
// 	}
// 	sb.WriteString("]\n")
// 	if a.Last != nil {
// 		fmt.Fprintf(&sb, "Last: %s\n", *a.Last)
// 	}
// 	fmt.Fprintf(&sb, ")")
// 	return sb.String()
// }
