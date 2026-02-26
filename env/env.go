package env

import (
	"os"
	"syscall"

	"github.com/aqua-aq/aqua-core/pkg/fatal"
)

const FILE_EXTENSION = "aq"
const MAIN = "main." + FILE_EXTENSION
const CONFIG = "project.toml"

var AQUA_PATH string

func init() {
	var ok bool
	AQUA_PATH, ok = syscall.Getenv("AQUA_PATH")
	if !ok {
		fatal.Fatal("environment value AQUA_PATH not found")
		os.Exit(0)
	}
}

const LIB_PATH = "libs"
const VERSION = "aqua v0.1.3-beta"
