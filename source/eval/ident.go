package eval

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/vm"
)

func (i IdentExpression) GetName(vm *vm.VM, scope scope.Scope[*object.Value]) (string, object.ExpressionResult) {
	if !i.HasAt {
		return i.Ident, object.ExpressionResult{}
	}
	val, ok := scope.Get(i.Ident)
	if !ok {
		return "", object.ExpressionResult{
			SignalVal: &object.Value{InnerValue: object.Error{
				Code:    errors.ValueError,
				Message: fmt.Sprintf("identifier %s is not defined", i.Ident),
			}},
		}
	}
	return IntoString(vm, val)
}
