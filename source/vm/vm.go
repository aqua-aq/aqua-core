package vm

import "github.com/vandi37/aqua/source/object"

type VM struct {
	Files map[string]map[string]*object.Value
}

func New() *VM {
	return &VM{
		Files: map[string]map[string]*object.Value{},
	}
}
