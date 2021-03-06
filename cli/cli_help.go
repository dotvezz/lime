package cli

import (
	"strings"
)

var helpFlags = map[string]bool{
	"-h":      true,
	"--help":  true,
	"--usage": true,
}

func (cli CLI) help() string {
	sb := new(strings.Builder)
	for _, c := range cli.commands {
		_, _ = sb.WriteString(describeRecursively(&c, make([]string, 0)))
	}
	return sb.String()
}
