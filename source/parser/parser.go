package parser

import (
	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/source/lexer/tokens"
)

type Parser struct {
	tokens []tokens.Token
	pos    pos.Pos
}

func New(tokens []tokens.Token, path string) (*Parser, error) {
	pos, err := pos.NewPos(1, 0, path)
	if err != nil {
		return nil, err
	}
	if len(tokens) != 0 {
		pos = tokens[0].Pos
	}
	return &Parser{tokens, pos}, nil
}
