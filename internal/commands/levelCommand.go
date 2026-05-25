package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/gamelogic"
)

func GetLevelCommandMap() map[string]cliCommand {
	commandMap := map[string]cliCommand{
		"done": {
			name:        "done",
			description: "Exit Level menu",
			Callback:    commandLevelExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			Callback:    commandLevelHelp,
		},
		"postcards": {
			name:        "postcards",
			description: "Review the collected postcards from the level",
			Callback:    commandPostCards,
		},
		"stats": {
			name:        "stats",
			description: "Displays a list of stats for the level",
			Callback:    commandLevelStats,
		},
	}
	return commandMap
}

func ReviewLevel(gs *gamelogic.GameState) {

	if len(gs.Escaped) == 0 {
		color.Green("Congratulations! You have caught all the Synidcate members in the area!")
	} else if len(gs.Caught) == 0 {
		color.Red("All the Syndicate members have escaped!")
	} else {
		color.Yellow("You caught %d Syndicate members and %d escaped!", len(gs.Caught), len(gs.Escaped))
	}
	commands := GetLevelCommandMap()

	green := color.New(color.FgGreen).SprintFunc()
	scanner := gs.Scanner
	for {
		// Get player input
		fmt.Printf("%s What do you wish to do? > ", green("Review:"))
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
}

func commandLevelExit(gs *gamelogic.GameState, _ []string) (bool, error) {
	return true, nil
}

func commandLevelHelp(gs *gamelogic.GameState, _ []string) (bool, error) {
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range GetLevelCommandMap() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return false, nil
}

func commandLevelStats(gs *gamelogic.GameState, _ []string) (bool, error) {

	fmt.Printf("%s Stats:\n", gs.Country.GetName())
	fmt.Printf("  Days: %d\n", gs.DayNumber)
	fmt.Printf("  Escaped: %d\n", len(gs.Escaped))
	fmt.Printf("  Caught: %d\n\n", len(gs.Caught))
	fmt.Println(" Caught:")
	functions.PrintFormattedList(gs.Caught)
	return false, nil
}
