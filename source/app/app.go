package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/vandi37/aqua/env"
	"github.com/vandi37/aqua/pkg/cli"
	"github.com/vandi37/aqua/source/config"
	"github.com/vandi37/aqua/source/object"
	"github.com/vandi37/aqua/source/run"
	"github.com/vandi37/aqua/source/vm"
)

var App = cli.CLI{
	Default: &cli.Args{
		Flags: map[string]bool{
			"quiet": false,
		},
		AllowMoreArgs: true,
	},
	Commands: map[string]cli.CLI{
		"help": {
			Default: &cli.Args{
				Flags: map[string]bool{},
				Run:   Help,
			},
		},
		"version": {
			Default: &cli.Args{
				Length: 0,
				Flags:  map[string]bool{},
				Run:    Version,
			},
		},
		"new": {
			Default: &cli.Args{
				Length: 1,
				Flags: map[string]bool{
					"version": true,
					"main":    true,
					"path":    true,
					"lib":     false,
				},
				Run: New,
			},
		},
		"run": {
			Default: &cli.Args{
				Flags: map[string]bool{
					"project": true,
					"file":    true,
				},
				AllowMoreArgs: true,
				Run:           Run,
			},
		},
		"exec": {
			Default: &cli.Args{
				Length:        1,
				Flags:         map[string]bool{},
				AllowMoreArgs: true,
				Run:           Exec,
			},
		},
	},
}

func Help(_ []string, _ map[string]string, _ map[string]struct{}, _ []string) error {
	fmt.Println("usage: aqua [--quiet|<command>] [arguments]")
	fmt.Println()
	fmt.Println("commands:")
	fmt.Printf("  %-30s %s\n", "help", "print help")
	fmt.Printf("  %-30s %s\n", "version", "print version")
	fmt.Printf("  %-30s %s\n", "new <name> [flags]", "create new project")
	fmt.Printf("      %-32s %s\n", "--version=<...>", "project version")
	fmt.Printf("      %-32s %s\n", "--main=<...>.aq", "path to main file")
	fmt.Printf("      %-32s %s\n", "--path=<...>", "path to directory where the project will be created")
	fmt.Printf("      %-32s %s\n", "--lib", "is the project a library")
	fmt.Printf("  %-30s %s\n", "run [flags] [arguments]", "run project")
	fmt.Printf("      %-32s %s\n", "--project=<...>", "path to directory where the project is")
	fmt.Printf("      %-32s %s\n", "--file=<...>.aq", "path to file that will be executed")
	fmt.Printf("  %-30s %s\n", "exec <script>.aq [arguments]", "execute file")
	return nil
}

func Version(_ []string, _ map[string]string, _ map[string]struct{}, _ []string) error {
	fmt.Println("aqua v0.1.0-alpha")
	return nil
}

func New(vals []string, optionals map[string]string, flags map[string]struct{}, _ []string) error {
	if len(vals) < 1 {
		return fmt.Errorf("expected project name")
	}
	config := config.DefaultConfig()
	config.Name = vals[0]
	if version, ok := optionals["version"]; ok {
		config.Version = version
	}
	if main, ok := optionals["main"]; ok {
		if filepath.Ext(main) != "."+env.FILE_EXTENSION {
			main += "." + env.FILE_EXTENSION
		}
		config.Main = main
	}
	path := config.Name
	if p, ok := optionals["path"]; ok {
		path = p
	}
	if _, ok := flags["lib"]; ok {
		config.Lib = &struct{}{}
	}
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(path, env.CONFIG), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	err = toml.NewEncoder(file).Encode(config)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(path, config.Main), []byte("println(\"Hello World!\")\n"), 0644)
}

func Run(_ []string, optionals map[string]string, _ map[string]struct{}, args []string) error {
	project := "."
	if p, ok := optionals["project"]; ok {
		project = p
	}
	vm, err := vm.NewAndLoadConfig(filepath.Join(project, env.CONFIG), &vm.Paths[*object.Value]{}, &vm.VMs[*object.Value]{}, run.Run, args)
	if err != nil {
		return err
	}

	path := vm.Config.Main
	if file, ok := optionals["file"]; ok {
		path = file
	}

	_, err = run.Run(filepath.Join(project, path), "main", vm)
	if err != nil {
		return err
	}
	return nil
}

func Exec(vals []string, _ map[string]string, _ map[string]struct{}, args []string) error {
	if len(vals) < 1 {
		return fmt.Errorf("expected project name")
	}
	path := vals[0]
	vm, err := vm.New(config.DefaultConfig(), &vm.Paths[*object.Value]{}, &vm.VMs[*object.Value]{}, run.Run, args)
	if err != nil {
		return err
	}

	_, err = run.Run(path, "main", vm)
	if err != nil {
		return err
	}
	return nil

}

func Start() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()
	args := os.Args[1:]
	err := App.Run(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
