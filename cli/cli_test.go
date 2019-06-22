package cli

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/options"
	"io"
	"log"
	"os"
	"sync"
	"testing"
)

func TestCLI_New(t *testing.T) {
	cc, ok := New().(cli)

	// Ensure that we're getting the expected internal implementation
	if !ok {
		t.Error("the `New` function did not return a `cli`")
	}

	// Ensure that all the pointers are initialized to protect against internal nil dereferences
	if cc.options == nil {
		t.Error("the `New` function returned a `cli` with nil options")
	}
	if cc.commands == nil {
		t.Error("the `New` function returned a `cli` with nil `commands`")
	}
	if cc.name == nil {
		t.Error("the `New` function returned a `cli` with nil name")
	}
	if cc.prompt == nil {
		t.Error("the `New` function returned a `cli` with nil prompt")
	}
	if cc.exitWord == nil {
		t.Error("the `New` function returned a `cli` with nil exitWord")
	}
}

func TestCLI_SetOptions(t *testing.T) {
	// Ensure we can set a single option
	{
		c := New()
		cc := c.(cli)

		_ = c.SetOptions(options.NoShell)
		if *cc.options&options.NoShell == 0 {
			t.Error("the `SetOptions` method did not save the option given")
		}
	}

	// Ensure we can set multiple options
	{
		c := New()
		cc := c.(cli)
		_ = c.SetOptions(1, 2, 4)
		if *cc.options&7 != 7 {
			t.Error("the `SetOptions` method did not save multiple options given")
		}
	}

	// Ensure invalid options are rejected
	{
		c := New()
		err := c.SetOptions(7)
		if err == nil {
			t.Error("the `SetOptions` method did not return an error when given invalid input")
		}
	}
}

func TestCLI_SetCommands(t *testing.T) {
	// Ensure we can set a single option
	{
		c := New()
		cc := c.(cli)
		_ = c.SetCommands(
			lime.Command{
				Keyword: "test",
			},
		)

		if len(*cc.commands) != 1 {
			t.Error("the `SetCommands` method did not save the command given")
		}

		if (*cc.commands)[0].Keyword != "test" {
			t.Error("the `SetCommands` method saved a command but its keyword value was lost")
		}
	}

	// Ensure we can set multiple options
	{
		c := New()
		cc := c.(cli)
		_ = c.SetCommands(
			lime.Command{
				Keyword: "test1",
			},
			lime.Command{
				Keyword: "test2",
			},
		)

		if len(*cc.commands) != 2 {
			t.Error("the `SetCommands` method did not save the commands given")
		}

		if (*cc.commands)[0].Keyword != "test1" {
			t.Error("the `SetCommands` method saved the first command but its keyword value was lost")
		}

		if (*cc.commands)[1].Keyword != "test2" {
			t.Error("the `SetCommands` method saved the second command but its keyword value was lost")
		}
	}
}

func TestCLI_SetName(t *testing.T) {
	c := New()
	cc := c.(cli)
	c.SetName("test")
	if *cc.name != "test" {
		t.Error("the `SetName` method did not save the name")
	}
}

func TestCLI_SetPrompt(t *testing.T) {
	c := New()
	cc := c.(cli)
	c.SetPrompt("test")
	if *cc.prompt != "test" {
		t.Error("the `SetPrompt` method did not save the prompt")
	}
}

func TestCLI_SetExitWord(t *testing.T) {
	c := New()
	cc := c.(cli)
	c.SetExitWord("test")
	if *cc.exitWord != "test" {
		t.Error("the `SetExitWord` method did not save the exit word")
	}
}

func TestCLI_Run(t *testing.T) {
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
	)

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

	// Ensure the succeeding command behaves as expected
	{
		os.Args = []string{"myCli", "succeed"}
		err := c.Run()

		if err != nil {
			t.Error("the `Run` returned an error for the command that should succeed")
		}
	}

	// Ensure the command output and injected args behave as expected
	{
		os.Args = []string{"myCli", "repeat", "blah"}
		output, err := captureOutput(c.Run)

		if err != nil {
			t.Error("the `Run` returned an error for a command that should succeed")
		}

		if output != fmt.Sprintln("[blah]") {
			t.Error("the `Run` command ran but its output was unexpected")
		}
	}
}

func TestCLI_RunShell(t *testing.T) {
	t.Skip("not implemented yet")
	c := New()
	_ = c.SetCommands(
		lime.Command{
			Keyword: "test",
			Func: func(_ []string) error {
				return nil
			},
		},
	)

	// Ensure the CLI enters and exits shell mode with no args
	//os.Args = []string{"myCli"}
}

// Modified from https://medium.com/@hau12a1/golang-capturing-log-println-and-fmt-println-output-770209c791b4
func captureOutput(f func() error) (string, error) {
	reader, writer, _ := os.Pipe()
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	err := f()
	_ = writer.Close()
	return <-out, err
}
