package commands

import (
	"fmt"
	"os"

	"github.com/sparrowhawk425/investigators/internal/gamelogic"
	"github.com/sparrowhawk425/investigators/internal/gameobjects"
)

// TODO: Use argparse package? https://pkg.go.dev/github.com/akamensky/argparse

// Callback performs command action. Returns true if time should be updated
type cliCommand struct {
	name         string
	description  string
	advancesTime bool
	Callback     func(*gamelogic.GameState, []string) (bool, error)
}

func GetCommandMap() map[string]cliCommand {
	commandMap := map[string]cliCommand{
		"exit": {
			name:         "exit",
			description:  "End the investigation",
			advancesTime: false,
			Callback:     commandExit,
		},
		"help": {
			name:         "help",
			description:  "Displays a help message",
			advancesTime: false,
			Callback:     commandHelp,
		},
		"places": {
			name:         "places",
			description:  "List the available locations. Accepts a list of optional arguments to filter by one or more location types",
			advancesTime: false,
			Callback:     commandLocations,
		},
		"people": {
			name:         "people",
			description:  "List the people",
			advancesTime: false,
			Callback:     commandPeople,
		},
		"crimes": {
			name:         "crimes",
			description:  "List the crimes that have occurred",
			advancesTime: false,
			Callback:     commandCrimes,
		},
		"dossiers": {
			name:         "dossiers",
			description:  "View your dossiers",
			advancesTime: false,
			Callback:     commandDossiers,
		},
		"next": {
			name:         "next",
			description:  "Move to the next time period",
			advancesTime: true,
			Callback:     commandNext,
		},
		"visit": {
			name:         "visit",
			description:  "Choose a location to visit from a menu",
			advancesTime: true,
			Callback:     commandVisitLocation,
		},
		"enemies": {
			name:         "enemies",
			description:  "View info about enemies",
			advancesTime: false,
			Callback:     commandDebugEnemies,
		},
	}
	return commandMap
}

func commandExit(gs *gamelogic.GameState, _ []string) (bool, error) {
	fmt.Println("Closing the Dossier... Goodbye!")
	os.Exit(0)
	return false, nil
}

func commandHelp(gs *gamelogic.GameState, _ []string) (bool, error) {
	fmt.Println("Welcome to the International Investigators!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range GetCommandMap() {
		advTime := ""
		if cmd.advancesTime {
			advTime = " (Advances time)"
		}
		fmt.Printf("%s%s: %s\n", cmd.name, advTime, cmd.description)

	}
	return false, nil
}

// TODO: Menu with filtering?
func commandLocations(gs *gamelogic.GameState, params []string) (bool, error) {
	locations := gs.Places
	if len(params) > 0 {
		if params[0] != "filter" {
			return false, fmt.Errorf("Unexpected argument '%s'", params[0])
		}
		//filters := getFilters()
		//locations = gs.GetLocations(filters)

		locTypes := []gameobjects.LocationType{}
		for _, param := range params {
			locType, err := gameobjects.GetLocationType(param)
			if err != nil {
				return false, err
			}
			locTypes = append(locTypes, locType)
		}
		locations = gs.GetLocationsByType(locTypes)
	}
	for _, loc := range locations {
		loc.Print()
	}
	return false, nil
}

func commandPeople(gs *gamelogic.GameState, _ []string) (bool, error) {
	for _, person := range gs.People {
		person.Print()
	}
	return false, nil
}

func commandCrimes(gs *gamelogic.GameState, _ []string) (bool, error) {
	fmt.Println("Reported Crimes:")
	for _, crime := range gs.Crimes {
		crime.Print()
	}
	return false, nil
}

func commandVisitLocation(gs *gamelogic.GameState, _ []string) (bool, error) {
	locations := []string{}
	for _, location := range gs.Places {
		locations = append(locations, location.GetAddress())
	}

	//idx := functions.MenuSelect(gs.Scanner, "Choose a location:", locations)
	filterTypes := []gamelogic.FilterType{gameobjects.Residence, gameobjects.Cheap, gameobjects.Jewelry}

	idx := gamelogic.CreateFilterableMenu(gs.Scanner, "Choose a location:", gs.Places, filterTypes)
	loc := gs.Places[idx]
	gs.Player.CurrentLocation = loc
	fmt.Println(loc.GetAddress())

	if len(loc.Visitors) > 0 {
		fmt.Println("People:")
		for _, person := range loc.Visitors {
			fmt.Printf(" - %s\n", person.GetName())
		}
	}
	if len(loc.GetAvailableLoot()) > 0 {
		fmt.Println("Loot:")
		for _, loot := range loc.GetAvailableLoot() {
			amt := loc.GetLootAmount(loot)
			fmt.Printf(" - %s: %d\n", loot, amt)
		}
	}
	if len(loc.GetClues()) > 0 {
		fmt.Println("Clues:")
		for _, clue := range loc.GetClues() {
			fmt.Printf(" - %s\n", clue)
		}
	}

	return true, nil
}

func commandDebugEnemies(gs *gamelogic.GameState, _ []string) (bool, error) {
	for _, c := range gs.Criminals {
		c.Print()
	}
	return false, nil
}

func commandNext(gs *gamelogic.GameState, _ []string) (bool, error) {
	return true, nil
}
