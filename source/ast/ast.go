package ast

import (
	"github.com/vandi37/aqua/pkg/pos"
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
			Pos     pos.Pos
		}
		Last *string
	}
	ObjectDec struct {
		Vals []struct {
			Name        IdentExpression
			Value       Expression
			IsContinuos bool
			Pos         pos.Pos
		}
		Pos pos.Pos
	}
	NumDec struct {
		Value float64
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
	ArrayDec struct {
		Pos      pos.Pos
		Elements []struct {
			Value       Expression
			IsContinuos bool
			Pos         pos.Pos
		}
	}
	SubroutineDec struct {
		Arguments Arguments
		Body      BlockExpression
		Prototype Expression
		Pos       pos.Pos
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
		Args       []Expression
		Pos        pos.Pos
	}

	LetExpression struct {
		IdentExpression
	}

	BlockExpression struct {
		Pos         pos.Pos
		Expressions []Expression
		Catch       *struct {
			Name        IdentExpression
			Expressions BlockExpression
			Pos         pos.Pos
		}
	}
	IfExpression struct {
		Pos       pos.Pos
		If        BlockExpression
		Condition Expression
		ElseIfs   []struct {
			Pos       pos.Pos
			Block     BlockExpression
			Condition Expression
		}
		Else *BlockExpression
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
	GlobalSubroutineDec struct {
		SubroutineDec SubroutineDec
		Name          IdentExpression
	}
	SignalExpression struct {
		Signal signal.Signal
		SigVal Expression
		Pos    pos.Pos
	}
	IdentExpression struct {
		Ident string
		HasAt bool
		Pos   pos.Pos
	}
	AssigmentExpression struct {
		Left     []Expression
		Right    []Expression
		Operator operators.Operator
		Pos      pos.Pos
	}
	ModExpression struct {
		Name  IdentExpression
		Pos   pos.Pos
		Body BlockExpression
	}
	ImportExpression struct {
		Path Expression
		Pos  pos.Pos
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
func (ModExpression) expression()       {}
func (ImportExpression) expression()    {}
