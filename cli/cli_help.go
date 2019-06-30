package cli

import (
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
		sb.WriteString(describeRecursively(&c, make([]string, 0)))
	}
	return sb.String()
}
