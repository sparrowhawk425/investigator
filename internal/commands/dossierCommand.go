package commands

import (
	"bufio"
	"fmt"

	"github.com/samber/lo"
	"github.com/sparrowhawk425/investigators/internal/characters"
	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/gamelogic"
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
	idx := gamelogic.MenuSelect(gs.Scanner, "Choose a Dossier:", lo.Map(gs.Player.Dossiers, func(d characters.Dossier, _ int) string { return d.Name }))
	d := gs.Player.Dossiers[idx]
	d.Print()

	return false, nil
}

func commandDossierCreate(gs *gamelogic.GameState, _ []string) (bool, error) {
	dossier := gs.Player.CreateDossier(gs.Scanner)
	updateDossier(gs.Scanner, dossier)

	return false, nil
}

func commandDossierUpdate(gs *gamelogic.GameState, _ []string) (bool, error) {

	if len(gs.Player.Dossiers) == 0 {
		fmt.Println("No Dossiers to update")
		return false, nil
	}
	idx := gamelogic.MenuSelect(gs.Scanner, "Choose a Dossier:", lo.Map(gs.Player.Dossiers, func(d characters.Dossier, _ int) string { return d.Name }))

	updateDossier(gs.Scanner, &gs.Player.Dossiers[idx])

	return false, nil
}

type dossierMenu string

const (
	dName      dossierMenu = "Dossier Name"
	dCharacter dossierMenu = "Character"
	dNotes     dossierMenu = "Notes"
	dDone      dossierMenu = "Done"
)

var dossierMenuItems = []dossierMenu{
	dName, dCharacter, dNotes, dDone,
}

func updateDossier(scanner *bufio.Scanner, dossier *characters.Dossier) {
	isDone := false
	for !isDone {
		idx := gamelogic.MenuSelect(scanner, "Choose Field:", lo.Map(dossierMenuItems, func(m dossierMenu, _ int) string { return string(m) }))
		option := dossierMenuItems[idx]
		switch option {
		case dName:
			fmt.Print("New Dossier Name: ")
			scanner.Scan()
			dossier.Name = scanner.Text()
		case dCharacter:
			updateCharacter(scanner, dossier)
		case dNotes:
			fmt.Print("Note: ")
			scanner.Scan()
			dossier.AddNote(scanner.Text())
		case dDone:
			isDone = true
		}
	}
}

type characterMenu string

const (
	cName        characterMenu = "Name"
	cGender      characterMenu = "Gender"
	cNationality characterMenu = "Nationality"
	cHeight      characterMenu = "Height"
	cWeight      characterMenu = "Weight"
	cEyeColor    characterMenu = "Eye Color"
	cHairColor   characterMenu = "Hair Color"
	cHairLength  characterMenu = "Hair Length"
	cShoeSize    characterMenu = "Shoe Size"
	cDone        characterMenu = "Done"
)

var characterMenuItems = []characterMenu{
	cName, cGender, cNationality, cHeight, cWeight, cEyeColor, cHairColor, cHairLength, cShoeSize, cDone,
}

func updateCharacter(scanner *bufio.Scanner, dossier *characters.Dossier) {

	isDone := false
	for !isDone {
		idx := gamelogic.MenuSelect(scanner, "Choose Field:", lo.Map(characterMenuItems, func(m characterMenu, _ int) string { return string(m) }))
		option := characterMenuItems[idx]
		switch option {
		case cName:
			fmt.Print("First Name: ")
			scanner.Scan()
			dossier.Target.SetFirstName(scanner.Text())
			fmt.Print("Last Name: ")
			scanner.Scan()
			dossier.Target.SetLastName(scanner.Text())
		case cGender:
			idx = gamelogic.MenuSelect(scanner, "Gender:", lo.Map(characters.Genders, func(g characters.Gender, _ int) string { return g.String() }))
			dossier.Target.Traits.Gender = characters.Genders[idx]
		case cNationality:
			idx = gamelogic.MenuSelect(scanner, "Nationality:", lo.Map(characters.Nationalities, func(n characters.Nationality, _ int) string { return n.String() }))
			dossier.Target.Traits.Nationality = characters.Nationalities[idx]
		case cHeight:
			idx = gamelogic.MenuSelect(scanner, "Height:", lo.Map(characters.Heights, func(h characters.Height, _ int) string { return h.String() }))
			dossier.Target.Traits.Height = characters.Heights[idx]
		case cWeight:
			idx = gamelogic.MenuSelect(scanner, "Weight:", lo.Map(characters.Weights, func(w characters.Weight, _ int) string { return w.String() }))
			dossier.Target.Traits.Weight = characters.Weights[idx]
		case cEyeColor:
			idx = gamelogic.MenuSelect(scanner, "Eye Color:", lo.Map(characters.EyeColors, func(ec characters.EyeColor, _ int) string { return ec.String() }))
			dossier.Target.Traits.EyeColor = characters.EyeColors[idx]
		case cHairColor:
			idx = gamelogic.MenuSelect(scanner, "Hair Color:", lo.Map(characters.HairColors, func(hc characters.HairColor, _ int) string { return hc.String() }))
			dossier.Target.Traits.HairColor = characters.HairColors[idx]
		case cHairLength:
			idx = gamelogic.MenuSelect(scanner, "Hair Length:", lo.Map(characters.HairLengths, func(hl characters.HairLength, _ int) string { return hl.String() }))
			dossier.Target.Traits.HairLength = characters.HairLengths[idx]
		case cShoeSize:
			idx = gamelogic.MenuSelect(scanner, "Shoe Size:", lo.Map(characters.ShoeSizes, func(ss characters.ShoeSize, _ int) string { return ss.String() }))
			dossier.Target.Traits.ShoeSize = characters.ShoeSizes[idx]
		case cDone:
			isDone = true
		}
	}
}
