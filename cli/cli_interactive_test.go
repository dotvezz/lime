package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/dotvezz/lime"
)

func TestCLI_RunInteractive(t *testing.T) {
	c := New()
	_ = c.SetCommands(
		lime.Command{
			Keyword: "test",
			Func: func(_ []string, _ io.Writer) error {
				return nil
			},
		},
		lime.Command{
			Keyword: "error",
			Func: func(_ []string, _ io.Writer) error {
				return errors.New("failed successfully")
			},
		},
	)

	output, out, _ := os.Pipe()
	in, input, _ := os.Pipe()

	c.SetOutput(out)
	c.SetInput(in)

	//Ensure the CLI enters and exits interactive mode with no args
	os.Args = []string{"myCli"}

	go func() {
		c.Run()
	}()

	if err := assertReadString("entering interactive mode\n> ", output); err != nil {
		t.Error(err)
	}

	_, _ = fmt.Fprintln(input, "")
	if err := assertReadString("> ", output); err != nil {
		t.Error(err)
	}

	_, _ = fmt.Fprintln(input, "test")
	if err := assertReadString("> ", output); err != nil {
		t.Error(err)
	}

	_, _ = fmt.Fprintln(input, "error")
	if err := assertReadString("failed successfully\n> ", output); err != nil {
		t.Error(err)
	}

	_, _ = fmt.Fprintln(input, "invalid")
	if err := assertReadString(fmt.Sprintf("%s\n> ", errNoMatch.Error()), output); err != nil {
		t.Error(err)
	}

	_, _ = fmt.Fprintln(input, c.exitWord)
}

func TestCLI_RunNamedInteractive(t *testing.T) {
	c := New()
	c.SetName("myCli")

	output, out, _ := os.Pipe()
	in, input, _ := os.Pipe()

	c.SetOutput(out)
	c.SetInput(in)

	go func() {
		c.Run()
	}()

	//Ensure the CLI enters and exits interactive mode with no args
	os.Args = []string{"myCli"}

	go func() {
		_ = c.Run()
	}()

	if err := assertReadString("entering interactive mode for myCli\n> ", output); err != nil {
		t.Error(err)
	}

	fmt.Fprintln(input, c.exitWord)

}

func assertReadString(want string, r io.Reader) (err error) {
	time.Sleep(100 * time.Millisecond)
	bs := make([]byte, len([]byte(want)))
	_, _ = r.Read(bs)
	got := string(bs)
	if got != want {
		return fmt.Errorf(`wanted "%s", got "%s"`, want, got)
	}

	return nil
}
