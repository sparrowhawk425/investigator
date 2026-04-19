package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
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
	boldGreen := color.New(color.FgGreen, color.Bold)
	boldGreen.Printf("Welcome to International Investigators, %s!\n", name)
	fmt.Println("We believe members of the Syndicate have spread all across the world and we need your help to track them down.")

	// Build the game
	gameState.BuildGame()
	fmt.Printf("We know %d Syndicate members are currently in the area.\n", len(gameState.Criminals))

	// REPL game loop
	commands := commands.GetCommandMap()
	var err error
	for {
		// TODO: Expand to move to another region, chasing bosses
		if len(gameState.Criminals) == 0 {
			if len(gameState.Escaped) == 0 {
				color.Green("Congratulations! You have caught all the Synidcate members in the area!")
				color.Green("You Win!")
			} else if len(gameState.Caught) == 0 {
				color.Red("All the Syndicate members have escaped!")
				color.Red("You lose!")
			} else {
				color.Yellow("You caught %d Syndicate members and %d escaped!", len(gameState.Caught), len(gameState.Escaped))
				color.Yellow("Try again!")
			}
			return
		}
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
