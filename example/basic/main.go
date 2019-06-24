package main

import (
	"errors"
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/cli"
	"os"
)

var command = lime.Command{
	Keyword: "tell",
	Commands: []lime.Command{
		{
			Keyword: "lie",
			Func: func(_ []string) error {
				fmt.Println("The author of this cli likes to eat oranges.")
				return nil
			},
		},
		{
			Keyword: "truth",
			Func: func(_ []string) error {
				fmt.Println("The author of this cli likes to eat apples.")
				return nil
			},
		},
	},
}

var commands = []lime.Command{
	{
		Keyword: "tell",
		Commands: []lime.Command{
			{
				Keyword: "lie",
				Func: func(_ []string) error {
					fmt.Println("The author of this cli likes to eat oranges.")
					return nil
				},
				Description: "Makes a preset statement which is factually untrue.",
			},
			{
				Keyword: "truth",
				Func: func(_ []string) error {
					fmt.Println("The author of this cli likes to eat apples.")
					return nil
				},
				Description: "Makes a preset statement which is factually true.",
			},
		},
	},
	{
		Keyword:     "repeat",
		Description: "Repeats all the words after the command.",
		Usage: []lime.Usage{
			{
				Example:     "mycli repeat the quick brown fox",
				Explanation: `outputs "[the quick brown fox]"`,
			},
			{
				Example:     "mycli repeat",
				Explanation: `returns an error: "there are no words to repeat"`,
			},
		},
		Func: func(args []string) error {
			if len(args) == 0 {
				return errors.New("there were no words to repeat")
			}
			fmt.Println(args)
			return nil
		},
	},
}

func main() {
	mycli := cli.New()
	_ = mycli.SetCommands(commands...)
	err := mycli.Run()
	if err != nil {
		os.Exit(1)
	}
}
