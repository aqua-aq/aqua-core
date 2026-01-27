package tokens

import "github.com/vandi37/aqua/pkg/pos"

type Token struct {
	Type  TokenType
	Value string
	Pos   pos.Pos
}
type TokenType uint16

const (
	TokenEof TokenType = iota
	TokenWhiteSpace
	TokenComment

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
	TokenReturn
	TokenBreak
	TokenContinue

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

	TokenParenthesisOpen
	TokenParenthesisClosed
	TokenSquareBracketOpen
	TokenSquareBracketClosed
	TokenBraceOpen
	TokenBraceClosed

	TokenIncrement
	TokenDecrement

	TokenIdentifier
	TokenNumber
	TokenString
)
