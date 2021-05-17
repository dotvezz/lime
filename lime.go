package lime

import "io"

// Command defines the structure of a cli command.
type Command struct {
	// The keyword which invokes this command
	Keyword string
	// A brief description of the command, used in all --help output
	Description string
	// A collection of examples and explanations for the command, used in command-specific --usage output
	Usage []Usage
	// A helpful bit of information about the command, used in all --help output
	Help string
	// Nested commands
	Commands []Command
	// The function to run when this command is invoked
	Func Func
}

// Usage defines the structure of a Usage entry
type Usage struct {
	// The example input
	Example string
	// The explanation of the example
	Explanation string
}

// Func is the signature of a function to run when a Command is invoked.
type Func func(args []string, out io.Writer) error

// Option is a bit mask value for setting options on a CLI
type Option int64
