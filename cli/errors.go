package cli

import "errors"

// errNoMatch is returned when Lime was unable to find a `lime.Command` matching the args/input
var errNoMatch = errors.New("no matching command found")

// errNoFunc is returned when Lime found a matching `lime.Command`, but the command had no `lime.Func` to invoke
var errNoFunc = errors.New("no function for command")

// errInvalidOption is returned when the option given to `CLI.SetOptions` is not a power of 2 (because `lime.Options` is a bit mask)
var errInvalidOption = errors.New("an invalid option value was given")

// errNoInput is returned when Lime was unable to find args/input to use
var errNoInput = errors.New("no command given")

// errNoHelp is returned when the `Help` property is needed but not present.
var errNoHelp = errors.New("no help provided for this command")

// errNoHelp is returned when the `Help` property is needed but not present.
var errNoUsage = errors.New("no usage provided for this command")