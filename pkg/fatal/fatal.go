package fatal

import (
	"fmt"
	"os"
)

func Fatal(args ...any) {
	fmt.Fprintln(os.Stderr, "\033[31m"+fmt.Sprint(args...)+"\033[0m")
}
