package parser

import (
	"math"

	"github.com/vandi37/aqua/source/ast"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/lexer/tokens"
	"github.com/vandi37/aqua/source/power"
	"github.com/vandi37/aqua/source/signal"
)

func (p *Parser) PrimaryExpression() (ast.Expression, error) {
	next, ok := p.Next()
	if !ok {
		return nil, UnexpectedEof(p.pos)
	}

	switch next.Type {
	case tokens.TokenLet:
		ident, err := p.Expect(tokens.TokenIdentifier)
		if err != nil {
			return nil, err
		}
		return ast.LetExpression{IdentExpression: ast.IdentExpression{
			Ident: ident.Value,
			Pos:   ident.Pos,
		}}, nil
	case tokens.TokenBegin:
		_, block, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
			tokens.TokenEnd: {},
		})
		return block, err
	case tokens.TokenIf:
		return p.ParseIfExpression()
	case tokens.TokenFor:
		return p.ParseForExpression()
	case tokens.TokenWhile:
		return p.ParseWhileExpression()
	case tokens.TokenRepeat:
		return p.ParseRepeatUntilExpression()
	case tokens.TokenUsing:
		return p.ParseUsingExpression()
	case tokens.TokenSub:
		return p.ParseSubroutineExpression()
	case tokens.TokenMod:
		return p.ParseModule()
	case tokens.TokenImport:
		return p.ParseImportExpression()
	case tokens.TokenReturn:
		return p.ParseSignalExpression(signal.SignalReturn)
	case tokens.TokenBreak:
		return p.ParseSignalExpression(signal.SignalBreak)
	case tokens.TokenContinue:
		return p.ParseSignalExpression(signal.SignalContinue)
	case tokens.TokenRaise:
		return p.ParseSignalExpression(signal.SignalRaise)
	case tokens.TokenTrue:
		return ast.BoolDec{Value: true, Pos: next.Pos}, nil
	case tokens.TokenFalse:
		return ast.BoolDec{Value: false, Pos: next.Pos}, nil
	case tokens.TokenNull:
		return ast.NullDec{Pos: next.Pos}, nil
	case tokens.TokenInfinity:
		return ast.NumDec{Value: math.Inf(1), Pos: next.Pos}, nil
	case tokens.TokenNan:
		return ast.NumDec{Value: math.NaN(), Pos: next.Pos}, nil
	case tokens.TokenStop:
		return ast.ErrorDec{Value: errors.Error{
			Code:    errors.IteratorStop,
			Message: "",
		}, Pos: next.Pos}, nil
	case tokens.TokenParenthesisOpened:
		expr, err := p.Expression(power.PowerLowest, false)
		if err != nil {
			return nil, err
		}
		_, err = p.Expect(tokens.TokenParenthesisClosed)
		return expr, err
	case tokens.TokenSquareBracketOpened:
		pos := p.pos
		elements, err := p.ParseArrayDeclaration(tokens.TokenSquareBracketClosed)
		return ast.ArrayDec{
			Elements: elements,
			Pos:      pos,
		}, err
	case tokens.TokenBraceOpened:
		return p.ParseObjectDeclaration()
	case tokens.TokenIdentifier:
		return ast.IdentExpression{
			Ident: next.Value,
			Pos:   next.Pos,
		}, nil
	case tokens.TokenNumber:
		return ast.NumDec{Value: next.NumValue, Pos: next.Pos}, nil
	case tokens.TokenString:
		return ast.StringDec{Value: next.Value, Pos: next.Pos}, nil
	case tokens.TokenEof:
		return nil, UnexpectedEof(next.Pos)
	}
	if prefix, ok := next.Type.IntoPrefix(); ok {
		expr, err := p.Expression(power.PowerPrefix, false)
		if err != nil {
			return nil, err
		}
		return ast.PrefixExpression{
			Operator: prefix,
			Pos:      next.Pos,
			Value:    expr,
		}, nil
	}
	panic(next)
	// return nil, Unexpected(next.String())
}
