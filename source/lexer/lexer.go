package lexer

import (
	"fmt"
	"unicode"

	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/lexer/tokens"
)

type Lexer struct {
	Pos        pos.Pos
	Data       []rune
	Tokens     []tokens.Token
	OneChar    map[rune]tokens.TokenType
	DoubleChar map[[2]rune]tokens.TokenType
	TripleChar map[[3]rune]tokens.TokenType
	KeyWords   map[string]tokens.TokenType
}

func New(data, path string) (*Lexer, error) {
	pos, err := pos.NewPos(1, 0, path)
	if err != nil {
		return nil, err
	}
	return &Lexer{
		Pos:  pos,
		Data: []rune(data),
	}, nil
}

func NewRelative(data string, relative pos.Pos) *Lexer {
	return &Lexer{
		Pos:  pos.NewRelative(relative, 1, 0),
		Data: []rune(data),
	}
}

func (l *Lexer) NextToken() (tokens.Token, error) {
	l.Trim()
	first, ok := l.Next()
	pos := l.Pos
	if !ok {
		return tokens.Token{Type: tokens.TokenEof, Value: "", Pos: pos}, nil
	}
	switch first {
	case '#':
		l.GetComment()
		return l.NextToken()
	case '"':
		return l.GetString(pos, tokens.TokenString, "string", '"')
	case '`':
		return l.GetString(pos, tokens.TokenIdentifier, "ident", '`')
	}
	if unicode.IsNumber(first) {
		return l.GetNumber(pos, first)
	}
	if unicode.IsLetter(first) || first == '_' {
		return l.GetIdent(pos, first)
	}

	second, ok := l.Peek(0)
	if ok {
		chars := [2]rune{first, second}
		t, ok := l.DoubleChar[chars]
		if ok {
			l.Next()
			return tokens.Token{Type: t, Value: string(chars[:]), Pos: pos}, nil
		}

		third, ok := l.Peek(1)
		if ok {
			chars := [3]rune{first, second, third}
			t, ok := l.TripleChar[chars]
			if ok {
				l.Next()
				l.Next()
				return tokens.Token{Type: t, Value: string(chars[:]), Pos: pos}, nil
			}
		}

	}
	t, ok := l.OneChar[first]
	if !ok {
		return tokens.Token{}, errors.Error{
			Code:    errors.SyntaxError,
			Message: fmt.Sprintf("invalid char '%s' at %s", string(first), pos.String()),
		}
	}
	return tokens.Token{Type: t, Value: string(first), Pos: pos}, nil
}

func (l *Lexer) Tokenize() error {
	token, err := l.NextToken()
	for ; err == nil; token, err = l.NextToken() {
		l.Tokens = append(l.Tokens, token)
		if token.Type == tokens.TokenEof {
			break
		}
	}
	return err
}
