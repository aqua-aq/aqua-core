package tokens

import (
	"fmt"

	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/source/operators"
)

type Token struct {
	Type     TokenType
	Value    string
	NumValue float64
	Pos      pos.Pos
}
type TokenType byte

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
	TokenUsing
	TokenSwitch
	TokenCase
	TokenDefault

	TokenSub
	TokenWith
	TokenMod
	TokenExport
	TokenImport
	TokenDelete
	TokenAs
	TokenReturn
	TokenBreak
	TokenContinue
	TokenRaise
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
	TokenQuestion

	TokenAssign
	TokenDots // ...
	TokenPtr
	TokenColumn
	TokenDot
	TokenQuestionDot
	TokenMethod
	TokenQuestionMethod
	TokenBind
	TokenComma

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
	case TokenSwitch:
		return "<switch>"
	case TokenCase:
		return "<case>"
	case TokenDefault:
		return "<default>"
	case TokenSub:
		return "<sub>"
	case TokenWith:
		return "<with>"
	case TokenMod:
		return "<mod>"
	case TokenExport:
		return "<export>"
	case TokenUsing:
		return "<using>"
	case TokenImport:
		return "<import>"
	case TokenAs:
		return "<as>"
	case TokenDelete:
		return "<delete>"
	case TokenReturn:
		return "<return>"
	case TokenBreak:
		return "<break>"
	case TokenContinue:
		return "<continue>"
	case TokenRaise:
		return "<raise>"
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
	case TokenQuestion:
		return "<question>"
	case TokenPtr:
		return "<ptr>"
	case TokenColumn:
		return "<column>"
	case TokenDot:
		return "<dot>"
	case TokenQuestionDot:
		return "<question dot>"
	case TokenQuestionMethod:
		return "<question method>"
	case TokenMethod:
		return "<method>"
	case TokenBind:
		return "<bind>"
	case TokenComma:
		return "<comma>"
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
	if t.Value == "" {
		return fmt.Sprintf(
			"%s at %s",
			t.Type.String(),
			t.Pos.String(),
		)
	}
	return fmt.Sprintf(
		"%s: \"%s\" at %s",
		t.Type.String(),
		t.Value,
		t.Pos.String(),
	)
}

func (t TokenType) IntoPrefix() (operators.PrefixOperator, bool) {
	switch t {
	case TokenNot:
		return operators.Not, true
	case TokenTypeof:
		return operators.Typeof, true
	case TokenMinus:
		return operators.Neg, true
	case TokenPtr:
		return operators.Ptr, true
	default:
		return 0, false
	}
}

func (t TokenType) IntoBin() (operators.Operator, bool) {
	switch t {
	case TokenAnd:
		return operators.And, true
	case TokenOr:
		return operators.Or, true
	case TokenXor:
		return operators.Xor, true
	case TokenIn:
		return operators.In, true
	case TokenPlus:
		return operators.Plus, true
	case TokenMinus:
		return operators.Minus, true
	case TokenMultiply:
		return operators.Multiply, true
	case TokenDivide:
		return operators.Divide, true
	case TokenModulus:
		return operators.Modulo, true
	case TokenStrongDivide:
		return operators.StrongDivide, true
	case TokenShl:
		return operators.Shl, true
	case TokenShr:
		return operators.Shr, true
	case TokenEq:
		return operators.Equal, true
	case TokenNotEq:
		return operators.NotEqual, true
	case TokenGt:
		return operators.Greater, true
	case TokenLt:
		return operators.Less, true
	case TokenGe:
		return operators.GreaterEqual, true
	case TokenLe:
		return operators.LessEqual, true
	case TokenQuestionDot:
		return operators.QuestionDot, true
	case TokenQuestion:
		return operators.Question, true
	case TokenQuestionMethod:
		return operators.QuestionMethod, true
	case TokenDot:
		return operators.Dot, true
	case TokenMethod:
		return operators.Method, true
	case TokenBind:
		return operators.Bind, true
	case TokenSquareBracketOpened:
		return operators.Index, true
	default:
		return 0, false
	}
}
