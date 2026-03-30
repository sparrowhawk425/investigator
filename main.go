package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/samber/lo"

	"github.com/sparrowhawk425/investigators/commands"
	"github.com/sparrowhawk425/investigators/gamelogic"
	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/gameobjects/enemies"
	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
	"github.com/sparrowhawk425/investigators/times"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	gameState := gamelogic.GameState{
		Scanner:   scanner,
		Day:       1,
		TimeOfDay: times.Morning,
	}
	fmt.Print("What is your name, investigator? > ")
	scanner.Scan()
	name := scanner.Text()
	fmt.Printf("Welcome to International Investigators, %s!\n", name)
	// Select country
	countryNames := lo.Map(nameapi.Countries, func(country nameapi.Country, i int) string { return country.Name })
	idx := gamelogic.MenuSelect(scanner, "Select a Country to begin your investigation:", countryNames)
	country := nameapi.Countries[idx]
	fmt.Printf("Travelling to %s...\n", country.Name)
	nameResults, err := nameapi.MakeHTTPGetRequest(country, 1)
	if err != nil {
		log.Fatalf("Error getting name from API: %v", err)
	}
	char := nameResults[0]
	target := enemies.CreateEnemy(char)
	fmt.Printf("Hunting for %s, known %s\n", target.Character.GetName(), target.Profile.Name)

	// Add locations and people to game
	gameState.Criminals = append(gameState.Criminals, target)
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
		gameState.PrintDay()

		// Get player input
		fmt.Print("What do you wish to do? > ")
		scanner.Scan()
		cleanText := functions.CleanInput(scanner.Text())

		// Run command
		update := false
		cmd, exists := commands[cleanText[0]]
		if exists {
			update, err = cmd.Callback(&gameState, cleanText[1:])
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
		if update {
			gameState.Update()
		}
	}
}
