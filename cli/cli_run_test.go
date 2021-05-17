package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/dotvezz/lime"
)

func TestCLI_Run_Basic(t *testing.T) {
	_, outputWriter, _ := os.Pipe()
	os.Stderr = outputWriter

	c := New()
	_ = c.SetCommands(
		lime.Command{
			Keyword: "fail",
			Func: func(_ []string, _ io.Writer) error {
				return errors.New("failed successfully")
			},
		},
		lime.Command{
			Keyword: "succeed",
			Func: func(_ []string, _ io.Writer) error {
				return nil
			},
		},
		lime.Command{
			Keyword: "repeat",
			Func: func(args []string, w io.Writer) error {
				fmt.Fprintln(w, args)
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
					Func: func(_ []string, w io.Writer) error {
						fmt.Fprintln(w, "success")
						return nil
					},
				},
			},
		},
	)

	// Ensure the succeeding command behaves as expected
	{
		err := c.Run("succeed")

		if err != nil {
			t.Error("the `Run` method returned an error for a command that should succeed")
		}
	}

	// Ensure the failing command behaves as expected
	{
		err := c.Run("fail")

		if err == nil {
			t.Error("the `Run` method did not return an error for the failing command")
		}

		if err != nil && err.Error() != "failed successfully" {
			t.Error("the `Run` method returned the wrong error for the failing command")
		}
	}

	// Ensure an error is returned when there is no func for the associated command
	{
		err := c.Run("noFunc")

		if err == nil {
			t.Error("the `Run` method did not return an error for the command with a nil func")
		}

		if err != nil && err.Error() != errNoFunc.Error() {
			t.Error("the `Run` method returned the wrong error for the command with a nil func")
		}
	}

	//// Ensure an error is returned when there is no command to run and interactive is disabled
	//{
	//	_ = c.SetOptions(options.NoInteractiveMode)
	//	err := c.Run()
	//
	//	if err == nil {
	//		t.Error("the `Run` method did not return an error with no input an interactive disabled")
	//	}
	//
	//	if err != nil && err.Error() != errNoInput.Error() {
	//		t.Error("the `Run` method returned the wrong error for the command with a nil func")
	//	}
	//}

	// Ensure an error is returned when there is no match found
	{
		err := c.Run("invalid")

		if err == nil {
			t.Error("the `Run` method did not return an error when there should be no matching command")
		}

		if err != nil && err.Error() != errNoMatch.Error() {
			t.Error("the `Run` method returned the wrong error when there should be no matching command")
		}
	}
}

func TestCLI_Run_Captured_IO(t *testing.T) {
	c := New()
	_ = c.SetCommands(
		lime.Command{
			Keyword: "fail",
			Func: func(_ []string, _ io.Writer) error {
				return errors.New("failed successfully")
			},
		},
		lime.Command{
			Keyword: "succeed",
			Func: func(_ []string, _ io.Writer) error {
				return nil
			},
		},
		lime.Command{
			Keyword: "repeat",
			Func: func(args []string, out io.Writer) error {
				fmt.Fprintln(out, args)
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
					Func: func(_ []string, out io.Writer) error {
						fmt.Fprintln(out, "success")
						return nil
					},
				},
			},
		},
	)
	outBuffer := &bytes.Buffer{}
	errBuffer := &bytes.Buffer{}

	c.SetOutput(outBuffer)
	c.SetErrOutput(errBuffer)

	// Ensure the command out and injected args behave as expected
	{
		outBuffer.Reset()
		err := c.Run("repeat", "blah")

		if err != nil {
			t.Error("the `Run` method returned an error for a command that should succeed")
		}

		if out := outBuffer.String(); out != fmt.Sprintln("[blah]") {
			t.Error("the `Run` command ran but its out was unexpected")
		}
	}

	// Ensure nested commands work
	{
		outBuffer.Reset()
		err := c.Run("nested", "test")

		if err != nil {
			t.Error("the `Run` method returned an error for the nested command")
		}

		if outBuffer.String() != fmt.Sprintln("success") {
			t.Error("the `Run` method did not run the nested command")
		}
	}
}
