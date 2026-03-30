package commands

import (
	"fmt"

	"github.com/samber/lo"
	"github.com/sparrowhawk425/investigators/gamelogic"
	"github.com/sparrowhawk425/investigators/internal/functions"
)

func GetDossierCommandMap() map[string]cliCommand {
	commandMap := map[string]cliCommand{
		"exit": {
			name:         "exit",
			description:  "Exit Dossier menu",
			advancesTime: false,
			Callback:     commandDossierExit,
		},
		"help": {
			name:         "help",
			description:  "Displays a help message",
			advancesTime: false,
			Callback:     commandDossierHelp,
		},
		"view": {
			name:         "view",
			description:  "Select a dossier to view in detail",
			advancesTime: false,
			Callback:     commandDossierView,
		},
		"create": {
			name:         "create",
			description:  "Create a new Dossier",
			advancesTime: false,
			Callback:     commandDossierCreate,
		},
		"update": {
			name:         "update",
			description:  "Update an existing Dossier",
			advancesTime: false,
			Callback:     commandDossierUpdate,
		},
	}
	return commandMap
}

func commandDossiers(gs *gamelogic.GameState, _ []string) (bool, error) {
	fmt.Print("Current Dossiers:")
	if len(gs.Player.Dossiers) == 0 {
		fmt.Println(" None")
	} else {
		fmt.Printf(" %d Active\n", len(gs.Player.Dossiers))
	}
	commands := GetDossierCommandMap()
	scanner := gs.Scanner

	for {
		// Get player input
		fmt.Print("Dossiers: What do you wish to do? > ")
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

func commandDossierExit(gs *gamelogic.GameState, _ []string) (bool, error) {
	return true, nil
}

func commandDossierHelp(gs *gamelogic.GameState, _ []string) (bool, error) {
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range GetDossierCommandMap() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return false, nil
}

func commandDossierView(gs *gamelogic.GameState, _ []string) (bool, error) {

	if len(gs.Player.Dossiers) == 0 {
		fmt.Println("No Dossiers to view")
		return false, nil
	}
	idx := gamelogic.MenuSelect(gs.Scanner, "Choose a Dossier:", lo.Map(gs.Player.Dossiers, func(d gamelogic.Dossier, _ int) string { return d.Name }))
	d := gs.Player.Dossiers[idx]
	d.Print()

	return false, nil
}

func commandDossierCreate(gs *gamelogic.GameState, _ []string) (bool, error) {
	gs.Player.CreateDossier(gs.Scanner)

	return false, nil
}

func commandDossierUpdate(gs *gamelogic.GameState, _ []string) (bool, error) {

	if len(gs.Player.Dossiers) == 0 {
		fmt.Println("No Dossiers to update")
		return false, nil
	}
	idx := gamelogic.MenuSelect(gs.Scanner, "Choose a Dossier:", lo.Map(gs.Player.Dossiers, func(d gamelogic.Dossier, _ int) string { return d.Name }))
	gs.Player.Dossiers[idx].Update(gs.Scanner)

	return false, nil
}
