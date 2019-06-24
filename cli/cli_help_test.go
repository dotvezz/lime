package cli

import (
	"fmt"
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
					Keyword: "test",
					Description: "A test",
					Help: "Used to do testing",
					Usage: []lime.Usage{
						{
							Example: "myCli nested test",
							Explanation: "Does nothing",
						},
					},
				},
			},
		},
	)

	// Capture the input and output
	inputReader, _, _ := os.Pipe()
	outputReader, outputWriter, _ := os.Pipe()
	os.Stdout = outputWriter
	os.Stderr = outputWriter
	os.Stdin = inputReader

	// Ensure the command output and injected args behave as expected
	{
		os.Args = []string{"myCli", "nested", "test", "--help"}
		err := c.Run()

		if err != nil {
			t.Error("the `Run` method returned an error for a command that should succeed")
		}

		if str := readString(outputReader); str != fmt.Sprintln("[blah]") {
			t.Error("the `Run` command ran but its output was unexpected:", str)
		}
	}
}