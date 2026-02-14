package lexer

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/aqua-aq/aqua/pkg/pos"
	"github.com/aqua-aq/aqua/source/errors"
	"github.com/aqua-aq/aqua/source/lexer/tokens"
)

func (l *Lexer) Trim() bool {
	res := false
	for next, ok := l.Peek(0); ok && unicode.IsSpace(next); next, ok = l.Peek(0) {
		l.Next()
		res = true
	}
	return res
}

func (l *Lexer) GetComment() {
	next, ok := l.Next()
	if !ok {
		return
	}
	if next == '-' {
		for next, ok := l.Next(); ok; next, ok = l.Next() {
			if peek, pOk := l.Peek(0); pOk && next == '-' && peek == '#' {
				l.Next()
				break
			}
		}
		return
	}
	for ; ok && next != '\n'; next, ok = l.Next() {
	}
}

func (l *Lexer) GetString(pos pos.Pos, t tokens.TokenType, name string, close rune) (tokens.Token, error) {
	var sb strings.Builder
	next, ok := l.Next()
	for ; ok && next != close; next, ok = l.Next() {
		if next == '\\' {
			next, ok = l.Peek(0)
			if !ok {
				break
			}
			if next == '@' || next == '}' {
				sb.WriteRune('\\')
				continue
			}
			var err error
			next, err = l.Escape(pos, l.Pos, name)
			if err != nil {
				return tokens.Token{}, err
			}
		}
		sb.WriteRune(next)
	}

	if !ok {
		return tokens.Token{}, errors.Error{
			Code:    errors.SyntaxError,
			Message: fmt.Sprintf("%s at %s is not closed", name, pos.String()),
		}
	}
	return tokens.Token{
		Type:  t,
		Value: sb.String(),
		Pos:   pos,
	}, nil
}

func (l *Lexer) GetNumber(pos pos.Pos, first rune) (tokens.Token, error) {
	var base = 10
	if next, ok := l.Peek(0); first == '0' && ok && (next == 'x' ||
		next == 'X' ||
		next == 'o' ||
		next == 'O' ||
		next == 'b' ||
		next == 'B') {
		l.Next()
		switch next {
		case 'x', 'X':
			base = 16
		case 'o', 'O':
			base = 8
		case 'b', 'B':
			base = 2
		}
	}
	ok := true
	var sb strings.Builder
	sb.WriteRune(first)
	for next, ok := l.Peek(0); ok && (isDigitInBase(next, base) || next == '_'); next, ok = l.Peek(0) {
		l.Next()
		if next == '_' {
			continue
		}
		sb.WriteRune(next)
	}
	next, ok := l.Peek(0)
	if !ok {
		return getNumToken(pos, sb.String(), base), nil
	}
	if next == '.' {
		l.Next()
		sb.WriteRune('.')
		for next, ok = l.Peek(0); ok && (isDigitInBase(next, base) || next == '_'); next, ok = l.Peek(0) {
			l.Next()
			if next == '_' {
				continue
			}
			sb.WriteRune(next)
		}
	}
	return getNumToken(pos, sb.String(), base), nil
}

func (l *Lexer) GetIdent(pos pos.Pos, first rune) (tokens.Token, error) {
	ok := true
	var sb strings.Builder
	sb.WriteRune(first)
	for next, ok := l.Peek(0); ok && (unicode.IsLetter(next) || unicode.IsDigit(next) || next == '_'); next, ok = l.Peek(0) {
		l.Next()
		sb.WriteRune(next)
	}
	res := sb.String()
	t, ok := l.KeyWords[res]
	if !ok {
		t = tokens.TokenIdentifier
	}
	return tokens.Token{
		Type:  t,
		Value: res,
		Pos:   pos,
	}, nil
}
