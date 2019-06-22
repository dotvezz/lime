package lime

// CLI is the public interface, which is exposed to packages consuming lime/cli
type CLI interface {
	SetOptions(opts ...Option) error
	SetCommands(commands ...Command) error
	SetName(name string)
	Run() error
	SetPrompt(prompt string)
	SetExitWord(exitWord string)
}

// Command defines the structure of a cli command.
type Command struct {
	// The keyword which invokes this command
	Keyword string
	// A brief description of the command, used in all --help and --usage output
	Description string
	// A collection of examples and explanations for the command, used in command-specific --usage output
	Examples []Example
	// A brief description of the command, used in command-specific --usage and --help output
	Help string
	// Nested commands
	Commands []Command
	// The function to run when this command is invoked
	Func CommandFunc
}

// Example defines the structure of a Examples entry
type Example struct {
	// The example input
	Input string
	// The explanation of the example
	Explanation string
}

// CommandFunc is the signature of a function to run when a Command is invoked.
type CommandFunc func(args []string) error

// Option is a bit mask value for setting options on a CLI
type Option int64
