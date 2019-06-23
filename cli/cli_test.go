package cli

import (
	"errors"
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/options"
	"io"
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

func TestCLI_RunShell(t *testing.T) {
	c := New()
	_ = c.SetCommands(
		lime.Command{
			Keyword: "test",
			Func: func(_ []string) error {
				return nil
			},
		},
		lime.Command{
			Keyword: "error",
			Func: func(_ []string) error {
				return errors.New("failed successfully")
			},
		},
	)
	cc := c.(cli)

	// Capture the input and output
	inputReader, inputWriter, _ := os.Pipe()
	outputReader, outputWriter, _ := os.Pipe()
	os.Stdout = outputWriter
	os.Stdin = inputReader

	//Ensure the CLI enters and exits shell mode with no args
	os.Args = []string{"myCli"}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_ = c.Run()
		wg.Done()
	}()

	if readString(outputReader) != "entering shell mode\n> " {
		t.Error("the shell mode initialization output was unexpected")
	}

	//test empty input
	writeLine(inputWriter, "")
	if readString(outputReader) != "> " {
		t.Error("the shell mode new line was unexpected")
	}

	writeLine(inputWriter, "test")
	if readString(outputReader) != "> " {
		t.Error("the shell non-error, empty output new line was unexpected")
	}

	writeLine(inputWriter, "error")
	if readString(outputReader) != "failed successfully\n> " {
		t.Error("an error from a `lime.Func` was not output in the shell")
	}

	writeLine(inputWriter, "invalid")
	if readString(outputReader) != fmt.Sprintf("%s\n> ", errNoMatch.Error()) {
		t.Error("a `limeErrors.ErrNoMatch` was not output in the shell")
	}

	writeLine(inputWriter, *cc.exitWord)

	wg.Wait()
}

func TestCLI_RunNamedShell(t *testing.T) {
	c := New()
	c.SetName("myCli")
	cc := c.(cli)

	// Capture the input and output
	inputReader, inputWriter, _ := os.Pipe()
	outputReader, outputWriter, _ := os.Pipe()
	os.Stdout = outputWriter
	os.Stdin = inputReader

	//Ensure the CLI enters and exits shell mode with no args
	os.Args = []string{"myCli"}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		_ = c.Run()
		wg.Done()
	}()

	if readString(outputReader) != "entering shell mode for myCli\n> " {
		t.Error("the shell mode initialization output was unexpected")
	}

	writeLine(inputWriter, *cc.exitWord)

	wg.Wait()
}

func readString(reader io.Reader) string {
	bs := make([]byte, 128)

	n, _ := reader.Read(bs)
	return string(bs[:n])
}

func writeLine(writer io.Writer, s string) {
	_, _ = writer.Write([]byte(fmt.Sprintln(s)))
}
