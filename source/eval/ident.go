package eval

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/pkg/stacktrace"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/vm"
)

func (i IdentExpression) GetName(vm *vm.VM, scope scope.Scope[*object.Value], pos pos.Pos) (string, object.ExpressionResult) {
	if !i.HasAt {
		return i.Ident, object.ExpressionResult{Trace: stacktrace.New(pos)}
	}
	val, ok := scope.Get(i.Ident)
	if !ok {
		return "", object.ExpressionResult{
			Trace: stacktrace.New(pos),
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.ValueError,
				Message: fmt.Sprintf("identifier %s is not defined", i.Ident),
			}},
		}
	}
	return IntoString(vm, val, pos)
}
