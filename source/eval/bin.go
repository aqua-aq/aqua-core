package eval

import (
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/operators"
	"github.com/vandi37/aqua/source/vm"
)

func RunBin(
	vm *vm.VM, scope scope.Scope[*object.Value], clone bool,
	left, right *object.Value, operator operators.Operator,
) object.ExpressionResult {
	panic("unimplemented")
}
