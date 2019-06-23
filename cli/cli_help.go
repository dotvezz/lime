package cli

import (
	"fmt"
	"github.com/dotvezz/lime"
	"strings"
)

func (cli cli) help() {
	for _, c := range *cli.commands {
		describeRecursively(&c, make([]string, 0))
	}
}

func describeRecursively(c *lime.Command, args []string) {
	keyword := strings.Trim(c.Keyword, " ")
	args = append(args, keyword)
	if len(keyword) == 0 {
		return
	}
	fmt.Println(strings.Join(args, " "))
	if len(c.Description) > 0 {
		fmt.Print(" - ", c.Description)
	}

	for _, c := range c.Commands {
		describeRecursively(&c, args)
	}
}