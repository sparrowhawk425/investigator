package gameobjects

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
)

type LocationType string

const (
	Residence LocationType = "Residence"
	Hotel     LocationType = "Hotel"
	Store     LocationType = "Store"
	Bank      LocationType = "Bank"
	Museum    LocationType = "Museum"
	Business  LocationType = "Business"
)

var locationTypes = []LocationType{
	Residence, Hotel, Store, Bank, Museum, Business,
}

func GetLocationType(locTypeStr string) (LocationType, error) {
	types := lo.Map(locationTypes, func(lt LocationType, i int) string {
		return strings.ToLower(string(lt))
	})
	idx := slices.Index(types, locTypeStr)
	if idx != -1 {
		return locationTypes[idx], nil
	}
	return "", fmt.Errorf("Unexpected LocationType: %s", locTypeStr)
}

type Address struct {
	Number int
	Name   string
}

type availableLoot struct {
	loot     Loot
	quantity int
}

type Quality int

const (
	cheap Quality = iota
	moderate
	expensive
)

type Location struct {
	Type     LocationType
	Address  Address
	City     string
	State    string
	Country  string
	PostCode string

	quality Quality
	loot    []availableLoot
}

func (loc Location) GetAvailableLoot() []Loot {
	available := []Loot{}
	for _, loot := range loc.loot {
		if loot.quantity > 0 {
			available = append(available, loot.loot)
		}
	}
	return available
}

func (loc Location) GetQuality() string {
	switch loc.quality {
	case cheap:
		return "Cheap"
	case moderate:
		return "Moderate"
	case expensive:
		return "Expensive"
	default:
		return "Indescribable"
	}
}

func CreateLocation(fromLoc nameapi.Location, locType LocationType) Location {
	qual := Quality(rand.IntN(3))
	availableLoot := setAvailableLoot(locType, qual)
	postCode := getPostalCode(fromLoc.Postcode)
	return Location{
		Type: locType,
		Address: Address{
			Number: fromLoc.Street.Number,
			Name:   fromLoc.Street.Name,
		},
		City:     fromLoc.City,
		State:    fromLoc.State,
		Country:  fromLoc.Country,
		PostCode: postCode,
		quality:  qual,
		loot:     availableLoot,
	}
}

func CreateRandomLocations(apiLocations []nameapi.Location) []Location {
	locations := make([]Location, len(apiLocations))
	for i, apiLoc := range apiLocations {
		locType := locationTypes[rand.IntN(len(locationTypes))]
		locations[i] = CreateLocation(apiLoc, locType)
	}
	return locations
}

// Have to do special conversion because the API returns a string or an int, but interface{} treats the ints as float64...
func getPostalCode(postCode any) string {
	strCode, ok := postCode.(string)
	if ok {
		return strCode
	}
	intCode, ok := postCode.(int)
	if ok {
		return strconv.Itoa(intCode)
	}
	floatCode, ok := postCode.(float64)
	if ok {
		return strconv.FormatFloat(floatCode, 'f', -1, 64)
	}
	return "12345"
}

func setAvailableLoot(locType LocationType, quality Quality) []availableLoot {
	maxAmt := 1
	switch quality {
	case cheap:
		maxAmt = 2
	case moderate:
		maxAmt = 6
	case expensive:
		maxAmt = 10
	}
	loot := []availableLoot{}
	switch locType {
	case Residence:
		loot = append(loot, availableLoot{loot: Jewelry, quantity: rand.IntN(maxAmt)})
		loot = append(loot, availableLoot{loot: Money, quantity: rand.IntN(maxAmt)})
		loot = append(loot, availableLoot{loot: Electronics, quantity: rand.IntN(maxAmt)})
		loot = append(loot, availableLoot{loot: Cars, quantity: rand.IntN(maxAmt)})
	case Bank:
		loot = append(loot, availableLoot{loot: Jewelry, quantity: rand.IntN(maxAmt * 3)})
		loot = append(loot, availableLoot{loot: Money, quantity: rand.IntN(maxAmt * 3)})
	case Museum:
		loot = append(loot, availableLoot{loot: Jewelry, quantity: rand.IntN(maxAmt * 2)})
		loot = append(loot, availableLoot{loot: Art, quantity: rand.IntN(maxAmt * 2)})
	case Hotel:
		loot = append(loot, availableLoot{loot: Jewelry, quantity: rand.IntN(maxAmt)})
		loot = append(loot, availableLoot{loot: Money, quantity: rand.IntN(maxAmt * 2)})
		loot = append(loot, availableLoot{loot: Electronics, quantity: rand.IntN(maxAmt * 2)})
	case Store:
		loot = append(loot, availableLoot{loot: Money, quantity: rand.IntN(maxAmt * 3)})
		loot = append(loot, availableLoot{loot: Electronics, quantity: rand.IntN(maxAmt * 3)})
		loot = append(loot, availableLoot{loot: Cars, quantity: rand.IntN(maxAmt * 3)})
	case Business:
		loot = append(loot, availableLoot{loot: Money, quantity: rand.IntN(maxAmt)})
		loot = append(loot, availableLoot{loot: Electronics, quantity: rand.IntN(maxAmt * 4)})
	}
	return loot
}
