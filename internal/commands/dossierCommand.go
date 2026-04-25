package commands

import (
	"bufio"
	"fmt"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/samber/lo"
	"github.com/sparrowhawk425/investigators/internal/characters"
	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/gamelogic"
)

func GetDossierCommandMap() map[string]cliCommand {
	commandMap := map[string]cliCommand{
		"close": {
			name:         "close",
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
		"create": {
			name:         "create",
			description:  "Create a new Dossier",
			advancesTime: false,
			Callback:     commandDossierCreate,
		},
		"view": {
			name:         "view",
			description:  "Select a dossier to view in detail. Optionally include the name of the dossier to match as an argument",
			advancesTime: false,
			Callback:     commandDossierView,
		},
		"update": {
			name:         "update",
			description:  "Update an existing Dossier. Optionally include the name of the dossier to match as an argument",
			advancesTime: false,
			Callback:     commandDossierUpdate,
		},
		"delete": {
			name:         "delete",
			description:  "Delete an existing Dossier. Optionally include the name of the dossier to match as an argument",
			advancesTime: false,
			Callback:     commandDossierDelete,
		},
		"match": {
			name:         "match",
			description:  "Get a list of matches for this Dossier. Optionally include the name of the dossier to match as an argument",
			advancesTime: false,
			Callback:     commandDossierMatch,
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
		for _, d := range gs.Player.Dossiers {
			fmt.Printf(" - %s\n", d.Name)
		}
	}
	commands := GetDossierCommandMap()
	scanner := gs.Scanner

	green := color.New(color.FgGreen).SprintFunc()
	for {
		// Get player input
		fmt.Printf("%s What do you wish to do? > ", green("Dossiers:"))
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

func commandDossierCreate(gs *gamelogic.GameState, _ []string) (bool, error) {
	dossier := gs.Player.CreateDossier(gs.Scanner)
	updateDossier(gs.Scanner, dossier)

	return false, nil
}

func commandDossierView(gs *gamelogic.GameState, params []string) (bool, error) {

	if len(gs.Player.Dossiers) == 0 {
		fmt.Println("No Dossiers to view")
		return false, nil
	}
	idx := getDossierIndex(gs, params)
	d := gs.Player.Dossiers[idx]
	d.Print()

	return false, nil
}

func commandDossierUpdate(gs *gamelogic.GameState, params []string) (bool, error) {

	if len(gs.Player.Dossiers) == 0 {
		fmt.Println("No Dossiers to update")
		return false, nil
	}
	idx := getDossierIndex(gs, params)
	if idx == -1 {
		fmt.Println("No dossier with that name exists")
		return false, nil
	}
	updateDossier(gs.Scanner, &gs.Player.Dossiers[idx])

	return false, nil
}

func commandDossierDelete(gs *gamelogic.GameState, params []string) (bool, error) {
	if len(gs.Player.Dossiers) == 0 {
		fmt.Println("No Dossiers to delete")
		return false, nil
	}
	idx := getDossierIndex(gs, params)
	if idx == -1 {
		fmt.Println("No Dossier with that name exists")
		return false, nil
	}
	gs.Player.Dossiers = slices.Delete(gs.Player.Dossiers, idx, idx)
	return false, nil
}

func commandDossierMatch(gs *gamelogic.GameState, params []string) (bool, error) {

	idx := getDossierIndex(gs, params)
	dossier := gs.Player.Dossiers[idx]
	matches := []characters.Character{}
	traits := dossier.Target.Traits
	matchTraits := "Matching on: "
	if traits.Gender != characters.UnknownGender {
		matchTraits += fmt.Sprintf("Gender: %s, ", traits.Gender)
	}
	if traits.Nationality != characters.UnknownNationality {
		matchTraits += fmt.Sprintf("From: %s, ", traits.Nationality)
	}
	if traits.EyeColor != characters.UnknownEyes {
		matchTraits += fmt.Sprintf("Eyes: %s, ", traits.EyeColor)
	}
	if traits.HairColor != characters.UnknownHairColor {
		matchTraits += fmt.Sprintf("Hair Color: %s, ", traits.HairColor)
	}
	if traits.HairLength != characters.UnknownHairLength {
		matchTraits += fmt.Sprintf("Hair Length: %s, ", traits.HairLength)
	}
	if traits.Height != characters.UnknownHeight {
		matchTraits += fmt.Sprintf("Height: %s, ", traits.Height)
	}
	if traits.Weight != characters.UnknownWeight {
		matchTraits += fmt.Sprintf("Weight: %s, ", traits.Weight)
	}
	if traits.ShoeSize != characters.UnknownShoe {
		matchTraits += fmt.Sprintf("Shoe Size: %s", traits.ShoeSize)
	}
	fmt.Println(matchTraits)
	for _, c := range gs.People {
		if traits.Gender != characters.UnknownGender {
			if c.Traits.Gender != traits.Gender {
				continue
			}
		}
		if traits.Nationality != characters.UnknownNationality {
			if c.Traits.Nationality != traits.Nationality {
				continue
			}
		}
		if traits.EyeColor != characters.UnknownEyes {
			if c.Traits.EyeColor != traits.EyeColor {
				continue
			}
		}
		if traits.HairColor != characters.UnknownHairColor {
			if c.Traits.HairColor != traits.HairColor {
				continue
			}
		}
		if traits.HairLength != characters.UnknownHairLength {
			if c.Traits.HairLength != traits.HairLength {
				continue
			}
		}
		if traits.Height != characters.UnknownHeight {
			if c.Traits.Height != traits.Height {
				continue
			}
		}
		if traits.Weight != characters.UnknownWeight {
			if c.Traits.Weight != traits.Weight {
				continue
			}
		}
		if traits.ShoeSize != characters.UnknownShoe {
			if c.Traits.ShoeSize != traits.ShoeSize {
				continue
			}
		}
		matches = append(matches, c)
	}
	if len(matches) == 0 {
		fmt.Println("No matches found for this dossier")
	} else {
		fmt.Println("Matches:")
		for _, match := range matches {
			fmt.Printf(" - %s\n", match.GetName())
		}
	}
	return false, nil
}

type dossierMenu string

const (
	dName      dossierMenu = "Dossier Name"
	dCharacter dossierMenu = "Character"
	dNotes     dossierMenu = "Notes"
	dDone      dossierMenu = "Done"
)

func (dm dossierMenu) String() string {
	return string(dm)
}

var dossierMenuItems = []dossierMenu{
	dName, dCharacter, dNotes, dDone,
}

func updateDossier(scanner *bufio.Scanner, dossier *characters.Dossier) {
	isDone := false
	for !isDone {
		idx := gamelogic.MenuSelect(scanner, "Select an option:", lo.Map(dossierMenuItems, func(dm dossierMenu, _ int) string { return dm.String() }))
		option := dossierMenuItems[idx]
		switch option {
		case dName:
			fmt.Print("New Dossier Name: ")
			scanner.Scan()
			dossier.Name = scanner.Text()
		case dCharacter:
			updateCharacter(scanner, dossier)
		case dNotes:
			updateNotes(scanner, dossier)
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

func (cm characterMenu) String() string {
	return string(cm)
}

var characterMenuItems = []characterMenu{
	cName, cGender, cNationality, cHeight, cWeight, cEyeColor, cHairColor, cHairLength, cShoeSize, cDone,
}

var clueGenders = append(characters.Genders, characters.UnknownGender)
var clueNationalities = append(characters.Nationalities, characters.UnknownNationality)
var clueHeights = append(characters.Heights, characters.UnknownHeight)
var clueWeights = append(characters.Weights, characters.UnknownWeight)
var clueEyeColors = append(characters.EyeColors, characters.UnknownEyes)
var clueHairColors = append(characters.HairColors, characters.UnknownHairColor)
var clueHairLengths = append(characters.HairLengths, characters.UnknownHairLength)
var clueShoeSizes = append(characters.ShoeSizes, characters.UnknownShoe)

func updateCharacter(scanner *bufio.Scanner, dossier *characters.Dossier) {

	isDone := false
	for !isDone {
		idx := gamelogic.MenuSelect(scanner, "Select an option:", lo.Map(characterMenuItems, func(cm characterMenu, _ int) string { return cm.String() }))
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
			idx = gamelogic.MenuSelect(scanner, "Gender:", lo.Map(clueGenders, func(g characters.Gender, _ int) string { return g.String() }))
			dossier.Target.Traits.Gender = clueGenders[idx]
		case cNationality:
			idx = gamelogic.MenuSelect(scanner, "Nationality:", lo.Map(clueNationalities, func(n characters.Nationality, _ int) string { return n.String() }))
			dossier.Target.Traits.Nationality = clueNationalities[idx]
		case cHeight:
			idx = gamelogic.MenuSelect(scanner, "Height:", lo.Map(clueHeights, func(h characters.Height, _ int) string { return h.String() }))
			dossier.Target.Traits.Height = clueHeights[idx]
		case cWeight:
			idx = gamelogic.MenuSelect(scanner, "Weight:", lo.Map(clueWeights, func(w characters.Weight, _ int) string { return w.String() }))
			dossier.Target.Traits.Weight = clueWeights[idx]
		case cEyeColor:
			idx = gamelogic.MenuSelect(scanner, "Eye Color:", lo.Map(clueEyeColors, func(ec characters.EyeColor, _ int) string { return ec.String() }))
			dossier.Target.Traits.EyeColor = clueEyeColors[idx]
		case cHairColor:
			idx = gamelogic.MenuSelect(scanner, "Hair Color:", lo.Map(clueHairColors, func(hc characters.HairColor, _ int) string { return hc.String() }))
			dossier.Target.Traits.HairColor = clueHairColors[idx]
		case cHairLength:
			idx = gamelogic.MenuSelect(scanner, "Hair Length:", lo.Map(clueHairLengths, func(hl characters.HairLength, _ int) string { return hl.String() }))
			dossier.Target.Traits.HairLength = clueHairLengths[idx]
		case cShoeSize:
			idx = gamelogic.MenuSelect(scanner, "Shoe Size:", lo.Map(clueShoeSizes, func(ss characters.ShoeSize, _ int) string { return ss.String() }))
			dossier.Target.Traits.ShoeSize = clueShoeSizes[idx]
		case cDone:
			isDone = true
		}
	}
}

type noteMenu string

const (
	addNote    noteMenu = "Add"
	deleteNote noteMenu = "Delete"
	doneNote   noteMenu = "Done"
)

func (nm noteMenu) String() string {
	return string(nm)
}

var noteMenuItems = []noteMenu{
	addNote, deleteNote, doneNote,
}

func updateNotes(scanner *bufio.Scanner, dossier *characters.Dossier) {
	option := addNote
	for option != doneNote {
		if len(dossier.Notes) > 0 {
			idx := gamelogic.MenuSelect(scanner, "Select an option:", lo.Map(noteMenuItems, func(nm noteMenu, _ int) string { return nm.String() }))
			option = noteMenuItems[idx]
		}
		switch option {
		case addNote:
			fmt.Print("Note: ")
			scanner.Scan()
			dossier.AddNote(scanner.Text())
		case deleteNote:
			idx := gamelogic.MenuSelect(scanner, "Select note to delete:", dossier.Notes)
			dossier.Notes = slices.Delete(dossier.Notes, idx, idx)
		}
	}
}

func getDossierIndex(gs *gamelogic.GameState, params []string) int {

	idx := -1
	if len(params) < 1 {
		idx = gamelogic.MenuSelect(gs.Scanner, "Choose a Dossier:", lo.Map(gs.Player.Dossiers, func(d characters.Dossier, _ int) string { return d.Name }))
	} else {
		dossierName := params[0]
		idx = slices.IndexFunc(gs.Player.Dossiers, func(d characters.Dossier) bool {
			return strings.ToLower(d.Name) == dossierName
		})
	}
	return idx
}
