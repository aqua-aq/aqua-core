package errors

import "fmt"

type Error struct {
	Code    Code
	Message string
}

type Code uint16

const (
	UnknownError Code = iota
	SyntaxError
	ImportError
	TypeError
	ValueError
	DivisionByZero
	InvalidSignal
	IteratorStop
)

func (c Code) Error() string {
	switch c {
	case TypeError:
		return "type error"
	case ValueError:
		return "value error"
	case InvalidSignal:
		return "invalid signal"
	case IteratorStop:
		return "iterator stop"
	case SyntaxError:
		return "syntax error"
	case ImportError:
		return "import error"
	default:
		return "unknown error"
	}
}
func (e Error) Error() string {
	msg := e.Message
	if msg != "" {
		msg = ": " + msg
	}
	return fmt.Sprintf("%d (%v)%s", e.Code, e.Code, msg)
}
