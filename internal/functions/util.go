package functions

import "strings"

// Remove excess space, split input and make it lowercase
func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
