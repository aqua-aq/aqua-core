package lexer

import (
	"unicode"

	"github.com/aqua-aq/aqua/pkg/float"
	"github.com/aqua-aq/aqua/pkg/pos"
	"github.com/aqua-aq/aqua/source/lexer/tokens"
)

func isDigitInBase(r rune, base int) bool {
	if base < 2 || base > 16 {
		return false // 2-16
	}

	if unicode.IsDigit(r) {
		val := int(r - '0')
		return val < base
	}

	if base > 10 {
		r = unicode.ToUpper(r)
		val := int(r - 'A' + 10)
		return val >= 10 && val < base
	}

	return false
}
func (l *Lexer) Next() (rune, bool) {
	if len(l.Data) == 0 {
		return 0, false
	}
	next := l.Data[0]
	if next == '\n' {
		l.Pos.NextLine()
	} else {
		l.Pos.AddOneColumn()
	}
	l.Data = l.Data[1:]
	return next, true
}

func (l *Lexer) Peek(i int) (rune, bool) {
	if len(l.Data) <= i {
		return 0, false
	}
	return l.Data[i], true
}

func (l *Lexer) EOF() bool {
	return len(l.Data) == 0
}

func hexRuneToNibble(r rune) (rune, bool) {
	switch {
	case r >= '0' && r <= '9':
		return r - '0', true
	case r >= 'a' && r <= 'f':
		return r - 'a' + 10, true
	case r >= 'A' && r <= 'F':
		return r - 'A' + 10, true
	default:
		return 0, false
	}
}

func getNumToken(pos pos.Pos, value string, base int) tokens.Token {
	return tokens.Token{
		Type:     tokens.TokenNumber,
		Value:    value,
		NumValue: float.ParseFloatBase(value, base),
		Pos:      pos,
	}
}
