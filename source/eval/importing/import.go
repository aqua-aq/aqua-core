package importing

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/vandi37/aqua/env"
	"github.com/vandi37/aqua/source/errors"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/vm"
)

type ImportType byte

const (
	TypeUnknown ImportType = iota
	TypeStd
	TypeDependence
	TypeDependenceMain
	TypeFile
)

func ParseImport(current string, s string, vm *vm.VM[*object.Value]) (string, string, ImportType, error) {
	if _, ok := vm.StdLib[s]; ok {
		return s, s, TypeStd, nil
	}

	if dep, ok := vm.Config.Dependencies[s]; ok {
		return dep, dep, TypeDependenceMain, nil
	}
	if filepath.Ext(s) != "."+env.FILE_EXTENSION {
		s += "." + env.FILE_EXTENSION
	}

	if first, second, _ := strings.Cut(s, ":"); second != "" {
		if dep, ok := vm.Config.Dependencies[first]; ok {
			return dep, filepath.Join(dep, second), TypeDependence, nil
		}

	}
	currentDir := filepath.Dir(current)
	path := filepath.Join(currentDir, s)
	path, err := filepath.Abs(path)
	if err != nil {
		return "", "", TypeUnknown, err
	}
	return "", path, TypeFile, nil
}

func GetImport(current string, s string, virtualMachine *vm.VM[*object.Value]) (string, map[string]*object.Value, error) {
	dep, full, t, err := ParseImport(current, s, virtualMachine)
	if err != nil {
		return "", nil, err
	}
	if v, ok := (*virtualMachine.Paths)[full]; ok {
		return v.Name, v.Vals, nil
	}
	switch t {
	case TypeStd:
		vals, err := virtualMachine.StdLib[full]()
		return full, vals, err
	case TypeDependenceMain, TypeDependence:
		depVm, ok := (*virtualMachine.LibVms)[dep]
		if !ok {
			depVm, err = vm.NewAndLoadConfig(filepath.Join(dep, env.CONFIG), virtualMachine.Paths, virtualMachine.LibVms, virtualMachine.Run, virtualMachine.Args)
			if err != nil {
				return "", nil, err
			}
		}

		main := depVm.Config.Main
		if t == TypeDependence {
			main = full
		}
		name := strings.TrimSuffix(filepath.Base(main), filepath.Ext(main))
		vals, err := depVm.Run(depVm.Config.Main, fmt.Sprintf("%s %s", dep, name), depVm)
		return name, vals, err
	case TypeFile:
		name := strings.TrimSuffix(filepath.Base(full), filepath.Ext(full))
		vals, err := virtualMachine.Run(full, name, virtualMachine)
		return name, vals, err
	}
	return "", nil, errors.Error{
		Code:    errors.ImportError,
		Message: fmt.Sprintf("invalid import type: %d with import string '%s'", t, current),
	}
}
