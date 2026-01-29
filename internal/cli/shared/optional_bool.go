package shared

import (
	"fmt"
	"strconv"
)

// OptionalBool tracks whether a boolean flag was explicitly set.
type OptionalBool struct {
	set   bool
	value bool
}

func (b *OptionalBool) Set(value string) error {
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("must be true or false")
	}
	b.value = parsed
	b.set = true
	return nil
}

func (b *OptionalBool) String() string {
	if !b.set {
		return ""
	}
	return strconv.FormatBool(b.value)
}

func (b *OptionalBool) IsBoolFlag() bool {
	return true
}

func (b OptionalBool) IsSet() bool {
	return b.set
}

func (b OptionalBool) Value() bool {
	return b.value
}
