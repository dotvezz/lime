package cli

import (
	"fmt"
	"github.com/dotvezz/lime"
	"strings"
)

// exec runs the Func from a `lime.Command`
func exec(c *lime.Command, depth int, args []string) error {
	if c.Func == nil {
		return errNoFunc
	}

	return c.Func(args[depth+1:])
}

func help(c *lime.Command, withUsage bool, args []string) error {
	if len(c.Help) == 0 {
		return errNoHelp
	}

	fmt.Println(c.Description)
	fmt.Println(c.Help)

	if !withUsage && len(c.Usage) > 0 {
		fmt.Printf("\nFor examples about this command, try %s --usage\n", strings.Join(args, " "))
	}

	return nil
}

func usage(c *lime.Command, withHelp bool, args []string) error {
	if len(c.Usage) == 0 {
		return errNoUsage
	}
	fmt.Println(c.Description)

	for _, usage := range c.Usage {
		fmt.Println(usage.Example)
		fmt.Println(usage.Explanation)
		fmt.Println()
	}

	if !withHelp && len(c.Usage) > 0 {
		fmt.Printf("\nFor help about this command, try %s --help\n", strings.Join(args, " "))
	}

	return nil
}
