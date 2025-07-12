package errors

import "fmt"

type Error struct {
	Code
	Message string
}

type Code uint

const (
	UnknownError Code = iota
	TypeError
	ValueError
	StackUnderflow
	NullPointer
	InvalidSubroutine
	InvalidSignal
)

func (c Code) Error() string {
	switch c {
	case UnknownError:
		return "unknown error"
	case TypeError:
		return "type error"
	case ValueError:
		return "value error"
	case InvalidSignal:
		return "invalid signal"
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
