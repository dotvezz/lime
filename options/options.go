package options

import (
	"math"

	"github.com/dotvezz/lime"
)

const (
	// NoInteractiveMode disables the interactive interface
	NoInteractiveMode lime.Option = 1 << iota
	// PrintErrors enables output of errors to stdout
	PrintErrors
)

// IsValid returns true if the option passed is a power of 2, or returns false otherwise
func IsValid(option lime.Option) bool {
	return math.Ceil(math.Log2(float64(option))) != math.Floor(math.Log2(float64(option)))
}
