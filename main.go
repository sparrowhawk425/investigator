package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/samber/lo"

	"github.com/sparrowhawk425/investigators/commands"
	"github.com/sparrowhawk425/investigators/gamelogic"
	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/gameobjects/enemies"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
	"github.com/sparrowhawk425/investigators/times"
)

func main() {
	day := 1
	timeOfDay := times.Morning
	scanner := bufio.NewScanner(os.Stdin)
	countryNames := lo.Map(nameapi.Countries, func(country nameapi.Country, i int) string {
		return country.Name
	})
	idx := menuSelect(scanner, "Select a Country to begin your investigation:", countryNames)
	country := nameapi.Countries[idx]
	fmt.Printf("Travelling to %s...\n", country.Name)
	nameResults, err := nameapi.MakeHTTPGetRequest(country, 1)
	if err != nil {
		log.Fatalf("Error getting name from API: %v", err)
	}
	char := nameResults[0]
	target := enemies.CreateEnemy(char)
	fmt.Printf("Hunting for %s %s, known %s\n", target.Character.Name.First, target.Character.Name.Last, target.Profiles[0].Name)

	// Add locations and people to game
	gameState := gamelogic.GameState{}
	results, err := nameapi.MakeHTTPGetRequest(country, 10)
	if err != nil {
		log.Fatalf("Error getting locations from API: %v", err)
	}
	// Create people
	for _, c := range results {
		gameState.People = append(gameState.People, gameobjects.CreateRandomCharacter(c))
	}
	// Create locations
	apiLocations := lo.Map(results, func(character nameapi.Character, i int) nameapi.Location {
		return character.Location
	})
	locations := gameobjects.CreateRandomLocations(apiLocations)
	gameState.Places = locations

	// REPL game loop
	commands := commands.GetCommandMap()
	for {
		fmt.Printf("Day: %d Time: %s\n", day, times.GetTimeOfDayName(timeOfDay))
		fmt.Print("What do you wish to do? > ")
		scanner.Scan()
		cleanText := cleanInput(scanner.Text())
		cmd, exists := commands[cleanText[0]]
		if exists {
			if err := cmd.Callback(&gameState, cleanText[1:]); err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
		timeOfDay = times.TransitionTimeOfDay(timeOfDay)
		if timeOfDay == times.Morning {
			day++
		}
	}
}

// Remove excess space, split input and make it lowercase
func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

// Given a list of items, allow the player to make a numeric selection
func menuSelect(scanner *bufio.Scanner, msg string, items []string) int {
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
		if idx < 0 || idx >= len(nameapi.Countries) {
			idx = -1
			fmt.Println("Invalid Choice")
		}
	}
	return idx
}
