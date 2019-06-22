package options

import "math"

// A bit mask value
type Option int64

const (
	// Disables the shell interface
	NoShell Option = 1 << iota
)

// IsValid returns true if the option passed is a power of 2, or returns false otherwise
func IsValid(option Option) bool {
	return math.Ceil(math.Log2(float64(option))) != math.Floor(math.Log2(float64(option)))
}
