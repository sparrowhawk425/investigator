package gamelogic

import (
	"bufio"
	"fmt"
	"strconv"

	"github.com/samber/lo"
	"github.com/sparrowhawk425/investigators/internal/characters"
	"github.com/sparrowhawk425/investigators/internal/gameobjects"
)

type Dossier struct {
	Name    string
	Target  *characters.Character
	Profile string
	Notes   []string
}

type Player struct {
	Name            string
	CurrentLocation gameobjects.Location
	Dossiers        []Dossier
}

const unknown = "Unknown"

func (p *Player) CreateDossier(scanner *bufio.Scanner) {
	fmt.Print("Choose a name for this Dossier > ")
	scanner.Scan()
	name := scanner.Text()

	target := characters.Character{
		Traits: characters.Characteristics{
			Dob:         characters.DateOfBirth{},
			Nationality: unknown,
			Gender:      unknown,
			Height:      characters.UnknownHeight,
			Weight:      characters.UnknownWeight,
			EyeColor:    characters.UnknownEyes,
			HairColor:   characters.UnknownHairColor,
			ShoeSize:    characters.UnknownShoe,
			HairLength:  characters.UnknownHairLength,
		},
	}
	target.SetName(unknown, unknown)
	dossier := Dossier{Name: name, Target: &target}
	dossier.Update(scanner)
	p.Dossiers = append(p.Dossiers, dossier)
}

func (d *Dossier) AddNote(note string) {
	d.Notes = append(d.Notes, note)
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

func (d *Dossier) Update(scanner *bufio.Scanner) {
	isDone := false
	for !isDone {
		idx := MenuSelect(scanner, "Choose Field:", lo.Map(dossierMenuItems, func(m dossierMenu, _ int) string { return string(m) }))
		option := dossierMenuItems[idx]
		switch option {
		case dName:
			fmt.Print("New Name: ")
			scanner.Scan()
			d.Name = scanner.Text()
		case dCharacter:
			d.UpdateCharacter(scanner)
		case dNotes:
			fmt.Print("Note: ")
			scanner.Scan()
			d.AddNote(scanner.Text())
		case dDone:
			isDone = true
		}
	}
}

type characterMenu string

const (
	cName        characterMenu = "Name"
	cGender      characterMenu = "Gender"
	cAge         characterMenu = "Age"
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
	cName, cGender, cAge, cNationality, cHeight, cWeight, cEyeColor, cHairColor, cHairLength, cShoeSize, cDone,
}

func (d *Dossier) UpdateCharacter(scanner *bufio.Scanner) {

	isDone := false
	for !isDone {
		idx := MenuSelect(scanner, "Choose Field:", lo.Map(characterMenuItems, func(m characterMenu, _ int) string { return string(m) }))
		option := characterMenuItems[idx]
		switch option {
		case cName:
			fmt.Print("First Name: ")
			scanner.Scan()
			d.Target.SetFirstName(scanner.Text())
			fmt.Print("Last Name: ")
			scanner.Scan()
			d.Target.SetLastName(scanner.Text())
		case cGender:
			fmt.Print("Gender: ")
			scanner.Scan()
			d.Target.Traits.Gender = scanner.Text()
		case cAge:
			d.Target.Traits.Dob.Age = getNum(scanner, "Age: ")
		case cNationality:
			fmt.Print("Nationality: ")
			scanner.Scan()
			d.Target.Traits.Nationality = scanner.Text()
		case cHeight:
			fmt.Print("Height: ")
			scanner.Scan()
			d.Target.Traits.Height = characters.Height(scanner.Text())
		case cWeight:
			fmt.Print("Weight: ")
			scanner.Scan()
			d.Target.Traits.Weight = characters.Weight(scanner.Text())
		case cEyeColor:
			idx = MenuSelect(scanner, "Eye Color:", lo.Map(characters.EyeColors, func(ec characters.EyeColor, _ int) string { return string(ec) }))
			d.Target.Traits.EyeColor = characters.EyeColors[idx]
		case cHairColor:
			idx = MenuSelect(scanner, "Hair Color:", lo.Map(characters.HairColors, func(hc characters.HairColor, _ int) string { return string(hc) }))
			d.Target.Traits.HairColor = characters.HairColors[idx]
		case cHairLength:
			idx = MenuSelect(scanner, "Hair Length:", lo.Map(characters.HairLengths, func(hl characters.HairLength, _ int) string { return string(hl) }))
			d.Target.Traits.HairLength = characters.HairLengths[idx]
		case cShoeSize:
			idx = MenuSelect(scanner, "Shoe Size:", lo.Map(characters.ShoeSizes, func(ss characters.ShoeSize, _ int) string { return string(ss) }))
			d.Target.Traits.ShoeSize = characters.ShoeSizes[idx]
		case cDone:
			isDone = true
		}
	}
}

func getNum(scanner *bufio.Scanner, prompt string) int {
	for {
		fmt.Print(prompt)
		scanner.Scan()
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Not a valid number")
			continue
		}
		return val
	}
}

func (d Dossier) Print() {
	fmt.Printf("Dossier %s\n", d.Name)
	d.PrintCharacter()
	d.PrintNotes()
}

func (d Dossier) PrintCharacter() {
	fmt.Println("Target:")
	if d.Target == nil {
		fmt.Println("  No target set")
		return
	}
	fmt.Printf("  First Name: %s\n", d.Target.GetFirstName())
	fmt.Printf("  Last Name: %s\n", d.Target.GetLastName())
	fmt.Printf("  Gender: %s\n", d.Target.Traits.Gender)
	fmt.Printf("  Nationality: %s\n", d.Target.Traits.Nationality)
	if d.Target.Traits.Dob.Age == 0 {
		fmt.Println("  Age: Unknown")
	} else {
		fmt.Printf("  Age: %d\n", d.Target.Traits.Dob.Age)
	}
	fmt.Printf("  Height: %s\n", d.Target.Traits.Height)
	fmt.Printf("  Weight: %s\n", d.Target.Traits.Weight)
	fmt.Printf("  Eye Color: %s\n", d.Target.Traits.EyeColor)
	fmt.Printf("  Hair Color: %s\n", d.Target.Traits.HairColor)
	fmt.Printf("  Hair Length: %s\n", d.Target.Traits.HairLength)
	fmt.Printf("  Shoe Size: %s\n", d.Target.Traits.ShoeSize)
}

func (d Dossier) PrintNotes() {
	fmt.Println("Notes:")
	for _, note := range d.Notes {
		fmt.Printf(" - %s\n", note)
	}
}
