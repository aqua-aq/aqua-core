package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/vandi37/aqua/source/config"
)

func main() {
	// path := "test.aq"
	// vm, err := vm.New[*object.Value](".")
	// if err != nil {
	// 	panic(err)
	// }
	// _, err = run.Run(path, "main", vm)
	// if err != nil {
	// 	panic(err)
	// }
	file, err := os.Open("config.toml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	config, err := config.NewConfig(file)
	if err != nil {
		panic(err)
	}
	spew.Dump(config)
}
