package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/sparrowhawk425/investigators/gameloop"
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
	config := config.Config{
		Difficulty: config.Easy,
		Level:      1,
		PlayerName: name,
		Countries:  config.Countries,
		PostCards:  config.LoadPostCards(),
	}

	gameloop.PlayGame(scanner, &config)

}
