package lexer

import (
	"fmt"
	"unicode/utf8"

	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/source/errors"
)



func (l *Lexer) Escape(pos pos.Pos, escPos pos.Pos, name string) (rune, error) {
	isNotClosed := errors.Error{
		Code:    errors.SyntaxError,
		Message: fmt.Sprintf("%s at %s is not closed", name, pos.String()),
	}
	next, ok := l.Next()
	if !ok {
		return 0, isNotClosed
	}
	switch next {
	case '\\', '"', '`':
		return next, nil
	case 'n':
		return '\n', nil
	case 'r':
		return '\r', nil
	case 't':
		return '\t', nil
	case 'a':
		return '\a', nil
	case 'b':
		return '\b', nil
	case 'v':
		return '\v', nil
	case 'x', 'u', 'U':
		first, err := l.readHex(escPos, isNotClosed)
		if err != nil {
			return 0, err
		}

		second, err := l.readHex(escPos, isNotClosed)
		if err != nil {
			return 0, err
		}
		if next == 'x' {
			return (first << 4) | second, nil
		}
		third, err := l.readHex(escPos, isNotClosed)
		if err != nil {
			return 0, err
		}
		forth, err := l.readHex(escPos, isNotClosed)
		if err != nil {
			return 0, err
		}
		if next == 'u' {
			val := first<<12 | second<<8 | third<<4 | forth
			if !utf8.ValidRune(val) {
				return 0, errors.Error{
					Code:    errors.SyntaxError,
					Message: fmt.Sprintf("invalid unicode code point U+%X at %s", val, escPos.String()),
				}
			}
			return val, nil
		}
		fifth, err := l.readHex(escPos, isNotClosed)
		if err != nil {
			return 0, err
		}
		sixth, err := l.readHex(escPos, isNotClosed)
		if err != nil {
			return 0, err
		}
		seventh, err := l.readHex(escPos, isNotClosed)
		if err != nil {
			return 0, err
		}
		eight, err := l.readHex(escPos, isNotClosed)
		if err != nil {
			return 0, err
		}
		val := first<<28 | second<<24 | third<<20 | forth<<16 |
			fifth<<12 | sixth<<8 | seventh<<4 | eight
		if !utf8.ValidRune(val) {
			return 0, errors.Error{
				Code:    errors.SyntaxError,
				Message: fmt.Sprintf("invalid unicode code point U+%X at %s", val, escPos.String()),
			}
		}
		return val, nil
	}
	return 0, errors.Error{
		Code:    errors.SyntaxError,
		Message: fmt.Sprintf("unknown escape sequence '\\%c' at %s", next, escPos.String()),
	}
}
func (l *Lexer) readHex(pos pos.Pos, err error) (rune, error) {
	r, ok := l.Next()
	if !ok {
		return 0, err
	}
	n, ok := hexRuneToNibble(r)
	if !ok {
		return 0, errors.Error{
			Code:    errors.SyntaxError,
			Message: fmt.Sprintf("invalid hex digit '%c' at %s", r, pos.String()),
		}
	}
	return n, nil
}
