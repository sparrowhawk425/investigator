package characters

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
)

type EyeColor string

const (
	BlueEyes    EyeColor = "Blue"
	GreenEyes   EyeColor = "Green"
	BrownEyes   EyeColor = "Brown"
	AmberEyes   EyeColor = "Amber"
	HazelEyes   EyeColor = "Hazel"
	GrayEyes    EyeColor = "Gray"
	UnknownEyes EyeColor = "Unknown"
	// Red/Violet (due to albinism)?
	// Heterochromia?
)

var EyeColors = []EyeColor{
	BlueEyes, GreenEyes, BrownEyes, AmberEyes, HazelEyes, GrayEyes,
}

type HairColor string

const (
	BlondHair        HairColor = "Blond"
	DarkBlondHair    HairColor = "Dark Blond"
	MediumBrownHair  HairColor = "Medium Brown"
	DarkBrownHair    HairColor = "Dark Brown"
	BlackHair        HairColor = "Black"
	AuburnHair       HairColor = "Auburn"
	RedHair          HairColor = "Red"
	GrayHair         HairColor = "Gray"
	WhiteHair        HairColor = "White"
	UnknownHairColor HairColor = "Unknown"
)

var HairColors = []HairColor{
	BlondHair, DarkBlondHair, MediumBrownHair, DarkBrownHair, BlackHair, AuburnHair, RedHair, GrayHair, WhiteHair,
}

type ShoeSize string

const (
	SmallShoe   ShoeSize = "Small"
	MediumShoe  ShoeSize = "Medium"
	LargeShoe   ShoeSize = "Large"
	UnknownShoe ShoeSize = "Unknown"
)

var ShoeSizes = []ShoeSize{
	SmallShoe, MediumShoe, LargeShoe,
}

type HairLength string

const (
	BaldHair          HairLength = "Bald"
	ShortHair         HairLength = "Short"
	MediumHair        HairLength = "Medium"
	LongHair          HairLength = "Long"
	UnknownHairLength HairLength = "Unknown"
)

var HairLengths = []HairLength{
	BaldHair, ShortHair, MediumHair, LongHair,
}

type Height string

const (
	ShortHeight   Height = "Short"
	AverageHeight Height = "Average"
	TallHeight    Height = "Tall"
	UnknownHeight Height = "Unknown"
)

type Weight string

const (
	ThinWeight    Weight = "Thin"
	AverageWeight Weight = "Average"
	OverWeight    Weight = "Heavy"
	UnknownWeight Weight = "Unknown"
)

type DateOfBirth struct {
	Date time.Time `json:"date"`
	Age  int
}

type Characteristics struct {
	Dob         DateOfBirth
	Nationality string
	Gender      string
	EyeColor    EyeColor
	HairColor   HairColor
	Height      Height
	Weight      Weight
	ShoeSize    ShoeSize
	HairLength  HairLength
}

type Goal struct {
	Progress int
	Target   int
}

func (g Goal) IsComplete() bool {
	return g.Progress >= g.Target
}

type name struct {
	first string
	last  string
}

type Character struct {
	name    name
	Traits  Characteristics
	Address gameobjects.Location

	Role        Role
	Behavior    Behavior
	Goal        Goal
	possessions []gameobjects.Loot
}

func (c Character) GetName() string {
	return fmt.Sprintf("%s %s", c.name.first, c.name.last)
}

func (c Character) GetFirstName() string {
	return c.name.first
}

func (c Character) GetLastName() string {
	return c.name.last
}

func (c *Character) SetName(first, last string) {
	c.name.first = first
	c.name.last = last
}

func (c *Character) SetFirstName(first string) {
	c.name.first = first
}

func (c *Character) SetLastName(last string) {
	c.name.last = last
}

func (c Character) Print() {

	fmt.Printf("%s\n", c.GetName())
	fmt.Printf(" - Gender: %s\n", c.Traits.Gender)
	fmt.Printf(" - Nationality: %s\n", c.Traits.Nationality)
	//fmt.Printf(" - Age: %d\n", c.Traits.Dob.Age)
	fmt.Printf(" - Height: %s\n", c.Traits.Height)
	fmt.Printf(" - Weight: %s\n", c.Traits.Weight)
	fmt.Printf(" - Eye Color: %s\n", c.Traits.EyeColor)
	fmt.Printf(" - Hair Color: %s\n", c.Traits.HairColor)
	fmt.Printf(" - Hair Length: %s\n", c.Traits.HairLength)
	fmt.Printf(" - Shoe Size: %s\n", c.Traits.ShoeSize)

	if len(c.possessions) > 0 {
		fmt.Println("\nItems:")
		for _, item := range c.possessions {
			fmt.Printf("%d %s\n", item.Quantity, item.Type)
		}
	}
	fmt.Println("")
}

func (c Character) GetPreferredLoot() []gameobjects.LootType {
	return c.Role.preferredLoot
}

func (c Character) HasTarget() bool {
	return c.Role.target != nil
}

func (c *Character) GetTarget() *gameobjects.Location {
	return c.Role.target
}

func (c *Character) SetTarget(loc *gameobjects.Location) {
	c.Role.target = loc
}

// TODO: This is redundant with ActionFSM (but on a different layer: role vs character)
func (c *Character) FindTarget(gs GameStateI) {

	targets := gs.GetLocationsByType(c.Role.targetLocations)
	targets = c.Behavior.FilterLocations(targets)
	if len(targets) == 0 {
		targets = gs.GetLocationsByLootType(c.Role.preferredLoot)
	}
	target := targets[rand.IntN(len(targets))]
	c.SetTarget(&target)
}

func (c Character) GetPossessions() []gameobjects.Loot {
	return c.possessions
}

func (c *Character) UpdatePossessions(loot gameobjects.Loot) {
	c.possessions = append(c.possessions, loot)
}

func (c *Character) PerformAction(gs GameStateI) {
	// Select an appropriate action
	action := c.selectAction(gs)
	// Perform the selected action
	action.Act(gs, c)
}

func (c *Character) selectAction(gs GameStateI) Action {

	// If it's time to sleep
	if gs.GetTimeOfDay() == c.Role.SleepDuring {
		return CreateSleepAction()
	}
	// If it's time to take action
	if gs.GetTimeOfDay() == c.Role.ActiveDuring {
		// If they have a target
		if c.Role.target != nil {
			return c.Role.RoleAction
		}
		// Find a target and perform recon
		targets := gs.GetLocationsByType(c.Role.targetLocations)
		targets = c.Behavior.FilterLocations(targets)
		if len(targets) == 0 {
			targets = gs.GetLocationsByLootType(c.Role.preferredLoot)
		}
		target := targets[rand.IntN(len(targets))]
		c.Role.target = &target
		return CreateReconAction()
	}
	// Not doing anything else, so rest
	return c.Role.RestAction
}

func CreateRandomCharacter(apiChar nameapi.Character) Character {

	eyeColor := EyeColors[rand.IntN(len(EyeColors))]
	hairColor := HairColors[rand.IntN(len(HairColors))]
	return Character{
		name: name{
			first: apiChar.Name.First,
			last:  apiChar.Name.Last,
		},
		Traits: Characteristics{
			Dob: DateOfBirth{
				Date: apiChar.DateOfBirth.Date,
				Age:  apiChar.DateOfBirth.Age,
			},
			Nationality: apiChar.Nationality,
			EyeColor:    eyeColor,
			HairColor:   hairColor,
			Gender:      apiChar.Gender,
			Height:      getHeight(apiChar.Gender),
			Weight:      getWeight(apiChar.Gender),
			ShoeSize:    ShoeSizes[rand.IntN(len(ShoeSizes))],
			HairLength:  HairLengths[rand.IntN(len(HairLengths))],
		},
		Role:     RegularRoles[rand.IntN(len(RegularRoles))],
		Behavior: RegularBehaviors[rand.IntN(len(RegularBehaviors))],
	}
}

// Avg Male weight: 200lb
// Avg Female weight: 170lb
func getWeight(gender string) Weight {
	avg := 170
	if strings.ToLower(gender) == "male" {
		avg = 200
	}
	w := calcDiff(avg)
	if w < avg {
		return ThinWeight
	} else if w > avg {
		return OverWeight
	}
	return AverageWeight
}

// Avg Male height: 171cm (5'7")
// Avg Female height: 159cm (5'3")
func getHeight(gender string) Height {
	avg := 159
	if strings.ToLower(gender) == "male" {
		avg = 171
	}
	h := calcDiff(avg)
	if h < avg {
		return ShortHeight
	} else if h > avg {
		return TallHeight
	}
	return AverageHeight
}

func calcDiff(avg int) int {
	pctDiff := rand.IntN(11) // 0-10%
	diff := avg * pctDiff / 100
	if rand.IntN(2)%2 == 0 { // + or -
		return avg + diff
	}
	return avg - diff
}

type clue int

const (
	clueGender clue = iota
	clueEyeColor
	clueHairColor
	clueHairLength
	clueHeight
	clueWeight
	clueShoeSize
)

var clueItems = []clue{
	clueGender, clueEyeColor, clueHairColor, clueHairLength, clueHeight, clueWeight, clueShoeSize,
}

func (c Character) CreateClue() string {
	clueItem := clueItems[rand.IntN(len(clueItems))]
	switch clueItem {
	case clueGender:
		return fmt.Sprintf("A shot from a security camera clearly shows the figure is %s.", c.Traits.Gender)
	case clueEyeColor:
		return fmt.Sprintf("A security guard got a look at the figure's face and saw they had %s eyes.", c.Traits.EyeColor)
	case clueHairColor | clueHairLength:
		if c.Traits.HairLength == BaldHair {
			return "A shot from a security camera shows the figure is bald"
		} else {
			return fmt.Sprintf("There's a %s %s stray hair recovered from the crime scene.", c.Traits.HairLength, c.Traits.HairColor)
		}
	case clueHeight:
		return ""
	case clueWeight:
		return ""
	case clueShoeSize:
		return fmt.Sprintf("A footprint recovered from the crime scene reveals the figure had %s shoes", c.Traits.ShoeSize)
	default:
		return "There's nothing to find"
	}
}
