package ast

import (
	"github.com/aqua-aq/aqua/pkg/pos"
	"github.com/aqua-aq/aqua/source/errors"
	"github.com/aqua-aq/aqua/source/object/signal"
	"github.com/aqua-aq/aqua/source/operators"
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
		IsGlobal  bool
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
		Pos       pos.Pos
		IsWhile   bool
		Condition Expression
		Block     BlockExpression
		Else      *BlockExpression
	}
	UsingExpression struct {
		Expression Expression
		Name       *IdentExpression
		Block      BlockExpression
		Pos        pos.Pos
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
	AssigmentPattern struct {
		Expression Expression
		Name       *IdentExpression
		Pos        pos.Pos
	}
	AssigmentExpression struct {
		ExpressionLeft *struct {
			Expression
			operators.Operator
		}
		Left  []AssigmentPattern
		Right Expression
		Pos   pos.Pos
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
	Case struct {
		Expression Expression
		Block      BlockExpression
		Pos        pos.Pos
	}
	SwitchExpression struct {
		Value   Expression
		Cases   []Case
		Default *BlockExpression
		Pos     pos.Pos
	}
	SliceExpression struct {
		Left  Expression
		Start Expression
		End   Expression
		Pos   pos.Pos
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
func (UsingExpression) expression()     {}
func (SubroutineDec) expression()       {}
func (SignalExpression) expression()    {}
func (IdentExpression) expression()     {}
func (AssigmentExpression) expression() {}
func (ModExpression) expression()       {}
func (ImportExpression) expression()    {}
func (SwitchExpression) expression()    {}
func (SliceExpression) expression()     {}
