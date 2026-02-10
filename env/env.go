package env

import "syscall"

const FILE_EXTENSION = "aq"
const MAIN = "main." + FILE_EXTENSION
const CONFIG = "project.toml"

var AQUA_PATH string

func init() {
	var ok bool
	AQUA_PATH, ok = syscall.Getenv("AQUA_PATH")
	if !ok {
		panic("environment value AQUA_PATH not found ")
	}
}
