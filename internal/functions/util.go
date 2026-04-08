package functions

import "strings"

// Remove excess space, split input and make it lowercase
func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func Filter[T any](items []T, fn func(item T, i int) bool) []T {
	filteredItems := []T{}
	for i, value := range items {
		if fn(value, i) {
			filteredItems = append(filteredItems, value)
		}
	}
	return filteredItems
}
