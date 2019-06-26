# lime
lime is a small CLI library for Go.

[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/dotvezz/lime)](https://goreportcard.com/report/github.com/dotvezz/lime)
[![codecov](https://codecov.io/gh/dotvezz/lime/branch/master/graph/badge.svg)](https://codecov.io/gh/dotvezz/lime)
[![CircleCI](https://circleci.com/gh/dotvezz/lime/tree/master.svg?style=shield)](https://circleci.com/gh/dotvezz/lime/tree/master)

## The Building Blocks

When using lime, there are two basic types that power your CLI application. `lime.Command` and `lime.Func`

`lime.Command` defines the basic structure and routing paths for the application, and `lime.Func` defines 
the signature of a function that is invoked by a command. Below is a basic example.

 ```go
package main
 
import (
	"fmt"
	"github.com/dotvezz/lime"
	"github.com/dotvezz/lime/cli"
)

func main() {
	var command = lime.Command{
		// Keyword is the value that will be matched against the command line arguments
		Keyword: "greet",
		// Func is the function that will be run for this command.
		Func: func(args []string) error {
			if len(args) > 0 {
				fmt.Printf("Hello, %s!\n", args[0])
			} else {
				fmt.Println("Hello, world!")
			}
			return nil
		},
	}

	mycli := cli.New()
	_ = mycli.SetCommands(command)
	_ = mycli.Run()
}
 ```

The program above runs like this:

```
> myCli greet
Hello, world!
> myCli greet John
Hello, John!
```

The `args` value passed to the function is a slice of `os.Args`, excluding the args used to match the
`lime.Command`.  So running `myCli greet John` provides the function with the value `[]string{"John"}`.

## Features

If you would like to suggest a feature for lime to support, please 
[open an issue](https://github.com/dotvezz/lime/issues) or submit a pull request for it.

### Interactive Shell

By default, lime gives your CLI an interactive shell. In a future release, one goal is for the 
shell to run as an interpreter for custom scripts.

### Basic Command Handling

Of course, lime supports plain old commands.

```go
var command = lime.Command{
	Keyword:	 "bark",
	Func: func(_ []string) error {
		fmt.Println("woof")
		return nil
	},
}
```

Commands can also be nested.

```go
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
```

#### Command Help, Usage, Description

When building your CLI with lime, you can provide usage examples as well as help and descriptions.

## Goals

The lime project has a number of goals. Some goals are general and intended as guidelines to the 
project's design, and others are specific features that it is intended to support at some point 
in the future.

### Guidelines

- Make it simple to write a powerful CLI without needing to read a bunch of documentation.
- Take as little control away from a CLI developer as possible.

### Feature Wish List

- Ability for the interactive shell mode run as an interpreter for custom scripts.
- Support for dynamic prompts in the interactive shell mode
- Support for command-line flags
- Support for `bash` auto-completion

## Release Status and Interface Stability

Lime is currently under early development and has not yet had any release. For the time being,
expect the `lime.CLI` interface to undergo changes.
