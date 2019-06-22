package errors

import "errors"

var NoMatch = errors.New("no matching command found")
var NoFunc = errors.New("no function for command")
var InvalidOption = errors.New("an invalid option value was given")
var NoInput = errors.New("no command given")
