package functions

import (
	"strings"

	"github.com/fatih/color"
)

type Printable interface {
	Print(func(format string, a ...any))
}

// Remove excess space, split input and make it lowercase
func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func Filter[T any](items []T, fn func(T, int) bool) []T {
	filteredItems := []T{}
	for i, value := range items {
		if fn(value, i) {
			filteredItems = append(filteredItems, value)
		}
	}
	return filteredItems
}

func PrintList[T Printable](items []T) {
	printColor1 := color.New(color.FgGreen).PrintfFunc()
	printColor2 := color.New(color.FgCyan).PrintfFunc()
	for i, item := range items {
		if i%2 == 0 {
			item.Print(printColor2)
		} else {
			item.Print(printColor1)
		}
	}
}

func PrintFormattedList[T any](items []T) {

	for i, item := range items {
		if i%2 == 0 {
			color.Cyan("% 3d. %s\n", i+1, item)
		} else {
			color.Green("% 3d. %s\n", i+1, item)
		}
	}
}

// slices.Delete doesn't seem to work as I need it to. It doesn't seem to update the slice length, leaving a bunch of zeroed out items at the end, which breaks my design
func DeleteFromSlice[T any](items []T, idx int) []T {
	return append(items[:idx], items[idx+1:]...)
}

func DeleteFromSliceFunc[T any](items []T, pred func(T) bool) []T {
	result := []T{}
	for _, item := range items {
		if !pred(item) {
			result = append(result, item)
		}
	}
	return result
}
