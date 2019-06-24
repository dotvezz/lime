package cli

import (
	"flag"
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/options"
	"os"
)

// cli is the private struct which holds pointers to the cli's internal values
type cli struct {
	options  *lime.Option
	commands *[]lime.Command
	name     *string
	prompt   *string
	exitWord *string
}

// New creates a new CLI
func New() lime.CLI {
	defaultPrompt := ">"
	defaultExitWord := "exit"
	return cli{
		commands: new([]lime.Command),
		options:  new(lime.Option),
		name:     new(string),
		prompt:   &defaultPrompt,
		exitWord: &defaultExitWord,
	}
}

// SetOptions takes a variadic list of options and applies them to the cli
// Returns an error if any of the options given are not valid
func (cli cli) SetOptions(opts ...lime.Option) error {
	for _, option := range opts {
		if options.IsValid(option) {
			return errInvalidOption
		}
		*cli.options |= option
	}

	return nil
}

// SetCommands takes a variadic list of Commands and stores them in the cli
func (cli cli) SetCommands(commands ...lime.Command) error {
	*cli.commands = append(*cli.commands, commands...)
	return nil
}

// SetName takes a string as the CLI application's name, used in some output
func (cli cli) SetName(name string) {
	*cli.name = name
}

// SetPrompt takes a string as the CLI application's prompt, used in interactive shell mode
func (cli cli) SetPrompt(prompt string) {
	*cli.prompt = prompt
}

// SetExitWord takes a string as the keyword to exit the interactive shell
func (cli cli) SetExitWord(exitWord string) {
	*cli.exitWord = exitWord
}

// Run finds a matching Command for the arguments given and invokes its Func.
func (cli cli) Run() error {
	// Go to shell mode if it's not disabled and there are no args
	if len(os.Args) == 1 {
		if *cli.options&options.NoShell == 0 {
			cli.shell()
		}
		_, _ = fmt.Fprintln(os.Stderr, errNoInput.Error())
		return errNoInput
	}
	c, depth, err := match(*cli.commands, os.Args[1:], 1)

	// Custom flag.Usage for extended help output
	flag.Usage = func() {
		if err == nil {
			_ = help(c, os.Args)
		} else {
			cli.help()
		}
		os.Exit(0)
	}
	flag.Parse()
	if triggerHelp(os.Args) {
		flag.Usage()
	}

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		return err
	}

	err = exec(c, depth, os.Args)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}
	return err
}
