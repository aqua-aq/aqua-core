package parser

import (
	"github.com/vandi37/aqua/source/ast"
	"github.com/vandi37/aqua/source/lexer/tokens"
	"github.com/vandi37/aqua/source/operators"
	"github.com/vandi37/aqua/source/power"
)

func (p *Parser) Expression(bp power.BindingPower, isBind bool) (ast.Expression, error) {
	left, err := p.PrimaryExpression()
	if err != nil {
		return nil, err
	}

	for {
		peek, ok := p.Peek()
		if !ok {
			break
		}
		if peek.Type == tokens.TokenComma {
			if bp >= power.PowerAssignment {
				break
			}
			leftArray := []ast.Expression{left}
			for ; peek.Type == tokens.TokenComma; peek, ok = p.Peek() {
				p.Move()
				expr, err := p.Expression(power.PowerAssignment, false)
				if err != nil {
					return nil, err
				}
				leftArray = append(leftArray, expr)
			}
			_, err = p.Expect(tokens.TokenAssign) // not allowing +=, -= etc in multi value assignment, may be added in the future
			if err != nil {
				return nil, err
			}
			pos := p.pos
			expr, err := p.Expression(power.PowerAssignment, false)
			if err != nil {
				return nil, err
			}
			rightArray := []ast.Expression{expr}
			for peek, ok = p.Peek(); peek.Type == tokens.TokenComma; peek, ok = p.Peek() {
				p.Move()
				expr, err := p.Expression(power.PowerAssignment, false)
				if err != nil {
					return nil, err
				}
				leftArray = append(leftArray, expr)
			}
			left = ast.AssigmentExpression{
				Left:  leftArray,
				Right: rightArray,
				Pos:   pos,
			}

		} else if peek.Type == tokens.TokenAssign {
			if bp >= power.PowerAssignment {
				break
			}
			p.Move()
			pos := p.pos
			right, err := p.Expression(power.PowerAssignment, false)
			if err != nil {
				return nil, err
			}
			left = ast.AssigmentExpression{
				Left:  []ast.Expression{left},
				Right: []ast.Expression{right},
				Pos:   pos,
			}
		} else if bin, ok := peek.Type.IntoBin(); ok {
			pos := p.pos
			p.Move()
			if peek, ok = p.Peek(); ok &&
				bin.IsValidInAssign() &&
				peek.Type == tokens.TokenAssign {
				if bp >= power.PowerAssignment {
					break
				}
				p.Move()
				right, err := p.Expression(power.PowerAssignment, false)
				if err != nil {
					return nil, err
				}
				left = ast.AssigmentExpression{
					Left:  []ast.Expression{left},
					Right: []ast.Expression{right},
					Pos:   pos,
				}
			} else {
				if bp >= bin.Power() {
					break
				}
				right, err := p.Expression(bin.Power(), bin == operators.Bind)
				if err != nil {
					return nil, err
				}
				if bin == operators.Index {
					_, err := p.Expect(tokens.TokenSquareBracketClosed)
					if err != nil {
						return nil, err
					}
				}

				left = ast.BinExpression{
					Pos:      pos,
					Operator: bin,
					Right:    right,
					Left:     left,
				}
			}
		} else if peek.Type == tokens.TokenIncrement || peek.Type == tokens.TokenDecrement {
			p.Move()
			pos := p.pos
			if bp >= power.PowerPostfix {
				break
			}
			operator := operators.Inc
			if peek.Type == tokens.TokenDecrement {
				operator = operators.Dec
			}
			left = ast.PrefixExpression{
				Operator: operator,
				Value:    left,
				Pos:      pos,
			}
		} else if peek.Type == tokens.TokenParenthesisOpened {
			if bp >= power.PowerPostfix || isBind {
				break
			}
			p.Move()
			pos := p.pos
			args, err := p.ParseArrayDeclaration(tokens.TokenParenthesisClosed)
			if err != nil {
				return nil, err
			}
			left = ast.CallExpression{
				Pos:        pos,
				Args:       args,
				Subroutine: left,
			}
		} else {
			break
		}
	}
	return left, nil
}
