package characters

import (
	"bufio"
	"fmt"

	"github.com/sparrowhawk425/investigators/internal/gameobjects"
)

type Dossier struct {
	Name    string
	Target  *Character
	Profile string
	Notes   []string
}

type Player struct {
	Character

	Name            string
	CurrentLocation *gameobjects.Location
	Dossiers        []Dossier
	Action          *Action
}

const unknown = "Unknown"

func (p *Player) CreateDossier(scanner *bufio.Scanner) *Dossier {
	fmt.Print("Choose a name for this Dossier > ")
	scanner.Scan()
	name := scanner.Text()

	target := Character{
		Traits: Characteristics{
			Dob:         DateOfBirth{},
			Nationality: UnknownNationality,
			Gender:      UnknownGender,
			Height:      UnknownHeight,
			Weight:      UnknownWeight,
			EyeColor:    UnknownEyes,
			HairColor:   UnknownHairColor,
			ShoeSize:    UnknownShoe,
			HairLength:  UnknownHairLength,
		},
	}
	target.SetName(unknown, unknown)
	dossier := Dossier{Name: name, Target: &target}
	p.Dossiers = append(p.Dossiers, dossier)
	return &p.Dossiers[len(p.Dossiers)-1]
}

func (d *Dossier) AddNote(note string) {
	d.Notes = append(d.Notes, note)
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
