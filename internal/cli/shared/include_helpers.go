package shared

import "slices"

// HasInclude returns true when include is present in values.
func HasInclude(values []string, include string) bool {
	return slices.Contains(values, include)
}
