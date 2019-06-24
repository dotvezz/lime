package cli

import (
	"errors"
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/options"
	"os"
	"testing"
)

func TestCLI_Run_Basic(t *testing.T) {
	_, outputWriter, _ := os.Pipe()
	os.Stderr = outputWriter

	c := New()
	_ = c.SetCommands(
		lime.Command{
			Keyword: "fail",
			Func: func(_ []string) error {
				return errors.New("failed successfully")
			},
		},
		lime.Command{
			Keyword: "succeed",
			Func: func(_ []string) error {
				return nil
			},
		},
		lime.Command{
			Keyword: "repeat",
			Func: func(args []string) error {
				fmt.Println(args)
				return nil
			},
		},
		lime.Command{
			Keyword: "noFunc",
		},
		lime.Command{
			Keyword: "nested",
			Commands: []lime.Command{
				{
					Keyword: "test",
					Func: func(_ []string) error {
						fmt.Println("success")
						return nil
					},
				},
			},
		},
	)

	// Ensure the succeeding command behaves as expected
	{
		os.Args = []string{"myCli", "succeed"}
		err := c.Run()

		if err != nil {
			t.Error("the `Run` method returned an error for a command that should succeed")
		}
	}

	// Ensure the failing command behaves as expected
	{
		os.Args = []string{"myCli", "fail"}
		err := c.Run()

		if err == nil {
			t.Error("the `Run` method did not return an error for the failing command")
		}

		if err != nil && err.Error() != "failed successfully" {
			t.Error("the `Run` method returned the wrong error for the failing command")
		}
	}

	// Ensure an error is returned when there is no func for the associated command
	{
		os.Args = []string{"myCli", "noFunc"}
		err := c.Run()

		if err == nil {
			t.Error("the `Run` method did not return an error for the command with a nil func")
		}

		if err != nil && err.Error() != errNoFunc.Error() {
			t.Error("the `Run` method returned the wrong error for the command with a nil func")
		}
	}

	// Ensure an error is returned when there is no command to run and interactive shell is disabled
	{
		_ = c.SetOptions(options.NoShell)
		os.Args = []string{"myCli"}
		err := c.Run()

		if err == nil {
			t.Error("the `Run` method did not return an error with no input an shell disabled")
		}

		if err != nil && err.Error() != errNoInput.Error() {
			t.Error("the `Run` method returned the wrong error for the command with a nil func")
		}
	}

	// Ensure an error is returned when there is no match found
	{
		os.Args = []string{"myCli", "invalid"}
		err := c.Run()

		if err == nil {
			t.Error("the `Run` method did not return an error when there should be no matching command")
		}

		if err != nil && err.Error() != errNoMatch.Error() {
			t.Error("the `Run` method returned the wrong error when there should be no matching command")
		}
	}
}

func TestCLI_Run_Captured_IO(t *testing.T) {
	_, outputWriter, _ := os.Pipe()
	os.Stderr = outputWriter

	c := New()
	_ = c.SetCommands(
		lime.Command{
			Keyword: "fail",
			Func: func(_ []string) error {
				return errors.New("failed successfully")
			},
		},
		lime.Command{
			Keyword: "succeed",
			Func: func(_ []string) error {
				return nil
			},
		},
		lime.Command{
			Keyword: "repeat",
			Func: func(args []string) error {
				fmt.Println(args)
				return nil
			},
		},
		lime.Command{
			Keyword: "noFunc",
		},
		lime.Command{
			Keyword: "nested",
			Commands: []lime.Command{
				{
					Keyword: "test",
					Func: func(_ []string) error {
						fmt.Println("success")
						return nil
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
		os.Args = []string{"myCli", "repeat", "blah"}
		err := c.Run()

		if err != nil {
			t.Error("the `Run` method returned an error for a command that should succeed")
		}

		if readString(outputReader) != fmt.Sprintln("[blah]") {
			t.Error("the `Run` command ran but its output was unexpected")
		}
	}

	// Ensure nested commands work
	{
		os.Args = []string{"myCli", "nested", "test"}
		err := c.Run()

		if err != nil {
			t.Error("the `Run` method returned an error for the nested command")
		}

		if readString(outputReader) != fmt.Sprintln("success") {
			t.Error("the `Run` method did not run the nested command")
		}
	}
}
