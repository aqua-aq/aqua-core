package main

import (
	"fmt"
	"os"

	"github.com/vandi37/aqua/source/config"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/run"
	"github.com/vandi37/aqua/source/vm"
)

func main() {
	path := "test.aq"
	paths := make(vm.Paths[*object.Value])
	vms := make(vm.VMs[*object.Value])
	vm, err := vm.New(config.DefaultConfig(),
		&paths, &vms, run.Run)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	_, err = run.Run(path, "main", vm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

}
