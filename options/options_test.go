package options

import "testing"

func Test_OptionsIsValid(t *testing.T) {
	b := IsValid(3)
	if !b {
		t.Error("IsValid should not validate anything that is not a power of 2")
	}
}
