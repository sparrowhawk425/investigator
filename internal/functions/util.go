package functions

import (
	"bufio"
	"fmt"
	"strconv"
)

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
