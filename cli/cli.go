package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/options"
)

// CLI is the private struct which holds pointers to the CLI's internal values
type CLI struct {
	options  lime.Option
	commands []lime.Command
	name     string
	prompt   string
	exitWord string
	out      io.Writer
	in       io.Reader
	err      io.Writer
}

// New creates a new CLI
func New() *CLI {
	defaultPrompt := ">"
	defaultExitWord := "exit"
	return &CLI{
		commands: make([]lime.Command, 0),
		prompt:   defaultPrompt,
		exitWord: defaultExitWord,
		out:      os.Stdout,
		in:       os.Stdin,
	}
}

// SetOptions takes a variadic list of options and applies them to the CLI
// Returns an error if any of the options given are not valid
func (cli *CLI) SetOptions(opts ...lime.Option) error {
	for _, option := range opts {
		if options.IsValid(option) {
			return errInvalidOption
		}
		cli.options |= option
	}

	return nil
}

// SetCommands takes a variadic list of Commands and stores them in the CLI
func (cli *CLI) SetCommands(commands ...lime.Command) error {
	cli.commands = append(cli.commands, commands...)
	return nil
}

// SetName takes a string as the CLI application's name, used in some out
func (cli *CLI) SetName(name string) {
	cli.name = name
}

// SetPrompt takes a string as the CLI application's prompt, used in interactive mode
func (cli *CLI) SetPrompt(prompt string) {
	cli.prompt = prompt
}

// SetExitWord takes a string as the keyword to exit the interactive
func (cli *CLI) SetExitWord(exitWord string) {
	cli.exitWord = exitWord
}

// SetOutput takes an io.Writer and treats it as the output stream target (Default is os.Stdout)
func (cli *CLI) SetOutput(w io.Writer) {
	cli.out = w
}

// SetInput takes an io.Reader and treats it as the input stream target (Default is os.Stdin)
func (cli *CLI) SetInput(w io.Reader) {
	cli.in = w
}

// SetErrOutput takes an io.Writer and treats it as the error output stream target (Default is os.Stderr)
func (cli *CLI) SetErrOutput(w io.Writer) {
	cli.err = w
}

// Run finds a matching Command for the arguments given and invokes its Func.
func (cli CLI) Run(args ...string) error {
	fromOS := false
	if len(args) == 0 {
		args = os.Args[1:]
		fromOS = true
	}
	// Go to interactive mode if it's not disabled and there are no args
	if len(args) == 1 && fromOS || len(args) == 0 {
		if cli.options&options.NoInteractiveMode == 0 {
			cli.interactive()
		}
		if cli.options&options.PrintErrors > 0 {
			_, _ = fmt.Fprintln(cli.err, errNoInput.Error())
		}
		return errNoInput
	}
	c, depth, err := match(cli.commands, args, 1)

	// Custom flag.Usage for extended help out
	flag.CommandLine = flag.NewFlagSet(args[0], flag.ContinueOnError)
	flag.Usage = func() {
		var helpStr string
		if err == nil {
			helpStr, _ = help(c)
		} else {
			helpStr = cli.help()
		}
		fmt.Fprint(cli.out, helpStr)
	}

	flag.CommandLine.Parse(args)
	if triggerHelp(args) {
		flag.Usage()
		return nil
	}

	if err != nil {
		if cli.options&options.PrintErrors > 0 {
			_, _ = fmt.Fprintln(cli.err, err.Error())
		}
		return err
	}

	err = exec(c, depth, args, cli.out)
	if err != nil {
		if cli.options&options.PrintErrors > 0 {
			_, _ = fmt.Fprintln(cli.err, err.Error())
		}
	}
	return err
}

// interactive launches the program in interactive mode
func (cli CLI) interactive() {
	sb := &strings.Builder{}

	sb.WriteString("entering interactive mode")
	if len(cli.name) > 1 {
		_, _ = fmt.Fprintf(sb, " for %s\n", cli.name)
	} else {
		_, _ = fmt.Fprintln(sb)
	}

	_, err := fmt.Fprint(cli.out, sb.String())
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(cli.in)
	for {
		_, err = fmt.Fprintf(cli.out, "%s ", cli.prompt)
		if err != nil {
			panic(err)
		}
		scanner.Scan()
		input := scanner.Text()
		if input == cli.exitWord {
			break
		}
		args := strings.Split(input, " ")
		c, depth, err := match(cli.commands, args, 0)
		if err != nil {
			if err != errNoMatch || len(input) > 0 {
				_, _ = fmt.Fprintln(cli.out, err)
			}
			continue
		}
		err = exec(c, depth, args, cli.out)
		if err != nil {
			_, _ = fmt.Fprintln(cli.out, err)
		}
	}
}
