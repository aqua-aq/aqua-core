package main

import (
	"fmt"

	"github.com/vandi37/aqua/pkg/pos"
	"github.com/vandi37/aqua/pkg/scope"
	"github.com/vandi37/aqua/source/eval"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/operators"
	"github.com/vandi37/aqua/source/vm"
)

func num(f float64) *object.Value {
	return &object.Value{InnerValue: object.Number{Value: f}}
}

func main() {
	vm := &vm.VM{}
	pos := pos.Pos{}
	scope := scope.New[*object.Value]()
	array := &object.Value{InnerValue: object.Array{Elements: []*object.Value{num(1), num(1), num(3), num(4), num(5)}}}

	newArray := eval.RunBin(vm, scope, false, array, num(6), operators.Plus, pos)
	if newArray.Signal.Has() {
		panic(newArray)
	}
	*array = *newArray.SignalVal
	elem := eval.RunBin(vm, scope, false, array, num(1), operators.Index, pos)
	if elem.Signal.Has() {
		panic(elem)
	}
	newElem := eval.RunBin(vm, scope, false, elem.SignalVal, num(1), operators.Plus, pos)
	if newElem.Signal.Has() {
		panic(newElem)
	}
	*elem.SignalVal = *newElem.SignalVal
	fmt.Println("Array: ", array.String())

}
