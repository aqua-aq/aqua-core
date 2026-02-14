package global

import (
	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/object/signal"
)

func Raise(name string, err object.Error) object.SubroutineResult {
	return object.SubroutineResult{
		Signal:    signal.SubroutineSignalRaise,
		SignalVal: &object.Value{InnerValue: err},
		Trace:     stacktrace.New(pos.BuiltInPos(name)),
	}
}
