package cli

import (
	"bufio"
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/errors"
	"github.com/dotvezz/lime/options"
	"os"
	"strings"
)

// The private struct which holds pointers to the cli's internal values
type cli struct {
	options  *lime.Option
	commands *[]lime.Command
	appName  *string
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
		appName:  new(string),
		prompt:   &defaultPrompt,
		exitWord: &defaultExitWord,
	}
}

// SetOptions takes a variadic list of options and applies them to the cli
// Returns an error if any of the options given are not valid
func (cli cli) SetOptions(opts ...lime.Option) error {
	for _, option := range opts {
		if options.IsValid(option) {
			return errors.InvalidOption
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

func (cli cli) SetName(name string) {
	*cli.appName = name
}

func (cli cli) SetPrompt(prompt string) {
	*cli.prompt = prompt
}

func (cli cli) SetExitWord(exitWord string) {
	*cli.exitWord = exitWord
}

// Run finds a matching Command for the arguments given and invokes its Func.
func (cli cli) Run() error {
	// Go to shell mode if it's not disabled and there are no args
	if len(os.Args) == 1 {
		if *cli.options&options.NoShell == 0 {
			return cli.shell()
		} else {
			return errors.NoInput
		}
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
	if len(args) == 0 {
		return nil, depth, errors.NoMatch
	}

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
	return nil, depth, errors.NoMatch
}

// shell launches the cli in an interactive shell
func (cli cli) shell() error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("entering shell mode")
	if len(*cli.appName) > 1 {
		fmt.Println("for", *cli.appName)
	} else {
		fmt.Println()
	}

	for {
		fmt.Printf("%s ", *cli.prompt)
		input, _ := reader.ReadString('\n')
		input = strings.Replace(input, "\n", "", -1)
		if input == *cli.exitWord {
			break
		}
		args := strings.Split(input, " ")
		c, depth, err := match(*cli.commands, args, 0)
		if err != nil {
			if err != errors.NoMatch || len(input) > 0 {
				fmt.Println(err)
			}
			continue
		}
		err = cli.exec(c, depth, args)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

// exec runs the function
func (cli cli) exec(c *lime.Command, depth int, args []string) error {
	if c.Func == nil {
		return errors.NoFunc
	}

	return c.Func(args[depth+1:])
}
