package cli

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/dotvezz/lime"
)

func TestCLI_Help(t *testing.T) {
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
			Func: func(_ []string, _ io.Writer) error {
				return errors.New("something")
			},
		},
		lime.Command{
			Description: "no keyword",
		},
	)
	buffer := &bytes.Buffer{}
	c.SetOutput(buffer)

	//// Ensure the help for a command works
	//{
	//	buffer.Reset()
	//	err := c.Run("myCli", "nested", "test", "--help")
	//
	//	if err != nil {
	//		t.Error("the `Run` method returned an error for a command that should succeed")
	//	}
	//
	//	expect := "A test\nUsed to do testing\n >  myCli nested test\n    Does nothing\n"
	//	if str := buffer.String(); str != expect {
	//		t.Errorf("\nExpected: \n%s\nBut Got:\n%s\n", expect, str)
	//	}
	//}
	//
	//// Ensure the help for the whole app works
	//{
	//	buffer.Reset()
	//	err := c.Run("--help")
	//
	//	if err != nil {
	//		t.Error("the `Run` method returned an error for a command that should succeed")
	//	}
	//
	//	expect := "Usage of myCli:\nnested test\n - A test\n"
	//	if str := buffer.String(); str != expect {
	//		t.Errorf("\nExpected: \n%s\nBut Got:\n%s\n", expect, str)
	//	}
	//}

	// Ensure a command with no help fails properly
	{
		buffer.Reset()
		err := c.Run("noHelp", "--help")

		if err != nil {
			t.Error("the `Run` method returned an error for a command that succeed")
		}

		expect := "No information for this command\n"
		if str := buffer.String(); str != expect {
			t.Errorf("\nExpected: \n%s\nBut Got:\n%s\n", expect, str)
		}
	}
}
