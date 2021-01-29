package options

import (
	"github.com/dotvezz/lime"
	"math"
)

const (
	// NoShell disables the shell interface
	NoShell     lime.Option = 1 << iota
	// PrintErrors enables output of errors to stdout
	PrintErrors
)

// IsValid returns true if the option passed is a power of 2, or returns false otherwise
func IsValid(option lime.Option) bool {
	return math.Ceil(math.Log2(float64(option))) != math.Floor(math.Log2(float64(option)))
}
