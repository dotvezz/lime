package cli

import (
	"bufio"
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/options"
	"os"
	"strings"
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
		return errNoInput
	}
	c, depth, err := match(*cli.commands, os.Args[1:], 1)
	if err != nil {
		return err
	}

	return cli.exec(c, depth, os.Args)
}

// match finds a matching command for a given set of arguments.
// Also returns the nesting depth of the matched command.
func match(commands []lime.Command, args []string, depth int) (*lime.Command, int, error) {
	var c *lime.Command
	for i := range commands {
		c = &commands[i]
		if c.Keyword == args[0] {
			if len(args) > 1 && len(c.Commands) > 0 {
				return match(c.Commands, args[1:], depth+1)
			}
			return c, depth, nil
		}
	}
	return nil, depth, errNoMatch
}

// shell launches the cli in an interactive shell
func (cli cli) shell() {
	fmt.Print("entering shell mode")
	if len(*cli.name) > 1 {
		fmt.Println(" for", *cli.name)
	} else {
		fmt.Println()
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf("%s ", *cli.prompt)
		scanner.Scan()
		input := scanner.Text()
		if string(input) == *cli.exitWord {
			break
		}
		args := strings.Split(string(input), " ")
		c, depth, err := match(*cli.commands, args, 0)
		if err != nil {
			if err != errNoMatch || len(input) > 0 {
				fmt.Println(err)
			}
			continue
		}
		err = cli.exec(c, depth, args)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// exec runs the function
func (cli cli) exec(c *lime.Command, depth int, args []string) error {
	if c.Func == nil {
		return errNoFunc
	}

	return c.Func(args[depth+1:])
}
