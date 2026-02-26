package cli

import (
	"fmt"
	"strings"
)

type CLI struct {
	Default  *Args
	Commands map[string]CLI
}

type Args struct {
	Length        int
	Flags         map[string]bool // does it need a value or not
	AllowMoreArgs bool            // allow more args
	Run           func([]string, map[string]string, map[string]struct{}, []string) error
}

func (c CLI) Run(args []string) error {
	if len(args) <= 0 && c.Default != nil {
		return c.Default.Parse(args)
	}
	if len(args) <= 0 {
		return nil
	}
	if c, ok := c.Commands[args[0]]; ok {
		return c.Run(args[1:])
	}
	if c.Default != nil {
		return c.Default.Parse(args)
	}
	return nil
}

func (a Args) Parse(args []string) error {
	if len(args) < a.Length {
		return fmt.Errorf("expected at least %d arguments\n", a.Length)
	}
	values := args[:a.Length]
	args = args[a.Length:]
	flags := make(map[string]struct{})
	optionals := make(map[string]string)
	i := 0
	for ; i < len(args); i++ {
		v := args[i]
		if !strings.HasPrefix(v, "--") {
			break
		}
		v = strings.TrimPrefix(v, "--")
		name, value, ok := strings.Cut(v, "=")
		flag, found := a.Flags[name]
		if !found {
			args = args[i:]
			break
		}
		if flag && ok {
			optionals[name] = value
		} else if !flag && !ok {
			flags[name] = struct{}{}
		} else if flag && !ok {
			return fmt.Errorf("expected --%s=<value>, got --%s", name, name)
		} else if !flag && ok {
			return fmt.Errorf("expected --%s, got --%s", name, v)
		}
	}
	args = args[i:]
	if len(args) > 0 && !a.AllowMoreArgs {
		return fmt.Errorf("unexpected args remained: %s", strings.Join(args, " "))
	}

	return a.Run(values, optionals, flags, args)
}
