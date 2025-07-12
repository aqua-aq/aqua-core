package vm


type VM struct {
}

func New() *VM {
	return &VM{}
}

// func (vm *VM) Call() object.Value {
// 	top, err := vm.Pop()
// 	if err != nil {
// 		return err
// 	}
// 	if object.IsNull(top) {
// 		return &object.Error{Code: errors.NullPointer}
// 	}
// 	length, ok := top.(*object.Int)
// 	if !ok {
// 		return &object.Error{Code: errors.TypeError, Message: "expected an integer for argument numbers"}
// 	}
// 	args := vm.stack.Cut(length.Value)
// 	last, err := vm.Pop()
// 	if err != nil {
// 		return err
// 	}
// 	if object.IsNull(last) {
// 		return &object.Error{Code: errors.NullPointer}
// 	}
// 	proc, ok := last.(*object.Subroutine)
// 	if !ok {
// 		return &object.Error{Code: errors.TypeError, Message: "expected a procedure for calling"}
// 	}
// 	scope := proc.Scope
// 	object.ParseArgs(proc.Elements, proc.Last, args, &scope)
// 	if proc.It != nil {
// 		proc.Scope.Set(IT, proc.It)
// 	}
// 	if proc.BuildIn != nil {
// 		obj, err := proc.BuildIn(&scope)
// 		if err != nil {
// 			return err
// 		}
// 		vm.Push(obj)
// 	} else if proc.Code == nil {
// 		return &object.Error{Code: errors.InvalidProcedure, Message: "expected to have either a build in function or code"}
// 	}
// 	// doing code
// 	return nil
// }
