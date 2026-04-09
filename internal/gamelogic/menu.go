package gamelogic

import (
	"bufio"
	"fmt"
	"strconv"
)

type ToString interface {
	String() string
}

type FilterType interface {
	IsType() bool
	String() string
}

type FilterFunc[T any] func(T, int) bool

// TODO: Pass map of filtertype to filter function?
// func GetLocationFilterFunc[T any](scanner *bufio.Scanner, filterType FilterType) FilterFunc[T] {
// 	switch filterType.(type) {
// 	case gameobjects.LocationType:
// 		fmt.Println("location type")
// 		locTypes := []gameobjects.LocationType{}
// 		availableTypes := gameobjects.LocationTypes
// 		done := false
// 		for !done {
// 			idx := MenuSelect(scanner, "Select Location Type:", lo.Map(availableTypes, func(lt gameobjects.LocationType, _ int) string { return lt.String() }))
// 			locTypes = append(locTypes, availableTypes[idx])
// 			availableTypes = slices.Delete(availableTypes, idx, idx)
// 			fmt.Print("Add another? > ")
// 			scanner.Scan()
// 			done = scanner.Text() == "y"
// 		}
// 		return gameobjects.FilterLocationsByType(locTypes)
// 	case gameobjects.LootType:
// 		fmt.Println("LootType")
// 		return gameobjects.FilterLocationsByLootType(nil)
// 	case gameobjects.Quality:
// 		fmt.Println("Quality")

// 		return gameobjects.FilterLocationsByQuality(nil)
// 	}
// 	return nil
// }

func CreateFilterableMenu[T ToString](scanner *bufio.Scanner, prompt string, items []T, filterTypes []FilterType) int {

	// TODO: Figure out how to set up generic filter function selectors
	// filterIdx := MenuSelect(scanner, "Choose Filter:", lo.Map(filterTypes, func(ft FilterType, _ int) string { return ft.String() }))
	// filterFunc := GetLocationFilterFunc[T](scanner, filterTypes[filterIdx])
	// availableItems := filter(items, filterFunc)
	idx := -1
	for idx < 0 {
		menuItems := make([]string, len(items))
		for i, item := range items {
			menuItems[i] = item.String()
		}
		fmt.Println(prompt)
		for i, mItem := range menuItems {
			fmt.Printf("%d. %s\n", i+1, mItem)
		}
		fmt.Print("Select an item > ")
		scanner.Scan()
		var err error
		idx, err = strconv.Atoi(scanner.Text())
		if err != nil {
			idx = -1
			fmt.Println("Invalid choice")
			continue
		}
		idx--
		if idx < 0 || idx >= len(items) {
			idx = -1
			fmt.Println("Invalid Choice")
		}
	}
	return idx
}

// Given a list of items, allow the player to make a numeric selection
func MenuSelect(scanner *bufio.Scanner, msg string, items []string) int {
	idx := -1
	for idx < 0 {
		fmt.Println(msg)
		for i, item := range items {
			fmt.Printf("%d. %s\n", i+1, item)
		}
		fmt.Print("Which number? > ")
		scanner.Scan()
		var err error
		idx, err = strconv.Atoi(scanner.Text())
		if err != nil {
			idx = -1
			fmt.Println("Invalid Choice")
			continue
		}
		idx--
		if idx < 0 || idx >= len(items) {
			idx = -1
			fmt.Println("Invalid Choice")
		}
	}
	return idx
}
