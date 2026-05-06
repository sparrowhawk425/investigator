package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	"github.com/fatih/color"
	"github.com/sparrowhawk425/investigators/gameloop"
	"github.com/sparrowhawk425/investigators/internal/characters"
	"github.com/sparrowhawk425/investigators/internal/config"
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
	cfg := config.Config{
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
		// TODO: Set progression based on travelling to boss' countries
		if slices.ContainsFunc(bosses, func(c characters.Character) bool { return c.Traits.Nationality == country.GetCode() }) {
			color.Green("You have successfully tracked a Syndicate boss to their home base.")
		}
		level := gameloop.PlayLevel(country, bosses, &cfg)
		bosses = level.Bosses
		levels = append(levels, level)
		// Advance to next level and remove the latest country from the list
		cfg.Level++
		cfg.Countries = slices.DeleteFunc(cfg.Countries, func(c config.Country) bool { return c.Equals(level.Country) })
	}
}
