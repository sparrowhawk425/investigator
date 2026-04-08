package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/samber/lo"

	"github.com/sparrowhawk425/investigators/internal/commands"
	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/gamelogic"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
	"github.com/sparrowhawk425/investigators/internal/times"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	gameState := gamelogic.GameState{
		Scanner:   scanner,
		DayNumber: 1,
		TimeOfDay: times.Morning,
	}
	fmt.Print("What is your name, Investigator? > ")
	scanner.Scan()
	name := scanner.Text()
	fmt.Printf("Welcome to International Investigators, %s!\n", name)
	// Select country
	countryNames := lo.Map(nameapi.Countries, func(country nameapi.Country, i int) string { return country.Name })
	idx := gamelogic.MenuSelect(scanner, "Select a Country to begin your investigation:", countryNames)
	country := nameapi.Countries[idx]
	fmt.Printf("Travelling to %s...\n", country.Name)

	gameState.BuildGame(country)
	fmt.Printf("We estimate %d Syndicate members are currently in the area.\n", len(gameState.Criminals))

	// REPL game loop
	commands := commands.GetCommandMap()
	var err error
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
