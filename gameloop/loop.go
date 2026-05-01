package gameloop

import (
	"bufio"
	"fmt"

	"github.com/fatih/color"
	"github.com/sparrowhawk425/investigators/internal/commands"
	"github.com/sparrowhawk425/investigators/internal/config"
	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/gamelogic"
	"github.com/sparrowhawk425/investigators/internal/times"
)

func PlayGame(scanner *bufio.Scanner, cfg *config.Config) {

	gs := gamelogic.GameState{
		Scanner:   scanner,
		DayNumber: 1,
		TimeOfDay: times.Morning,
	}
	gs.Player.Name = cfg.PlayerName
	// Build the game
	gs.BuildGame(cfg)
	fmt.Printf("We know %d Syndicate members are currently in the area.\n", len(gs.Criminals))

	// REPL game loop
	commands := commands.GetCommandMap()
	var err error
	for len(gs.Criminals) > 0 {
		gs.PrintDay()

		// Get player input
		fmt.Print("What do you wish to do? > ")
		scanner.Scan()
		cleanText := functions.CleanInput(scanner.Text())

		// Run command
		update := false
		cmd, exists := commands[cleanText[0]]
		if exists {
			update, err = cmd.Callback(&gs, cleanText[1:])
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		} else {
			fmt.Println("Unknown command")
		}
		if update {
			gs.Update()
		}
	}

	// TODO: Expand to move to another region, chasing bosses. Each member you catch has a few postcards with locations from around the world. One points to the correct next area. Catch enough of the members, you can narrow down the location
	if len(gs.Escaped) == 0 {
		color.Green("Congratulations! You have caught all the Synidcate members in the area!")
		color.Green("You Win!")
	} else if len(gs.Caught) == 0 {
		color.Red("All the Syndicate members have escaped!")
		color.Red("You lose!")
	} else {
		color.Yellow("You caught %d Syndicate members and %d escaped!", len(gs.Caught), len(gs.Escaped))
		color.Yellow("Try again!")
	}
}
