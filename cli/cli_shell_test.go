package cli

import (
	"errors"
	"fmt"
	"github.com/dotvezz/lime"
	"os"
	"sync"
	"testing"
)

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
		c.Run()
		wg.Done()
	}()

	if readString(outputReader) != "entering shell mode\n> " {
		t.Error("the shell mode initialization output was unexpected:")
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
