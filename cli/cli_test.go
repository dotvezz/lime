package cli

import (
	"testing"

	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/options"
)

func TestCLI_SetOptions(t *testing.T) {
	// Ensure we can set a single option
	{
		c := New()

		_ = c.SetOptions(options.NoInteractiveMode)
		if c.options&options.NoInteractiveMode == 0 {
			t.Error("the `SetOptions` method did not save the option given")
		}
	}

	// Ensure we can set multiple options
	{
		c := New()
		_ = c.SetOptions(1, 2, 4)
		if c.options&7 != 7 {
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
		_ = c.SetCommands(
			lime.Command{
				Keyword: "test",
			},
		)

		if len(c.commands) != 1 {
			t.Error("the `SetCommands` method did not save the command given")
		}

		if (c.commands)[0].Keyword != "test" {
			t.Error("the `SetCommands` method saved a command but its keyword value was lost")
		}
	}

	// Ensure we can set multiple options
	{
		c := New()
		_ = c.SetCommands(
			lime.Command{
				Keyword: "test1",
			},
			lime.Command{
				Keyword: "test2",
			},
		)

		if len(c.commands) != 2 {
			t.Error("the `SetCommands` method did not save the commands given")
		}

		if (c.commands)[0].Keyword != "test1" {
			t.Error("the `SetCommands` method saved the first command but its keyword value was lost")
		}

		if (c.commands)[1].Keyword != "test2" {
			t.Error("the `SetCommands` method saved the second command but its keyword value was lost")
		}
	}
}

func TestCLI_SetName(t *testing.T) {
	c := New()
	c.SetName("test")
	if c.name != "test" {
		t.Error("the `SetName` method did not save the name")
	}
}

func TestCLI_SetPrompt(t *testing.T) {
	c := New()
	c.SetPrompt("test")
	if c.prompt != "test" {
		t.Error("the `SetPrompt` method did not save the prompt")
	}
}

func TestCLI_SetExitWord(t *testing.T) {
	c := New()
	c.SetExitWord("test")
	if c.exitWord != "test" {
		t.Error("the `SetExitWord` method did not save the exit word")
	}
}
