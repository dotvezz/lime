package errors

import "errors"

// ErrNoMatch is returned when Lime was unable to find a `lime.Command` matching the args/input
var ErrNoMatch = errors.New("no matching command found")

// ErrNoFunc is returned when Lime found a matching `lime.Command`, but the command had no `lime.Func` to invoke
var ErrNoFunc = errors.New("no function for command")

// ErrInvalidOption is returned when the option given to `CLI.SetOptions` is not a power of 2 (because `lime.Options` is a bit mask)
var ErrInvalidOption = errors.New("an invalid option value was given")

// ErrNoInput is returned when Lime was unable to find args/input to use
var ErrNoInput = errors.New("no command given")
