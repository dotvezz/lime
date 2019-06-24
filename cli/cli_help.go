package cli

var helpFlags = map[string]bool{
	"-h":      true,
	"--help":  true,
	"--usage": true,
}

func (cli cli) help() {
	for _, c := range *cli.commands {
		describeRecursively(&c, make([]string, 0))
	}
}
