package lexer

import "github.com/vandi37/aqua/source/lexer/tokens"

func (l *Lexer) InitOneChar() {
	l.OneChar['+'] = tokens.TokenPlus
	l.OneChar['-'] = tokens.TokenMinus
	l.OneChar['*'] = tokens.TokenMultiply
	l.OneChar['/'] = tokens.TokenDivide
	l.OneChar['%'] = tokens.TokenModulus
	l.OneChar['>'] = tokens.TokenGt
	l.OneChar['<'] = tokens.TokenLt
	l.OneChar['='] = tokens.TokenAssign
	l.OneChar['&'] = tokens.TokenPtr
	l.OneChar[':'] = tokens.TokenColumn
	l.OneChar['.'] = tokens.TokenDot
	l.OneChar[','] = tokens.TokenComma
	l.OneChar['('] = tokens.TokenParenthesisOpened
	l.OneChar[')'] = tokens.TokenParenthesisClosed
	l.OneChar['['] = tokens.TokenSquareBracketOpened
	l.OneChar[']'] = tokens.TokenSquareBracketClosed
	l.OneChar['{'] = tokens.TokenBraceOpened
	l.OneChar['}'] = tokens.TokenBraceClosed
}

func (l *Lexer) InitDoubleChar() {
	l.DoubleChar[[2]rune{'/', '/'}] = tokens.TokenStrongDivide
	l.DoubleChar[[2]rune{'<', '<'}] = tokens.TokenShl
	l.DoubleChar[[2]rune{'>', '>'}] = tokens.TokenShr
	l.DoubleChar[[2]rune{'=', '='}] = tokens.TokenEq
	l.DoubleChar[[2]rune{'~', '='}] = tokens.TokenNotEq
	l.DoubleChar[[2]rune{'>', '='}] = tokens.TokenGe
	l.DoubleChar[[2]rune{'<', '='}] = tokens.TokenLe
	l.DoubleChar[[2]rune{'.', '>'}] = tokens.TokenMethod
	l.DoubleChar[[2]rune{'-', '>'}] = tokens.TokenBind
	l.DoubleChar[[2]rune{'+', '+'}] = tokens.TokenIncrement
	l.DoubleChar[[2]rune{'-', '-'}] = tokens.TokenDecrement
	l.DoubleChar[[2]rune{'?', '.'}] = tokens.TokenQuestionDot
	l.DoubleChar[[2]rune{'?', '?'}] = tokens.TokenQuestion
}

func (l *Lexer) InitTripleChar() {
	l.TripleChar[[3]rune{'.', '.', '.'}] = tokens.TokenDots
}

func (l *Lexer) InitKeywords() {
	l.KeyWords["let"] = tokens.TokenLet
	l.KeyWords["begin"] = tokens.TokenBegin
	l.KeyWords["end"] = tokens.TokenEnd
	l.KeyWords["catch"] = tokens.TokenCatch
	l.KeyWords["if"] = tokens.TokenIf
	l.KeyWords["elif"] = tokens.TokenElif
	l.KeyWords["else"] = tokens.TokenElse
	l.KeyWords["for"] = tokens.TokenFor
	l.KeyWords["while"] = tokens.TokenWhile
	l.KeyWords["repeat"] = tokens.TokenRepeat
	l.KeyWords["until"] = tokens.TokenUntil
	l.KeyWords["switch"] = tokens.TokenSwitch
	l.KeyWords["case"] = tokens.TokenCase
	l.KeyWords["default"] = tokens.TokenDefault
	l.KeyWords["using"] = tokens.TokenUsing
	l.KeyWords["sub"] = tokens.TokenSub
	l.KeyWords["with"] = tokens.TokenWith
	l.KeyWords["mod"] = tokens.TokenMod
	l.KeyWords["export"] = tokens.TokenExport
	l.KeyWords["import"] = tokens.TokenImport
	l.KeyWords["as"] = tokens.TokenAs
	l.KeyWords["return"] = tokens.TokenReturn
	l.KeyWords["break"] = tokens.TokenBreak
	l.KeyWords["continue"] = tokens.TokenContinue
	l.KeyWords["raise"] = tokens.TokenRaise
	l.KeyWords["enum"] = tokens.TokenEnum
	l.KeyWords["and"] = tokens.TokenAnd
	l.KeyWords["or"] = tokens.TokenOr
	l.KeyWords["xor"] = tokens.TokenXor
	l.KeyWords["in"] = tokens.TokenIn
	l.KeyWords["not"] = tokens.TokenNot
	l.KeyWords["typeof"] = tokens.TokenTypeof
	l.KeyWords["true"] = tokens.TokenTrue
	l.KeyWords["false"] = tokens.TokenFalse
	l.KeyWords["null"] = tokens.TokenNull
	l.KeyWords["infinity"] = tokens.TokenInfinity
	l.KeyWords["nan"] = tokens.TokenNan
}

func (l *Lexer) Init() {
	l.OneChar = make(map[rune]tokens.TokenType)
	l.DoubleChar = make(map[[2]rune]tokens.TokenType)
	l.TripleChar = make(map[[3]rune]tokens.TokenType)
	l.KeyWords = make(map[string]tokens.TokenType)

	l.InitOneChar()
	l.InitDoubleChar()
	l.InitTripleChar()
	l.InitKeywords()
}
