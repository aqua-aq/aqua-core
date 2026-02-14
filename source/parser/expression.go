package parser

import (
	"github.com/aqua-aq/aqua-core/source/ast"
	"github.com/aqua-aq/aqua-core/source/lexer/tokens"
	"github.com/aqua-aq/aqua-core/source/operators"
	"github.com/aqua-aq/aqua-core/source/power"
)

func (p *Parser) Expression(bp power.BindingPower, isBind bool) (ast.Expression, error) {
	left, err := p.PrimaryExpression()
	if err != nil {
		return nil, err
	}
	for {
		peek, ok := p.Peek(0)
		if !ok {
			break
		}
		if peek.Type == tokens.TokenComma || peek.Type == tokens.TokenColumn {
			if bp >= power.PowerPatternAssigment {
				break
			}
			var name *ast.IdentExpression
			if peek.Type == tokens.TokenColumn {
				p.Move()
				ident, err := p.Expect(tokens.TokenIdentifier)
				if err != nil {
					return nil, err
				}
				name = &ast.IdentExpression{Ident: ident.Value, Pos: ident.Pos}
			}
			leftArray := []ast.AssigmentPattern{{Name: name, Expression: left, Pos: peek.Pos}}
			for peek, _ = p.Peek(0); peek.Type == tokens.TokenComma; peek, _ = p.Peek(0) {
				p.Move()
				expr, err := p.Expression(power.PowerDirectAssignment, false)
				if err != nil {
					return nil, err
				}
				peek, _ = p.Peek(0)
				if peek.Type == tokens.TokenColumn {
					p.Move()
					ident, err := p.Expect(tokens.TokenIdentifier)
					if err != nil {
						return nil, err
					}
					name = &ast.IdentExpression{Ident: ident.Value, Pos: ident.Pos}
				}
				leftArray = append(leftArray, ast.AssigmentPattern{Expression: expr, Name: name, Pos: peek.Pos})
			}
			_, err = p.Expect(tokens.TokenAssign)
			if err != nil {
				return nil, err
			}
			pos := p.pos
			right, err := p.Expression(power.PowerLowest, false)
			if err != nil {
				return nil, err
			}

			left = ast.AssigmentExpression{
				Left:  leftArray,
				Right: right,
				Pos:   pos,
			}

		} else if peek.Type == tokens.TokenAssign {
			if bp >= power.PowerDirectAssignment {
				break
			}
			p.Move()
			pos := p.pos
			right, err := p.Expression(power.PowerPatternAssigment, false)
			if err != nil {
				return nil, err
			}
			left = ast.AssigmentExpression{
				ExpressionLeft: &struct {
					ast.Expression
					operators.Operator
				}{
					Expression: left,
				},
				Right: right,
				Pos:   pos,
			}
		} else if bin, ok := peek.Type.IntoBin(); ok {
			pos := p.pos
			if peek, ok = p.Peek(1); ok &&
				bin.IsValidInAssign() &&
				peek.Type == tokens.TokenAssign {
				if bp >= power.PowerDirectAssignment {
					break
				}
				p.MoveN(2)
				right, err := p.Expression(power.PowerPatternAssigment, false)
				if err != nil {
					return nil, err
				}
				left = ast.AssigmentExpression{
					ExpressionLeft: &struct {
						ast.Expression
						operators.Operator
					}{
						Expression: left,
						Operator:   bin,
					},
					Right: right,
					Pos:   pos,
				}
			} else {
				if bp >= bin.Power() {
					break
				}
				p.Move()
				var hasColumn bool
				if bin == operators.Index {
					peek, _ := p.Peek(0)
					hasColumn = peek.Type == tokens.TokenColumn
				}
				if hasColumn {
					p.Move()
					if peek, _ := p.Peek(0); peek.Type == tokens.TokenSquareBracketClosed {
						p.Move()
						left = ast.SliceExpression{
							Pos:  pos,
							Left: left,
						}
						continue
					}
				}
				power := bin.Power()
				if bin.IsRight() {
					power--
				}
				if bin == operators.Index {
					power = 0
				}
				right, err := p.Expression(power, bin == operators.Bind)
				if err != nil {
					return nil, err
				}
				if hasColumn {
					_, err := p.Expect(tokens.TokenSquareBracketClosed)
					if err != nil {
						return nil, err
					}
					left = ast.SliceExpression{
						Pos:  pos,
						Left: left,
						End:  right,
					}
					continue
				}
				if bin == operators.Index {
					if peek, _ := p.Peek(0); peek.Type == tokens.TokenColumn {
						p.Move()
						if peek, _ := p.Peek(0); peek.Type == tokens.TokenSquareBracketClosed {
							p.Move()
							left = ast.SliceExpression{
								Pos:   pos,
								Left:  left,
								Start: right,
							}
							continue
						}
						end, err := p.Expression(0, bin == operators.Bind)
						if err != nil {
							return nil, err
						}
						_, err = p.Expect(tokens.TokenSquareBracketClosed)
						if err != nil {
							return nil, err
						}
						left = ast.SliceExpression{
							Pos:   pos,
							Left:  left,
							Start: right,
							End:   end,
						}
						continue
					}
					_, err = p.Expect(tokens.TokenSquareBracketClosed)
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
			if bp >= power.PowerPostfix {
				break
			}
			p.Move()
			pos := p.pos

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
