package parser

import (
	"fmt"

	"github.com/aqua-aq/aqua-core/source/ast"
	"github.com/aqua-aq/aqua-core/source/lexer/tokens"
	"github.com/aqua-aq/aqua-core/source/power"
)

func (p *Parser) ParseImportExpression() (ast.ImportExpression, error) {
	pos := p.pos
	path, err := p.Expression(power.PowerLowest, false)
	if err != nil {
		return ast.ImportExpression{}, err
	}
	var name *ast.IdentExpression
	if peek, ok := p.Peek(0); ok && peek.Type == tokens.TokenAs {
		p.Move()
		ident, err := p.Expect(tokens.TokenIdentifier)
		if err != nil {
			return ast.ImportExpression{}, err
		}
		name = &ast.IdentExpression{
			Ident: ident.Value,
			Pos:   ident.Pos,
		}
	}

	return ast.ImportExpression{
		Pos:  pos,
		Name: name,
		Path: path,
	}, nil
}

func (p *Parser) ParseArrayDeclaration(end tokens.TokenType) ([]ast.ArrayElement, error) {
	var elements []ast.ArrayElement
	for {
		var isContinuos bool
		if len(p.tokens) == 0 {
			return nil, UnexpectedEof(p.pos)
		}
		curPos := p.tokens[0].Pos
		if peek, ok := p.Peek(0); ok && peek.Type == tokens.TokenDots {
			p.Move()
			isContinuos = true
		} else if peek.Type == end {
			break
		}
		expr, err := p.Expression(power.PowerPatternAssigment, false)
		if err != nil {
			return nil, err
		}
		elements = append(elements, ast.ArrayElement{
			Pos:         curPos,
			Value:       expr,
			IsContinuos: isContinuos,
		})
		if peek, ok := p.Peek(0); !ok || peek.Type != tokens.TokenComma {
			break
		}
		p.Move()
	}
	_, err := p.Expect(end)
	if err != nil {
		return nil, err
	}
	return elements, nil
}

func (p *Parser) ParseObjectDeclaration() (ast.ObjectDec, error) {
	pos := p.pos
	var vals []ast.ObjectVal
	for {
		if peek, ok := p.Peek(0); ok && peek.Type == tokens.TokenDots {
			p.Move()
			expr, err := p.Expression(power.PowerPatternAssigment, false)
			if err != nil {
				return ast.ObjectDec{}, err
			}
			vals = append(vals, ast.ObjectVal{
				Value:       expr,
				IsContinuos: true,
				Pos:         peek.Pos,
			})
		} else if peek.Type == tokens.TokenBraceClosed {
			break
		} else if peek.Type == tokens.TokenSub {
			p.Move()
			sub, err := p.ParseSubroutineExpression()
			if err != nil {
				return ast.ObjectDec{}, nil
			}
			if sub.Name == nil {
				return ast.ObjectDec{}, Expected("subroutine", fmt.Sprintf("anonyms subroutine at %v", sub.Pos))
			}
			sub.IsGlobal = false
			vals = append(vals, ast.ObjectVal{
				Name:  *sub.Name,
				Value: sub,
				Pos:   sub.Pos,
			})
		} else {
			ident, err := p.Expect(tokens.TokenIdentifier)
			if err != nil {
				return ast.ObjectDec{}, err
			}
			_, err = p.Expect(tokens.TokenColumn)
			if err != nil {
				return ast.ObjectDec{}, err
			}
			expr, err := p.Expression(power.PowerPatternAssigment, false)
			if err != nil {
				return ast.ObjectDec{}, err
			}
			vals = append(vals, ast.ObjectVal{
				Name: ast.IdentExpression{
					Ident: ident.Value,
					Pos:   ident.Pos,
				},
				Value: expr,
				Pos:   ident.Pos,
			})
		}
		if peek, ok := p.Peek(0); !ok || peek.Type != tokens.TokenComma {
			break
		}
		p.Move()
	}
	_, err := p.Expect(tokens.TokenBraceClosed)
	if err != nil {
		return ast.ObjectDec{}, err
	}
	return ast.ObjectDec{
		Pos:  pos,
		Vals: vals,
	}, nil
}
