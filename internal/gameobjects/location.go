package gameobjects

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"

	"github.com/samber/lo"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
)

// Cycle mapping for Character
type Person interface {
	GetName() string
	GetFirstName() string
	GetLastName() string
}

type LocationType string

const (
	Residence  LocationType = "Residence"
	Hotel      LocationType = "Hotel"
	Store      LocationType = "Store"
	Bank       LocationType = "Bank"
	Museum     LocationType = "Museum"
	Business   LocationType = "Business"
	Casino     LocationType = "Casino"
	PawnShop   LocationType = "Pawn Shop"
	Restaurant LocationType = "Restaurant"
)

var LocationTypes = []LocationType{
	Residence, Hotel, Store, Bank, Museum, Business, Casino, PawnShop,
}

func (LocationType) IsType() bool {
	return true
}

func (lt LocationType) String() string {
	return string(lt)
}

func GetLocationType(locTypeStr string) (LocationType, error) {
	types := lo.Map(LocationTypes, func(lt LocationType, i int) string {
		return strings.ToLower(string(lt))
	})
	idx := slices.Index(types, locTypeStr)
	if idx != -1 {
		return LocationTypes[idx], nil
	}
	return "", fmt.Errorf("Unexpected LocationType: %s", locTypeStr)
}

type Address struct {
	Number int
	Name   string
}

type Quality int

const (
	Cheap Quality = iota
	Moderate
	Expensive
)

var QualityTypes = []Quality{
	Cheap, Moderate, Expensive,
}

func (Quality) IsType() bool {
	return true
}

func (q Quality) String() string {
	switch q {
	case Cheap:
		return "Cheap"
	case Moderate:
		return "Moderate"
	case Expensive:
		return "Expensive"
	default:
		return "Indescribable"
	}
}

type Location struct {
	Type     LocationType
	Address  Address
	City     string
	State    string
	Country  string
	PostCode string

	IsOccupied bool
	Visitors   []Person

	quality Quality
	money   int
	loot    []Loot
	clues   []string
}

func (loc Location) Equals(other Location) bool {
	if loc.Type != other.Type {
		return false
	}
	if loc.Address.Number != other.Address.Number {
		return false
	}
	if loc.Address.Name != other.Address.Name {
		return false
	}
	if loc.City != other.City {
		return false
	}
	if loc.State != other.State {
		return false
	}
	if loc.Country != other.Country {
		return false
	}
	if loc.PostCode != other.PostCode {
		return false
	}
	return true
}

func (loc Location) Print() {
	fmt.Printf("%s\n", loc.GetAddress())

	if len(loc.GetAvailableLoot()) == 0 {
		fmt.Println("Notable Loot: None")
	} else {
		fmt.Println("Notable Loot:")
		for _, loot := range loc.GetAvailableLoot() {
			fmt.Printf(" - %s\n", loot)
		}
	}
	if len(loc.Visitors) > 0 {
		fmt.Println("People:")
		for _, person := range loc.Visitors {
			fmt.Printf(" - %s\n", person.GetName())
		}
	}
	if len(loc.GetClues()) > 0 {
		fmt.Println("Clues:")
		for _, clue := range loc.GetClues() {
			fmt.Printf(" - %s\n", clue)
		}
	}
	fmt.Println("")
}

func (loc Location) GetAddress() string {
	desc := fmt.Sprintf("%s %s", loc.quality, loc.Type)
	return fmt.Sprintf("%s: %d %s, %s, %s", desc, loc.Address.Number, loc.Address.Name, loc.City, loc.State)
}

func (loc Location) GetQuality() Quality {
	return loc.quality
}

func (loc Location) GetQualityStr() string {
	return loc.quality.String()
}

func (loc Location) GetAvailableLoot() []LootType {
	available := []LootType{}
	for _, loot := range loc.loot {
		if loot.Quantity > 0 {
			available = append(available, loot.Type)
		}
	}
	return available
}

func (loc Location) GetLootAmount(lootType LootType) int {
	for _, availableLoot := range loc.loot {
		if availableLoot.Type == lootType {
			return availableLoot.Quantity
		}
	}
	return 0
}

func (loc *Location) AddLoot(lootType LootType, amount int) {
	loc.UpdateLoot(lootType, amount)
}

func (loc *Location) GiveLoot(lootType LootType, amount int) Loot {
	value := loc.UpdateLoot(lootType, -1*amount)
	return Loot{Type: lootType, Quantity: amount, Value: value}
}

func (loc *Location) UpdateLoot(lootType LootType, amount int) int {
	for i := range loc.loot {
		if loc.loot[i].Type == lootType {
			loc.loot[i].Quantity += amount
			value := loc.loot[i].Value
			if loc.loot[i].Quantity == 0 {
				loc.loot = slices.Delete(loc.loot, i, i)
			}
			return value
		}
	}
	// If it wasn't found, add a new one
	if amount > 0 {
		value := lootType.GetValue()
		loc.loot = append(loc.loot, Loot{
			Type: lootType, Quantity: amount, Value: value,
		})
		return value
	}
	return 0
}

func (loc *Location) AddClue(clue string) {
	loc.clues = append(loc.clues, clue)
}

func (loc Location) GetClues() []string {
	return loc.clues
}

func (loc Location) GetRiskPercent() int {
	risk := 1
	switch loc.Type {
	case Residence:
		risk = 10
	case Store:
		risk = 15
	case Hotel:
		risk = 20
	case Bank:
		risk = 25
	case Museum:
		risk = 20
	case Business:
		risk = 15
	case Casino:
		risk = 30
	}
	switch loc.quality {
	case Moderate:
		risk *= 2
	case Expensive:
		risk *= 3
	}

	return risk
}

func (loc Location) String() string {
	return loc.GetAddress()
}

// Filters

func FilterLocationsByType(locTypes []LocationType) func(Location, int) bool {
	return func(loc Location, _ int) bool {
		return slices.Contains(locTypes, loc.Type)
	}
}

func FilterLocationsByLootType(lootTypes []LootType) func(Location, int) bool {
	return func(loc Location, _ int) bool {
		for _, lootType := range lootTypes {
			if slices.Contains(loc.GetAvailableLoot(), lootType) {
				return true
			}
		}
		return false
	}
}

func FilterLocationsByQuality(quality []Quality) func(Location, int) bool {
	return func(loc Location, _ int) bool {
		return slices.Contains(quality, loc.quality)
	}
}

// Helpers

func CreateLocation(fromLoc nameapi.Location, locType LocationType, isOccupied bool) Location {
	qual := Quality(rand.IntN(3))
	availableLoot := setAvailableLoot(locType, qual)
	return Location{
		Type: locType,
		Address: Address{
			Number: fromLoc.Street.Number,
			Name:   fromLoc.Street.Name,
		},
		City:       fromLoc.City,
		State:      fromLoc.State,
		Country:    fromLoc.Country,
		PostCode:   parsePostCode(fromLoc.Postcode),
		IsOccupied: isOccupied,
		quality:    qual,
		loot:       availableLoot,
	}
}

func CreateRandomLocations(apiLocations []nameapi.Location) []Location {
	locations := make([]Location, len(apiLocations))
	for i, apiLoc := range apiLocations {
		locType := LocationTypes[rand.IntN(len(LocationTypes))]
		occupiedPct := rand.IntN(100)
		locations[i] = CreateLocation(apiLoc, locType, occupiedPct > 5)
	}
	return locations
}

// Have to do special conversion because the API returns a string or an int
func parsePostCode(data []byte) string {
	if len(data) > 1 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}
	return string(data)
}

func setAvailableLoot(locType LocationType, quality Quality) []Loot {
	maxAmt := 1
	switch quality {
	case Cheap:
		maxAmt = 2
	case Moderate:
		maxAmt = 6
	case Expensive:
		maxAmt = 10
	}
	loot := []Loot{}
	switch locType {
	case Residence:
		loot = append(loot, Loot{Type: Jewelry, Quantity: rand.IntN(maxAmt), Value: Jewelry.GetValue()})
		loot = append(loot, Loot{Type: Money, Quantity: rand.IntN(maxAmt), Value: Money.GetValue()})
		loot = append(loot, Loot{Type: Electronics, Quantity: rand.IntN(maxAmt), Value: Electronics.GetValue()})
		loot = append(loot, Loot{Type: Cars, Quantity: rand.IntN(maxAmt), Value: Cars.GetValue()})
	case Bank:
		loot = append(loot, Loot{Type: Jewelry, Quantity: rand.IntN(maxAmt * 3), Value: Jewelry.GetValue()})
		loot = append(loot, Loot{Type: Money, Quantity: rand.IntN(maxAmt * 3), Value: Money.GetValue()})
	case Museum:
		loot = append(loot, Loot{Type: Jewelry, Quantity: rand.IntN(maxAmt * 2), Value: Jewelry.GetValue()})
		loot = append(loot, Loot{Type: Art, Quantity: rand.IntN(maxAmt * 2), Value: Art.GetValue()})
	case Hotel:
		loot = append(loot, Loot{Type: Jewelry, Quantity: rand.IntN(maxAmt * 2), Value: Jewelry.GetValue()})
		loot = append(loot, Loot{Type: Money, Quantity: rand.IntN(maxAmt), Value: Money.GetValue()})
		loot = append(loot, Loot{Type: Electronics, Quantity: rand.IntN(maxAmt * 2), Value: Electronics.GetValue()})
	case Store:
		loot = append(loot, Loot{Type: Money, Quantity: rand.IntN(maxAmt * 2), Value: Money.GetValue()})
		loot = append(loot, Loot{Type: Electronics, Quantity: rand.IntN(maxAmt * 3), Value: Electronics.GetValue()})
		loot = append(loot, Loot{Type: Cars, Quantity: rand.IntN(maxAmt * 3), Value: Cars.GetValue()})
	case Business:
		loot = append(loot, Loot{Type: Money, Quantity: rand.IntN(maxAmt), Value: Money.GetValue()})
		loot = append(loot, Loot{Type: Electronics, Quantity: rand.IntN(maxAmt * 4), Value: Electronics.GetValue()})
	case Casino:
		loot = append(loot, Loot{Type: Money, Quantity: rand.IntN(maxAmt * 4), Value: Money.GetValue()})
		loot = append(loot, Loot{Type: Jewelry, Quantity: rand.IntN(maxAmt * 4), Value: Jewelry.GetValue()})
	case Restaurant:
		loot = append(loot, Loot{Type: Money, Quantity: rand.IntN(maxAmt), Value: Money.GetValue()})
	case PawnShop:
		loot = append(loot, Loot{Type: Money, Quantity: rand.IntN(maxAmt * 2), Value: Money.GetValue()})
		loot = append(loot, Loot{Type: Jewelry, Quantity: rand.IntN(maxAmt * 3), Value: Jewelry.GetValue()})
		loot = append(loot, Loot{Type: Art, Quantity: rand.IntN(maxAmt * 2), Value: Art.GetValue()})
		loot = append(loot, Loot{Type: Electronics, Quantity: rand.IntN(maxAmt * 2), Value: Electronics.GetValue()})
	}
	return loot
}
