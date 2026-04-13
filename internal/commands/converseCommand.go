package commands

import (
	"fmt"

	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/gamelogic"
)

func GetConverseCommandMap() map[string]cliCommand {
	commandMap := map[string]cliCommand{
		"exit": {
			name:         "exit",
			description:  "Exit Conversation menu",
			advancesTime: false,
			Callback:     commandConverseExit,
		},
		"help": {
			name:         "help",
			description:  "Displays a help message",
			advancesTime: false,
			Callback:     commandConverseHelp,
		},
	}
	return commandMap
}

func commandConverse(gs *gamelogic.GameState, _ []string) (bool, error) {

	commands := GetConverseCommandMap()
	scanner := gs.Scanner
	for {
		// Get player input
		fmt.Print("What do you want to ask? > ")
		scanner.Scan()
		cleanText := functions.CleanInput(scanner.Text())
		cmd, exists := commands[cleanText[0]]
		if exists {
			exit, err := cmd.Callback(gs, cleanText[1:])
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			if exit {
				break
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
	return false, nil
}

func commandConverseExit(gs *gamelogic.GameState, _ []string) (bool, error) {
	return true, nil
}

func commandConverseHelp(gs *gamelogic.GameState, _ []string) (bool, error) {
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range GetConverseCommandMap() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return false, nil
}
