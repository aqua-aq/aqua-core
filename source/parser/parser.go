package parser

import (
	"fmt"

	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/pkg/errors"
	"github.com/aqua-aq/aqua-core/source/lexer/tokens"
)

func UnexpectedEof(pos pos.Pos) errors.Error {
	return errors.Error{
		Code:    errors.SyntaxError,
		Message: fmt.Sprintf("unexpected <eof> at %s", pos.String()),
	}
}

func Expected(expected, got string) errors.Error {
	return errors.Error{
		Code:    errors.SyntaxError,
		Message: fmt.Sprintf("expected %s, got %s", expected, got),
	}
}

func Unexpected(got string) errors.Error {
	return errors.Error{
		Code:    errors.SyntaxError,
		Message: fmt.Sprintf("unexpected %s", got),
	}
}

type Parser struct {
	tokens []tokens.Token
	pos    pos.Pos
}

func New(tokens []tokens.Token, pos pos.Pos) *Parser {
	return &Parser{tokens, pos}
}

func (p *Parser) MoveN(n int) bool {
	if len(p.tokens) < n {
		return false
	}
	next := p.tokens[n-1]
	p.tokens = p.tokens[n:]
	p.pos = next.Pos
	return true
}

func (p *Parser) Move() bool {
	return p.MoveN(1)
}

func (p *Parser) Next() (tokens.Token, bool) {
	if len(p.tokens) == 0 {
		return tokens.Token{}, false
	}
	next := p.tokens[0]
	p.Move()
	return next, true
}
func (p *Parser) Peek(n int) (tokens.Token, bool) {
	if len(p.tokens) <= n {
		return tokens.Token{}, false
	}
	return p.tokens[n], true
}

func (p *Parser) Expect(token tokens.TokenType) (tokens.Token, error) {
	next, ok := p.Next()
	if !ok {
		return tokens.Token{}, UnexpectedEof(p.pos)
	}

	if next.Type != token {
		return next, Expected(token.String(), next.String())
	}
	return next, nil
}

func (p *Parser) Pos() pos.Pos {
	return p.pos
}
