package characters

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/sparrowhawk425/investigators/internal/config"
	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
)

type Goal struct {
	Progress int
	Target   int
	Alerts   []string
}

func (g *Goal) Update(loot gameobjects.Loot) {
	startingPercent := 100 * g.Progress / g.Target
	if loot.Type == gameobjects.Money {
		g.Progress += loot.Value * loot.Quantity
		if g.Progress < 0 {
			g.Progress = 0
		}
	}
	curPercent := 100 * g.Progress / g.Target
	if curPercent != startingPercent {
		if startingPercent < 25 && curPercent >= 25 {
			g.Alerts = append(g.Alerts, "Our sources indicate a Syndicate member is about 25% done with their goal")
		} else if startingPercent < 50 && curPercent >= 50 {
			g.Alerts = append(g.Alerts, "Our sources indicate a Syndicate member is about 50% done with their goal")
		} else if startingPercent < 75 && curPercent >= 75 {
			g.Alerts = append(g.Alerts, "Our sources indicate a Syndicate member is about 75% done with their goal")
		}
	}

}

func (g Goal) IsComplete() bool {
	return g.Progress >= g.Target
}

type name struct {
	first string
	last  string
}

func (n name) Equals(other name) bool {
	return n.first == other.first && n.last == other.last
}

type CriminalPostCard []string

func (pc CriminalPostCard) Print() {
	fmt.Println("Post Card:")
	for _, line := range pc {
		fmt.Println(line)
	}
	fmt.Println("")
}

type Character struct {
	name
	Traits  Characteristics
	Address gameobjects.Location

	Role
	Behavior
	Goal

	possessions gameobjects.Inventory

	target     *gameobjects.Location
	idleTarget *gameobjects.Location
	FindTarget func([]gameobjects.Location) *gameobjects.Location

	GetLootAmount func(int) int
	GetRisk       func() int

	reconCount    int
	GetReconTimes func() int

	postCards []CriminalPostCard
}

func (c Character) Equals(other Character) bool {
	return c.name.Equals(other.name) && c.Address.Equals(other.Address)
}

func (c Character) String() string {
	return c.GetName()
}

func getLootAmount(maxAmount int) int {
	return rand.IntN(maxAmount)
}

func getReconTimes() int {
	return 1
}

func findTarget(locations []gameobjects.Location) *gameobjects.Location {
	if len(locations) == 0 {
		slog.Warn("No valid locations found for target")
		return nil
	}
	return &locations[rand.IntN(len(locations))]
}

func (c Character) getRisk() int {
	return c.Role.RoleAction.Risk
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

func (c Character) Print(printFunc func(string, ...any)) {

	printFunc("%s\n", c.GetName())
	printFunc(" - Gender: %s\n", c.Traits.Gender)
	printFunc(" - Nationality: %s\n", c.Traits.Nationality)
	printFunc(" - Height: %s\n", c.Traits.Height)
	printFunc(" - Weight: %s\n", c.Traits.Weight)
	printFunc(" - Eye Color: %s\n", c.Traits.EyeColor)
	printFunc(" - Hair Color: %s\n", c.Traits.HairColor)
	printFunc(" - Hair Length: %s\n", c.Traits.HairLength)
	printFunc(" - Shoe Size: %s\n", c.Traits.ShoeSize)
	printFunc(" - Address: %d %s\n", c.Address.Address.Number, c.Address.Address.Name) // TODO: Instead of listing their address, make it easier to track a criminal (Need to make it possible to find criminals elsewhere)?

	printFunc("")
}

func (c Character) GetPreferredLoot() []gameobjects.LootType {
	return c.Role.preferredLoot
}

func (c Character) HasTarget() bool {
	return c.target != nil
}

func (c *Character) GetTarget() *gameobjects.Location {
	return c.target
}

func (c *Character) SetTarget(loc *gameobjects.Location) {
	c.target = loc
}

func (c *Character) GetIdleLocation() *gameobjects.Location {
	return c.idleTarget
}

func (c Character) GetJobLocation() *gameobjects.Location {
	return c.Role.JobLocation
}

func (c *Character) SetJobLocation(location *gameobjects.Location) {
	c.Role.JobLocation = location
}

func (c Character) GetItems() gameobjects.Inventory {
	return c.possessions
}

func (c *Character) AddItems(lootType gameobjects.LootType, amount int) {
	c.updatePossessions(lootType, amount)
}

func (c *Character) RemoveItems(lootType gameobjects.LootType, amount int) {
	c.updatePossessions(lootType, -amount)
}

func (c *Character) updatePossessions(lootType gameobjects.LootType, amount int) {
	loot, ok := c.possessions[lootType]
	if !ok {
		c.possessions[lootType] = gameobjects.Loot{Type: lootType, Value: lootType.GetValue()}
		loot = c.possessions[lootType]
	}
	loot.Quantity += amount
	c.possessions[lootType] = loot

	c.Goal.Update(loot)
}

func (c *Character) PerformAction(gs GameStateI) {
	// Select an appropriate action
	action := c.selectAction(gs)
	// Perform the selected action
	action.Act(gs, c)
}

func (c *Character) selectAction(gs GameStateI) Action {

	isCriminal := slices.ContainsFunc(CriminalRoles, func(r Role) bool { return c.Role.Name == r.Name })

	// If it's time to sleep
	if gs.GetTimeOfDay() == c.Role.SleepDuring {
		return CreateSleepAction()
	}
	// If it's time to take action
	if gs.GetTimeOfDay() == c.Role.ActiveDuring {
		// If they are a Syndicate member and completed their goal, escape!
		if isCriminal && c.Goal.IsComplete() {
			return CreateEscapeAction()
		}
		// Find a target
		if c.target == nil {
			c.target = c.FindTarget(gs.GetLocations())
			c.reconCount = 0
		}
		// Check if should perform recon
		if c.ShouldRecon() {
			return CreateReconAction()
		}
		// Perform role action
		return c.Role.RoleAction
	}
	// Spare time, so recreate, maybe
	recreation := rand.IntN(100)
	if recreation < 34 {
		// Shop
		shops := gameobjects.ShopLocations
		if isCriminal {
			shops = []gameobjects.LocationType{gameobjects.PawnShop, gameobjects.Fence}
		}
		c.idleTarget = c.FindTarget(gs.GetLocationsByType(shops))
		return CreateSellingAction()
	} else if recreation < 67 {
		c.idleTarget = c.FindTarget(gs.GetLocationsByType(gameobjects.RecreationLocations))
		return CreateVisitingAction()
	}
	// Not doing anything else, so rest
	return c.Role.IdleAction
}

func (c *Character) ShouldRecon() bool {
	if c.reconCount < c.GetReconTimes() {
		c.reconCount++
		return true
	}
	return false
}

func (c Character) GetPostCards() []CriminalPostCard {
	return c.postCards
}

func (c *Character) SetPostCards(postCards []CriminalPostCard) {
	c.postCards = postCards
}

func (c *Character) AddPostCard(postCard CriminalPostCard) {
	c.postCards = append(c.postCards, postCard)
}

func CreateRandomCharacter(apiChar nameapi.Character, role Role) Character {

	goal := Goal{Target: 1500}
	eyeColor := EyeColors[rand.IntN(len(EyeColors))]
	hairColor := HairColors[rand.IntN(len(HairColors))]
	behavior := RegularBehaviors[rand.IntN(len(RegularBehaviors))]
	gender := Gender(strings.ToUpper(apiChar.Gender[:1]) + strings.ToLower(apiChar.Gender[1:]))
	if rand.IntN(100) > 94 {
		gender = NonbinaryGender
	}
	resType := gameobjects.Residence
	if rand.IntN(100) > 80 {
		resType = gameobjects.Hotel
	}
	c := Character{
		name: name{
			first: apiChar.Name.First,
			last:  apiChar.Name.Last,
		},
		Traits: Characteristics{
			Dob: DateOfBirth{
				Date: apiChar.DateOfBirth.Date,
				Age:  apiChar.DateOfBirth.Age,
			},
			Nationality: config.CountryCode(apiChar.Nationality),
			EyeColor:    eyeColor,
			HairColor:   hairColor,
			Gender:      gender,
			Height:      getHeight(apiChar.Gender),
			Weight:      getWeight(apiChar.Gender),
			ShoeSize:    ShoeSizes[rand.IntN(len(ShoeSizes))],
			HairLength:  HairLengths[rand.IntN(len(HairLengths))],
		},
		Address:     gameobjects.CreateLocation(apiChar.Location, resType, true),
		Goal:        goal,
		Role:        role,
		Behavior:    behavior,
		possessions: make(map[gameobjects.LootType]gameobjects.Loot),
	}
	// A little kludgy, but it does allow me to wrap the methods in one place
	c.FindTarget = behavior.FindTarget(role.FindTarget(findTarget))
	c.GetLootAmount = behavior.GetLootAmount(getLootAmount)
	c.GetRisk = behavior.GetRiskPercent(c.getRisk)
	c.GetReconTimes = behavior.GetReconModifier(getReconTimes)
	return c
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
	case clueHairColor, clueHairLength:
		if c.Traits.HairLength == BaldHair {
			return "A shot from a security camera shows the figure is bald"
		} else {
			return fmt.Sprintf("There's a %s %s stray hair recovered from the crime scene.", c.Traits.HairLength, c.Traits.HairColor)
		}
	case clueHeight:
		return fmt.Sprintf("A security guard says the figure's height was %s", c.Traits.Height)
	case clueWeight:
		switch c.Traits.Weight {
		case ThinWeight:
			return fmt.Sprintf("A security guard says the figure was %s", c.Traits.Weight)
		case AverageWeight:
			return fmt.Sprintf("A security guard says the figure was %s weight", c.Traits.Weight)
		case OverWeight:
			return fmt.Sprintf("A security guard says the figure was %sset", c.Traits.Weight)
		default:
			return "The figure was of indeterminable weight"
		}
	case clueShoeSize:
		return fmt.Sprintf("A footprint recovered from the crime scene reveals the figure had %s shoes", c.Traits.ShoeSize)
	default:
		// Shouldn't get here
		return "There's nothing to find"
	}
}
