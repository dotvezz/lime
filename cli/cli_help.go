package cli

import (
	"fmt"
	"strings"
)

var helpFlags = map[string]bool{
	"-h":      true,
	"--help":  true,
	"--usage": true,
}

func (cli cli) help() string {
	sb := new(strings.Builder)
	for _, c := range *cli.commands {
		_, _ = fmt.Fprint(sb, describeRecursively(&c, make([]string, 0)))
	}
	return sb.String()
}
