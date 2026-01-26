package ast

import (
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/operators"
	"github.com/vandi37/aqua/source/signal"
)

type Expression interface {
	expression()
}

type (
	Arguments struct {
		Elements []struct {
			Name    IdentExpression
			Default Expression
		}
		Last *string
	}
	ObjectDec struct {
		Vals []struct {
			Name        IdentExpression
			Value       Expression
			IsContinuos bool
		}
	}
	NumDec    float64
	StringDec string
	NullDec   struct{}
	ErrorDec  errors.Error
	ArrayDec  struct {
		Elements []struct {
			Value       Expression
			IsContinuos bool
		}
	}
	SubroutineDec struct {
		Arguments Arguments
		Body      BlockExpression
		Prototype Expression
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
	CallExpression struct {
		Subroutine Expression
		Args       []Expression
	}

	LetExpression struct {
		IdentExpression
	}

	BlockExpression struct {
		Expressions []Expression
		Catch       *struct {
			Name        IdentExpression
			Expressions BlockExpression
		}
	}
	IfExpression struct {
		If        BlockExpression
		Condition Expression
		ElseIfs   []struct {
			Block     BlockExpression
			Condition Expression
		}
		Else *BlockExpression
	}
	ForExpression struct {
		Arguments  Arguments
		Expression Expression
		IsEnum     bool
		Block      BlockExpression
		Else       *BlockExpression
	}
	WhileExpression struct {
		// a marker is a expression while or repeat-until
		IsWhile   bool
		Condition Expression
		After     Expression
		Block     BlockExpression
		Else      *BlockExpression
	}
	GlobalSubroutineDec struct {
		SubroutineDec SubroutineDec
		Name          IdentExpression
	}
	SignalExpression struct {
		Signal signal.Signal
		SigVal Expression
	}
	IdentExpression struct {
		Ident string
		HasAt bool
	}
	AssigmentExpression struct {
		Left     []Expression
		Right    []Expression
		Operator operators.Operator
	}
)

func (ObjectDec) expression()           {}
func (NumDec) expression()              {}
func (NullDec) expression()             {}
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
func (GlobalSubroutineDec) expression() {}
func (SubroutineDec) expression()       {}
func (SignalExpression) expression()    {}
func (IdentExpression) expression()     {}
func (AssigmentExpression) expression() {}
