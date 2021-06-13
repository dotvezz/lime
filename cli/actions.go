package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/dotvezz/lime"
)

const (
	noInfo            = "No information for this command\n"
	explanationPrefix = "   "
	examplePrefix     = " > "
	descriptionPrefix = " - "
	argumentSeparator = " "
)

// exec runs the Func from a `lime.Command`
func exec(c *lime.Command, args []string, out io.Writer) error {
	if c.Func == nil {
		return errNoFunc
	}
	return c.Func(args, out)
}

func help(c *lime.Command) (string, error) {
	sb := new(strings.Builder)
	if len(c.Help) == 0 && len(c.Description) == 0 && len(c.Usage) == 0 {
		return noInfo, errNoHelp
	}

	if len(c.Description) > 0 {
		_, _ = fmt.Fprintln(sb, c.Description)
	}

	if len(c.Help) > 0 {
		_, _ = fmt.Fprintln(sb, c.Help)
	}

	for i := range c.Usage {
		_, _ = fmt.Fprintln(sb, examplePrefix, c.Usage[i].Example)
		_, _ = fmt.Fprintln(sb, explanationPrefix, c.Usage[i].Explanation)
	}

	return sb.String(), nil
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
func describeRecursively(c *lime.Command, args []string) string {
	sb := new(strings.Builder)
	keyword := strings.Trim(c.Keyword, " ")
	args = append(args, keyword)
	if len(keyword) == 0 {
		return ""
	}

	if len(c.Description) > 0 {
		_, _ = fmt.Fprintf(sb, "%s\n%s%s\n", strings.Join(args, argumentSeparator), descriptionPrefix, c.Description)
	}

	for _, com := range c.Commands {
		_, _ = sb.WriteString(describeRecursively(&com, args))
	}

	return sb.String()
}
