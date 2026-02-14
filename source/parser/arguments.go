package parser

import (
	"github.com/aqua-aq/aqua/source/ast"
	"github.com/aqua-aq/aqua/source/lexer/tokens"
	"github.com/aqua-aq/aqua/source/power"
)

func (p *Parser) ParseArguments() (ast.Arguments, error) {
	var res ast.Arguments
	for {
		peek, _ := p.Peek(0)
		if peek.Type == tokens.TokenDots {
			p.Move()
			ident, err := p.Expect(tokens.TokenIdentifier)
			if err != nil {
				return ast.Arguments{}, err
			}
			res.Last = &ident.Value
			return res, nil
		}
		if peek.Type != tokens.TokenIdentifier {
			return res, nil
		}
		ident, err := p.Expect(tokens.TokenIdentifier)
		if err != nil {
			return ast.Arguments{}, err
		}
		peek, ok := p.Peek(0)
		if !ok {
			return res, nil
		}
		var def ast.Expression
		if peek.Type == tokens.TokenAssign {
			p.Move()
			def, err = p.Expression(power.PowerAssignment, false)
			if err != nil {
				return ast.Arguments{}, err
			}
			peek, ok = p.Peek(0)
			if !ok {
				return res, nil // the next step will think about it
			}
		}
		res.Elements = append(res.Elements, struct {
			Name    ast.IdentExpression
			Default ast.Expression
		}{
			Name:    ast.IdentExpression{Ident: ident.Value, Pos: ident.Pos},
			Default: def,
		})

		if peek.Type != tokens.TokenComma {
			return res, nil
		}
		p.Move()
	}
}
