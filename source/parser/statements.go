package parser

import (
	"github.com/vandi37/aqua/source/ast"
	"github.com/vandi37/aqua/source/lexer/tokens"
	"github.com/vandi37/aqua/source/power"
	"github.com/vandi37/aqua/source/signal"
)

func (p *Parser) ParseBlockExpression(
	endings map[tokens.TokenType]struct{},
) (tokens.TokenType, ast.BlockExpression, error) {
	endings[tokens.TokenCatch] = struct{}{}
	pos := p.pos
	expressions := []ast.Expression{}
	peek, ok := p.Peek(0)
	isEnding := func() bool {
		_, ok := endings[peek.Type]
		return ok
	}
	for ; ok && !isEnding(); peek, ok = p.Peek(0) {
		expr, err := p.Expression(power.PowerLowest, false)
		if err != nil {
			return tokens.TokenEof, ast.BlockExpression{}, err
		}
		expressions = append(expressions, expr)

	}
	if !ok {
		return tokens.TokenEof, ast.BlockExpression{}, UnexpectedEof(p.pos)
	}
	p.Move() // move after peek
	ending := peek.Type
	var catch *ast.CatchBlock
	if ending == tokens.TokenCatch {
		catchPos := p.pos
		ident, err := p.Expect(tokens.TokenIdentifier)
		if err != nil {
			return tokens.TokenEof, ast.BlockExpression{}, err
		}
		var block ast.BlockExpression
		ending, block, err = p.ParseBlockExpression(endings)
		if err != nil {
			return tokens.TokenEof, ast.BlockExpression{}, err
		}
		catch = &ast.CatchBlock{
			Name: ast.IdentExpression{
				Ident: ident.Value,
				Pos:   ident.Pos,
			},
			Pos:         catchPos,
			Expressions: block,
		}
	}
	delete(endings, tokens.TokenCatch)
	return ending, ast.BlockExpression{
		Expressions: expressions,
		Pos:         pos,
		Catch:       catch,
	}, nil
}

func (p *Parser) ParseIfExpression() (ast.IfExpression, error) {
	pos := p.pos
	condition, err := p.Expression(power.PowerLowest, false)
	if err != nil {
		return ast.IfExpression{}, err
	}
	ending, block, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
		tokens.TokenEnd:  {},
		tokens.TokenElif: {},
		tokens.TokenElse: {},
	})
	if err != nil {
		return ast.IfExpression{}, err
	}
	ifBlock := block
	var elifs []ast.ElifBlock
	for ending == tokens.TokenElif {
		elifPos := p.pos
		elifCondition, err := p.Expression(power.PowerLowest, false)
		if err != nil {
			return ast.IfExpression{}, err
		}

		ending, block, err = p.ParseBlockExpression(map[tokens.TokenType]struct{}{
			tokens.TokenEnd:  {},
			tokens.TokenElif: {},
			tokens.TokenElse: {},
		})
		if err != nil {
			return ast.IfExpression{}, err
		}

		elifs = append(elifs, ast.ElifBlock{
			Pos:       elifPos,
			Block:     block,
			Condition: elifCondition,
		})
	}

	var elseBlock *ast.BlockExpression
	if ending == tokens.TokenElse {
		_, block, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
			tokens.TokenEnd: {},
		})
		if err != nil {
			return ast.IfExpression{}, err
		}
		elseBlock = &block
	}

	return ast.IfExpression{
		Pos:       pos,
		If:        ifBlock,
		Condition: condition,
		Elifs:     elifs,
		Else:      elseBlock,
	}, nil
}

func (p *Parser) ParseForExpression() (ast.ForExpression, error) {
	pos := p.pos
	args, err := p.ParseArguments()
	if err != nil {
		return ast.ForExpression{}, err
	}
	_, err = p.Expect(tokens.TokenIn)
	if err != nil {
		return ast.ForExpression{}, err
	}
	peek, ok := p.Peek(0)
	if !ok {
		return ast.ForExpression{}, UnexpectedEof(p.pos)
	}
	var isEnum bool
	if peek.Type == tokens.TokenEnum {
		isEnum = true
		p.Move()
	}
	expression, err := p.Expression(power.PowerLowest, false)
	if err != nil {
		return ast.ForExpression{}, err
	}
	ending, block, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
		tokens.TokenEnd:  {},
		tokens.TokenElse: {},
	})
	if err != nil {
		return ast.ForExpression{}, err
	}
	var elseBlock *ast.BlockExpression
	if ending == tokens.TokenElse {
		_, b, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
			tokens.TokenEnd: {},
		})
		if err != nil {
			return ast.ForExpression{}, err
		}
		elseBlock = &b
	}

	return ast.ForExpression{
		Pos:        pos,
		Arguments:  args,
		Expression: expression,
		IsEnum:     isEnum,
		Block:      block,
		Else:       elseBlock,
	}, nil
}
func (p *Parser) ParseWhileExpression() (ast.WhileExpression, error) {
	pos := p.pos
	condition, err := p.Expression(power.PowerLowest, false)
	if err != nil {
		return ast.WhileExpression{}, err
	}
	ending, block, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
		tokens.TokenEnd:  {},
		tokens.TokenElse: {},
	})
	if err != nil {
		return ast.WhileExpression{}, err
	}
	var elseBlock *ast.BlockExpression
	if ending == tokens.TokenElse {
		_, b, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
			tokens.TokenEnd: {},
		})
		if err != nil {
			return ast.WhileExpression{}, err
		}
		elseBlock = &b
	}
	return ast.WhileExpression{
		Pos:       pos,
		IsWhile:   true,
		Condition: condition,
		Block:     block,
		Else:      elseBlock,
	}, nil
}

func (p *Parser) ParseRepeatUntilExpression() (ast.WhileExpression, error) {
	pos := p.pos
	ending, block, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
		tokens.TokenUntil: {},
		tokens.TokenElse:  {},
	})
	if err != nil {
		return ast.WhileExpression{}, err
	}
	var elseBlock *ast.BlockExpression
	if ending == tokens.TokenElse {
		_, b, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
			tokens.TokenUntil: {},
		})
		if err != nil {
			return ast.WhileExpression{}, err
		}
		elseBlock = &b
	}
	condition, err := p.Expression(power.PowerLowest, false)
	if err != nil {
		return ast.WhileExpression{}, err
	}
	return ast.WhileExpression{
		Pos:       pos,
		IsWhile:   false,
		Condition: condition,
		Block:     block,
		Else:      elseBlock,
	}, nil
}

func (p *Parser) ParseUsingExpression() (ast.UsingExpression, error) {
	pos := p.pos
	expr, err := p.Expression(power.PowerLowest, false)
	if err != nil {
		return ast.UsingExpression{}, err
	}
	var name *ast.IdentExpression
	if peek, ok := p.Peek(0); ok && peek.Type == tokens.TokenAs {
		p.Move()
		ident, err := p.Expect(tokens.TokenIdentifier)
		if err != nil {
			return ast.UsingExpression{}, err
		}
		name = &ast.IdentExpression{
			Ident: ident.Value,
			Pos:   ident.Pos,
		}
	}

	_, block, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
		tokens.TokenEnd: {},
	})
	if err != nil {
		return ast.UsingExpression{}, err
	}
	return ast.UsingExpression{
		Pos:        pos,
		Name:       name,
		Expression: expr,
		Block:      block,
	}, nil

}

func (p *Parser) ParseSubroutineExpression() (ast.SubroutineDec, error) {
	pos := p.pos
	var ident *ast.IdentExpression
	if peek, ok := p.Peek(0); ok && peek.Type == tokens.TokenIdentifier {
		p.Move()
		ident = &ast.IdentExpression{
			Ident: peek.Value,
			Pos:   peek.Pos,
		}
	}

	if _, err := p.Expect(tokens.TokenParenthesisOpened); err != nil {
		return ast.SubroutineDec{}, err
	}
	args, err := p.ParseArguments()
	if err != nil {
		return ast.SubroutineDec{}, err
	}
	if _, err := p.Expect(tokens.TokenParenthesisClosed); err != nil {
		return ast.SubroutineDec{}, err
	}
	ending, block, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
		tokens.TokenEnd:  {},
		tokens.TokenWith: {},
	})
	if err != nil {
		return ast.SubroutineDec{}, err
	}
	var prototype ast.Expression
	if ending == tokens.TokenWith {
		prototype, err = p.Expression(power.PowerLowest, false)
		if err != nil {
			return ast.SubroutineDec{}, err
		}
	}

	return ast.SubroutineDec{
		Name:      ident,
		Pos:       pos,
		Arguments: args,
		Body:      block,
		Prototype: prototype,
		IsGlobal:  ident != nil,
	}, nil
}

func (p *Parser) ParseModule() (ast.ModExpression, error) {
	pos := p.pos
	ident, err := p.Expect(tokens.TokenIdentifier)
	if err != nil {
		return ast.ModExpression{}, err
	}
	_, block, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{
		tokens.TokenExport: {},
	})
	if err != nil {
		return ast.ModExpression{}, err
	}
	firstExport, err := p.Expect(tokens.TokenIdentifier)
	if err != nil {
		return ast.ModExpression{}, err
	}
	export := []string{firstExport.Value}
	peek, ok := p.Peek(0)
	for ; ok && peek.Type == tokens.TokenComma; peek, ok = p.Peek(0) {
		p.Move()
		next, err := p.Expect(tokens.TokenIdentifier)
		if err != nil {
			return ast.ModExpression{}, err
		}
		export = append(export, next.Value)
	}
	return ast.ModExpression{
		Pos: pos,
		Name: ast.IdentExpression{
			Ident: ident.Value,
			Pos:   ident.Pos,
		},
		Body:   block,
		Export: export,
	}, nil
}
func (p *Parser) ParseSignalExpression(signal signal.Signal) (ast.SignalExpression, error) {
	pos := p.pos
	var expr ast.Expression = ast.NullDec{Pos: pos}
	if peek, ok := p.Peek(0); ok && peek.Type == tokens.TokenColumn {
		p.Move()
		var err error
		expr, err = p.Expression(power.PowerLowest, false)
		if err != nil {
			return ast.SignalExpression{}, err
		}

	}

	return ast.SignalExpression{
		Signal: signal,
		SigVal: expr,
		Pos:    pos,
	}, nil
}

func (p *Parser) ParseSwitchExpression() (ast.SwitchExpression, error) {
	pos := p.pos
	value, err := p.Expression(power.PowerLowest, false)
	if err != nil {
		return ast.SwitchExpression{}, err
	}
	_, err = p.Expect(tokens.TokenCase)
	if err != nil {
		return ast.SwitchExpression{}, err
	}
	ending := tokens.TokenCase
	var cases []ast.Case
	endings := map[tokens.TokenType]struct{}{
		tokens.TokenCase:    {},
		tokens.TokenDefault: {},
		tokens.TokenEnd:     {},
	}
	for ending == tokens.TokenCase {
		curPos := p.pos
		expr, err := p.Expression(power.PowerLowest, false)
		if err != nil {
			return ast.SwitchExpression{}, err
		}
		var block ast.BlockExpression
		ending, block, err = p.ParseBlockExpression(endings)
		if err != nil {
			return ast.SwitchExpression{}, err
		}
		cases = append(cases, ast.Case{
			Pos:        curPos,
			Expression: expr,
			Block:      block,
		})
	}
	var d *ast.BlockExpression
	if ending == tokens.TokenDefault {
		_, block, err := p.ParseBlockExpression(map[tokens.TokenType]struct{}{tokens.TokenEnd: {}})
		if err != nil {
			return ast.SwitchExpression{}, err
		}
		d = &block
	}

	return ast.SwitchExpression{
		Pos:     pos,
		Value:   value,
		Cases:   cases,
		Default: d,
	}, nil
}
