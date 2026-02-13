package vm

import (
	"os"
	"path/filepath"

	"github.com/vandi37/aqua/env"
	"github.com/vandi37/aqua/source/config"
)

type Paths[T any] map[string]struct {
	Vals map[string]T
	Name string
}

type VMs[T any] map[string]*VM[T]

type VM[T any] struct {
	Paths  *Paths[T]
	StdLib map[string]func() (map[string]T, error)
	LibVms *VMs[T]
	Run    func(path, name string, vm *VM[T]) (map[string]T, error)
	Config config.Config
	Args   []string
}

func New[T any](config config.Config, paths *Paths[T], libVms *VMs[T], run func(path, name string, vm *VM[T]) (map[string]T, error), args []string) (*VM[T], error) {
	err := ResolveDependencies(config.Dependencies)
	if err != nil {
		return nil, err
	}
	std := map[string]func() (map[string]T, error){}
	return &VM[T]{
		Paths:  paths,
		StdLib: std,
		LibVms: libVms,
		Run:    run,
		Config: config,
		Args:   args,
	}, nil
}

func ResolveDependencies(dependencies map[string]string) error {
	for k, v := range dependencies {
		abs, err := filepath.Abs(filepath.Join(env.AQUA_PATH, env.LIB_PATH, v))
		if err != nil {
			return err
		}
		dependencies[k] = abs
	}
	return nil
}

func NewAndLoadConfig[T any](path string, paths *Paths[T], libVms *VMs[T], run func(path, name string, vm *VM[T]) (map[string]T, error), args []string) (*VM[T], error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	config, err := config.NewConfig(file)
	if err != nil {
		return nil, err
	}
	return New(config, paths, libVms, run, args)
}
