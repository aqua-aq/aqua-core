package global

import (
	"github.com/aqua-aq/aqua-core/pkg/pos"
	"github.com/aqua-aq/aqua-core/pkg/stacktrace"
	"github.com/aqua-aq/aqua-core/source/object"
	"github.com/aqua-aq/aqua-core/source/object/signal"
)

func Raise(name string, err object.Error) object.SubroutineResult {
	return object.SubroutineResult{
		Signal:    signal.SubroutineSignalRaise,
		SignalVal: object.New(err),
		Trace:     stacktrace.New(pos.BuiltInPos(name)),
	}
}
