package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/sparrowhawk425/investigators/internal/commands"
	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/gamelogic"
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
	gameState.Player.Name = name
	fmt.Printf("Welcome to International Investigators, %s!\n", name)
	fmt.Println("We believe members of the Syndicate have spread all across the world and we need your help to track them down.")

	// Build the game
	gameState.BuildGame()
	fmt.Printf("We know %d Syndicate members are currently in the area.\n", len(gameState.Criminals))

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
			// TODO: Currently character locations are off-by-one time (only added on update, so won't appear during visit, but will appear one time after)
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
