package cli

import (
	"github.com/dotvezz/lime"
)

// match finds a matching command for a given set of arguments.
// Also returns the nesting depth of the matched command.
func match(commands []lime.Command, args []string, depth int) (*lime.Command, int, error) {
	var c *lime.Command
	for i := range commands {
		c = &commands[i]
		if c.Keyword == args[0] {
			if len(args) > 1 && len(c.Commands) > 0 && len(args) > 1 {
				return match(c.Commands, args[1:], depth+1)
			}
			return c, depth, nil
		}
	}

	return nil, depth, errNoMatch
}
