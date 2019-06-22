# lime
lime is a small CLI library for Go.

## Features

If you would like to suggest a feature for lime to support, please [open an issue](https://github.com/dotvezz/lime/issues) or submit a pull request for it.

### Interactive Shell

By default, lime gives your CLI an interactive shell. In a future release, one goal is to add
support for running scripts for the custom shell

### Basic Command Handling

Of course, lime supports plain old commands.

```go
package main

import (
	"errors"
	"fmt"
	"github.com/dotvezz/lime"
)

var command = lime.Command{
    Keyword:     "bark",
    Func: func(_ []string) error {
        fmt.Println("woof")
        return nil
    },
}
```

Commands can also be nested.

```go
package main

import (
	"fmt"
	"github.com/dotvezz/lime"
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
```

#### Command Help, Usage, Description

When building your CLI with lime, you can provide usage examples as well as help and descriptions.

## Release Status and Interface Stability

Lime is currently under early development and has not yet had any release. For the time being,
expect the `lime.CLI` interface to undergo changes.