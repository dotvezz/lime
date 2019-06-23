package cli

import "github.com/dotvezz/lime"

// exec runs the Func from a `lime.Command`
func exec(c *lime.Command, depth int, args []string) error {
	if c.Func == nil {
		return errNoFunc
	}

	return c.Func(args[depth+1:])
}
