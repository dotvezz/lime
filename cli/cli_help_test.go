package cli

import (
	"errors"
	"github.com/dotvezz/lime"
	"os"
	"testing"
)

func TestCLI_Help(t *testing.T) {
	_, outputWriter, _ := os.Pipe()
	os.Stderr = outputWriter

	c := New()
	_ = c.SetCommands(
		lime.Command{
			Keyword: "nested",
			Commands: []lime.Command{
				{
					Keyword:     "test",
					Description: "A test",
					Help:        "Used to do testing",
					Usage: []lime.Usage{
						{
							Example:     "myCli nested test",
							Explanation: "Does nothing",
						},
					},
				},
			},
		},
		lime.Command{
			Keyword: "noHelp",
			Func: func(_ []string) error {
				return errors.New("something")
			},
		},
		lime.Command{
			Description: "no keyword",
		},
	)

	// Capture the input and output
	inputReader, _, _ := os.Pipe()
	outputReader, outputWriter, _ := os.Pipe()
	os.Stdout = outputWriter
	os.Stderr = outputWriter
	os.Stdin = inputReader

	// Ensure the help for a command works
	{
		os.Args = []string{"myCli", "nested", "test", "--help"}
		err := c.Run()

		if err != nil {
			t.Error("the `Run` method returned an error for a command that should succeed")
		}

		expect := "A test\nUsed to do testing\n >  myCli nested test\n    Does nothing\n"
		if str := readString(outputReader); str != expect {
			t.Errorf("\nExpected: \n%s\nBut Got:\n%s\n", expect, str)
		}
	}

	// Ensure the help for the whole app works
	{
		os.Args = []string{"myCli", "--help"}
		err := c.Run()

		if err != nil {
			t.Error("the `Run` method returned an error for a command that should succeed")
		}

		expect := "Usage of myCli:\nnested test\n -  A test\n"
		if str := readString(outputReader); str != expect {
			t.Errorf("\nExpected: \n%s\nBut Got:\n%s\n", expect, str)
		}
	}

	// Ensure a command with no help fails properly
	{
		os.Args = []string{"myCli", "noHelp", "--help"}
		err := c.Run()

		if err != nil {
			t.Error("the `Run` method returned an error for a command that succeed")
		}

		expect := "No information for this command\n"
		if str := readString(outputReader); str != expect {
			t.Errorf("\nExpected: \n%s\nBut Got:\n%s\n", expect, str)
		}
	}
}
