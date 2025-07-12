package signal

type Signal byte

const (
	SignalNone = iota
	SignalReturn
	SignalRaise
	SignalBreak
	SignalContinue
)

type SubroutineSignal bool

const (
	SubroutineSignalReturn = false
	SubroutineSignalRaise  = true
)

func (s Signal) String() string {
	switch s {
	case SignalNone:
		return "none"
	case SignalReturn:
		return "return"
	case SignalBreak:
		return "break"
	case SignalContinue:
		return "continue"
	case SignalRaise:
		return "raise"
	default:
		return "unknown"
	}
}

func (s Signal) Has() bool {
	return s != SignalNone
}

func (s Signal) IntoSubroutineSignal() (SubroutineSignal, bool) {
	return s == SignalRaise, s == SignalRaise || s == SignalReturn
}

func (s SubroutineSignal) IntoSignal() Signal {
	if s {
		return SignalRaise
	}
	return SignalReturn
}

func (s SubroutineSignal) String() string {
	return s.IntoSignal().String()
}
