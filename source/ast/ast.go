package ast

import (
	"github.com/vandi37/aqua/source/operators"
	"github.com/vandi37/aqua/source/signal"
)

type Expression interface {
	expression()
}

type (
	Argument struct {
		Name    string
		Default Expression
	}
	Arguments struct {
		Elements []Argument
		// optional
		Last *string
	}
	ObjectDec struct {
		Vals []struct {
			Name        string
			Value       Expression
			IsContinuos bool
		}
	}
	IntDec    int
	NumDec    float64
	StringDec string
	NullDec   struct{}
	ArrayDec  struct {
		Elements []struct {
			Value       Expression
			IsContinuos bool
		}
	}

	BinExpression struct {
		Left     Expression
		Operator operators.Operator
		Right    Expression
	}
	PrefixExpression struct {
		Operator operators.PrefixOperator
		Value    Expression
	}

	LetExpression struct {
		Name string
	}

	BlockExpression struct {
		Expressions []Expression
		Catch       *CatchBlock
	}

	CatchBlock struct {
		Name        string
		Expressions BlockExpression
	}

	IfExpression struct {
		If      BlockExpression
		ElseIfs []BlockExpression
		Else    *BlockExpression
		Catch   *CatchBlock
	}
	ForExpression struct {
		Arguments  Arguments
		Expression Expression
		Block      BlockExpression
		Else       *BlockExpression
		Catch      *CatchBlock
	}
	WhileExpression struct {
		// a marker is a expression while or repeat-until
		IsWhile   bool
		Condition Expression
		// optional, may be nil
		After Expression
		Block BlockExpression
		Else  *BlockExpression
		Catch *CatchBlock
	}
	SubroutineDec struct {
		Arguments Arguments
		Body      BlockExpression
		// optional
		Prototype *ObjectDec
	}
	SignalExpression struct {
		Signal signal.Signal
		SigVal Expression
	}
)

func (ObjectDec) expression()        {}
func (IntDec) expression()           {}
func (NumDec) expression()           {}
func (NullDec) expression()          {}
func (StringDec) expression()        {}
func (ArrayDec) expression()         {}
func (BinExpression) expression()    {}
func (PrefixExpression) expression() {}
func (LetExpression) expression()    {}
func (BlockExpression) expression()  {}
func (IfExpression) expression()     {}
func (ForExpression) expression()    {}
func (WhileExpression) expression()  {}
func (SubroutineDec) expression()    {}
func (SignalExpression) expression() {}
