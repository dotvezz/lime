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

func help(c *lime.Command) error {
	if len(c.Help) == 0 && len(c.Description) == 0 && len(c.Usage) == 0 {
		fmt.Println("No information for this command")
		return errNoHelp
	}

	if len(c.Description) > 0 {
		fmt.Println(c.Description)
	}

	if len(c.Help) > 0 {
		fmt.Println(c.Help)
	}

	for i := range c.Usage {
		fmt.Println(" > ", c.Usage[i].Example)
		fmt.Println("   ", c.Usage[i].Explanation)
	}

	return nil
}

// triggerHelp checks the args for any of the help flags. Returns true if there was a help flag, false otherwise
func triggerHelp(args []string) bool {
	for i := range args {
		if b, ok := helpFlags[args[i]]; ok {
			return b
		}
	}

	return false
}

// traverses a given command's Commands field, and any tree sprouting from it, to generate help for each command
func describeRecursively(c *lime.Command, args []string) {
	keyword := strings.Trim(c.Keyword, " ")
	args = append(args, keyword)
	if len(keyword) == 0 {
		return
	}
	if len(c.Description) > 0 {
		fmt.Println(strings.Join(args, " "))
		fmt.Println(" - ", c.Description)
	}

	for _, com := range c.Commands {
		describeRecursively(&com, args)
	}
}
