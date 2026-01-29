package tokens

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/pos"
)

type Token struct {
	Type     TokenType
	Value    string
	NumValue float64
	Pos      pos.Pos
}
type TokenType uint16

const (
	TokenEof TokenType = iota

	// Keywords
	TokenLet
	TokenBegin
	TokenEnd
	TokenCatch
	TokenIf
	TokenElif
	TokenElse
	TokenFor
	TokenWhile
	TokenRepeat
	TokenUntil

	TokenSub
	TokenWith
	TokenMod
	TokenImport
	TokenAs
	TokenReturn
	TokenBreak
	TokenContinue
	TokenEnum

	TokenAnd
	TokenOr
	TokenXor
	TokenIn
	TokenNot
	TokenTypeof

	TokenTrue
	TokenFalse
	TokenNull
	TokenInfinity
	TokenNan

	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenModulus
	TokenStrongDivide
	TokenShl
	TokenShr
	TokenEq
	TokenNotEq
	TokenGt
	TokenLt
	TokenGe
	TokenLe

	TokenAssign
	TokenDots // ...
	TokenAt
	TokenPtr
	TokenColumn
	TokenDot
	TokenMethod
	TokenBind
	TokenComma
	TokenArrow // =>

	TokenParenthesisOpened
	TokenParenthesisClosed
	TokenSquareBracketOpened
	TokenSquareBracketClosed
	TokenBraceOpened
	TokenBraceClosed

	TokenIncrement
	TokenDecrement

	TokenIdentifier
	TokenNumber
	TokenString
)

func (t TokenType) String() string {
	switch t {
	case TokenEof:
		return "<eof>"
	case TokenLet:
		return "<let>"
	case TokenBegin:
		return "<begin>"
	case TokenEnd:
		return "<end>"
	case TokenCatch:
		return "<catch>"
	case TokenIf:
		return "<if>"
	case TokenElif:
		return "<elif>"
	case TokenElse:
		return "<else>"
	case TokenFor:
		return "<for>"
	case TokenWhile:
		return "<while>"
	case TokenRepeat:
		return "<repeat>"
	case TokenUntil:
		return "<until>"
	case TokenSub:
		return "<sub>"
	case TokenWith:
		return "<with>"
	case TokenMod:
		return "<mod>"
	case TokenImport:
		return "<import>"
	case TokenAs:
		return "<as>"
	case TokenReturn:
		return "<return>"
	case TokenBreak:
		return "<break>"
	case TokenContinue:
		return "<continue>"
	case TokenEnum:
		return "<enum>"
	case TokenAnd:
		return "<and>"
	case TokenOr:
		return "<or>"
	case TokenXor:
		return "<xor>"
	case TokenIn:
		return "<in>"
	case TokenNot:
		return "<not>"
	case TokenTypeof:
		return "<typeof>"
	case TokenTrue:
		return "<true>"
	case TokenFalse:
		return "<false>"
	case TokenNull:
		return "<null>"
	case TokenInfinity:
		return "<infinity>"
	case TokenNan:
		return "<nan>"
	case TokenPlus:
		return "<plus>"
	case TokenMinus:
		return "<minus>"
	case TokenMultiply:
		return "<multiply>"
	case TokenDivide:
		return "<divide>"
	case TokenModulus:
		return "<modulus>"
	case TokenStrongDivide:
		return "<string divide>"
	case TokenShl:
		return "<shl>"
	case TokenShr:
		return "<shr>"
	case TokenEq:
		return "<eq>"
	case TokenNotEq:
		return "<ne>"
	case TokenGt:
		return "<gt>"
	case TokenLt:
		return "<lt>"
	case TokenGe:
		return "<ge>"
	case TokenLe:
		return "<le>"
	case TokenAssign:
		return "<assign>"
	case TokenDots:
		return "<dots>"
	case TokenAt:
		return "<at>"
	case TokenPtr:
		return "<ptr>"
	case TokenColumn:
		return "<column>"
	case TokenDot:
		return "<dot>"
	case TokenMethod:
		return "<method>"
	case TokenBind:
		return "<bind>"
	case TokenComma:
		return "<comma>"
	case TokenArrow:
		return "<arrow>"
	case TokenParenthesisOpened:
		return "<parenthesis opened>"
	case TokenParenthesisClosed:
		return "<parenthesis closed>"
	case TokenSquareBracketOpened:
		return "<square bracket opened>"
	case TokenSquareBracketClosed:
		return "<square bracket closed>"
	case TokenBraceOpened:
		return "<brace opened>"
	case TokenBraceClosed:
		return "<brace closed>"
	case TokenIncrement:
		return "<increment>"
	case TokenDecrement:
		return "<decrement>"

	case TokenIdentifier:
		return "<identifier>"
	case TokenNumber:
		return "<number>"
	case TokenString:
		return "<string>"
	default:
		return "<unknown>"
	}
}

func (t Token) String() string {
	return fmt.Sprintf(
		"%s: \"%s\" at %s",
		t.Type.String(),
		t.Value,
		t.Pos.String(),
	)
}
