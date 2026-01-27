package stacktrace

import (
	"fmt"
	"strings"

	"github.com/vandi37/aqua/pkg/pos"
)

type Frame struct {
	Func string
	Pos  pos.Pos
}
type StackTrace struct {
	first  pos.Pos
	frames []Frame
}

func New(pos pos.Pos) StackTrace {
	return StackTrace{first: pos, frames: nil}
}
func (s StackTrace) Add(funcName string, pos pos.Pos) StackTrace {
	return StackTrace{frames: append(s.frames, Frame{funcName, pos}), first: s.first}
}
func (s StackTrace) Frames() []Frame {
	return s.frames
}
func (s StackTrace) Pos() pos.Pos {
	return s.first
}
func (s StackTrace) IsEmpty() bool {
	return len(s.frames) == 0
}

func (s StackTrace) String() string {
	b := strings.Builder{}
	fmt.Fprintf(&b, "	at %s", s.first.String())
	for _, v := range s.frames {
		fmt.Fprintf(&b, "	at %s (%s)", v.Func, v.Pos.String())
	}
	return b.String()
}
