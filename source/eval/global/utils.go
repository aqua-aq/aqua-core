package global

import (
	"github.com/aqua-aq/aqua/pkg/pos"
	"github.com/aqua-aq/aqua/pkg/stacktrace"
	"github.com/aqua-aq/aqua/source/object"
	"github.com/aqua-aq/aqua/source/object/signal"
)

func Raise(name string, err object.Error) object.SubroutineResult {
	return object.SubroutineResult{
		Signal:    signal.SubroutineSignalRaise,
		SignalVal: &object.Value{InnerValue: err},
		Trace:     stacktrace.New(pos.BuiltInPos(name)),
	}
}
