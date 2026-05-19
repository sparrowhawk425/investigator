package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/fatih/color"
	"github.com/sparrowhawk425/investigators/gameloop"
	"github.com/sparrowhawk425/investigators/internal/characters"
	"github.com/sparrowhawk425/investigators/internal/config"
	"github.com/sparrowhawk425/investigators/internal/functions"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("What is your name, Investigator? > ")
	scanner.Scan()
	name := scanner.Text()
	boldGreen := color.New(color.FgGreen, color.Bold)
	boldGreen.Printf("Welcome to International Investigators, %s!\n", name)
	fmt.Println("We believe members of the Syndicate have spread all across the world and we need your help to track them down.")

	// TODO: incorporate loop to allow player to pursue next ring of targets in a new location
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.Config{
		Logger:     logger,
		Difficulty: config.Easy,
		Level:      1,
		PlayerName: name,
		Scanner:    scanner,
		Countries:  config.Countries,
		PostCards:  config.LoadPostCards(),
	}
	levels := []gameloop.Level{}
	bosses := []characters.Character{}
	for cfg.Level <= cfg.GetFinalLevel() {
		country := gameloop.SelectLevel(&cfg)
		if len(bosses) > 0 {
			if slices.ContainsFunc(bosses, func(boss characters.Character) bool { return boss.Traits.Nationality == country.Code }) {
				color.Green("You have successfully tracked a Syndicate boss to their headquarters.")
			} else {
				color.Yellow("There are no known Syndicate bosses in this region.")
			}
		}
		level := gameloop.PlayLevel(country, bosses, &cfg)
		bosses = level.Bosses
		levels = append(levels, level)
		// Advance to next level and remove the latest country from the list
		cfg.Level++
		cfg.Countries = functions.DeleteFromSliceFunc(cfg.Countries, func(c config.Country) bool { return c.Equals(level.Country) })
	}
	color.Green("Your investigation has come to an end.\n")
	totalCaught := 0
	totalEscaped := 0
	fmt.Println("Summary:")
	for _, level := range levels {
		totalCaught += len(level.Caught)
		totalEscaped += len(level.Escaped)
		fmt.Printf("While in %s:\n", level.Country.Name)
		fmt.Printf(" - Criminals Caught: %d\n - Criminals Escaped: %d\n", len(level.Caught), len(level.Escaped))
	}
	fmt.Printf("In total, you caught %d Syndicate members and %d escaped\n", totalCaught, totalEscaped)
	if totalEscaped == 0 {
		color.Green("Congratulations! You caught every Syndicate member. Well done!")
	} else if totalCaught == 0 {
		color.Red("Every Syndicate member escaped! Please do better next time.")
	} else {
		color.Yellow("You managed to catch some of the Syndicate members, but others escaped. We'll get them next time.")
	}
	fmt.Println("Thanks for playing!")
}
