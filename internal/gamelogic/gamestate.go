package gamelogic

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"slices"

	"github.com/samber/lo"

	"github.com/sparrowhawk425/investigators/internal/characters"
	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
	"github.com/sparrowhawk425/investigators/internal/times"
)

type GameState struct {
	Scanner   *bufio.Scanner
	DayNumber int
	WeekDay   times.DayOfTheWeek
	TimeOfDay times.TimeOfDay
	Player    characters.Player
	Places    []gameobjects.Location
	People    []characters.Character
	Criminals []characters.Character
	Crimes    []Crime
}

func (gs GameState) PrintDay() {
	fmt.Printf("\n%s, Day: %d Time: %s\n", gs.WeekDay, gs.DayNumber, gs.TimeOfDay.GetName())
}

func (gs GameState) GetTimeOfDay() times.TimeOfDay {
	return gs.TimeOfDay
}

func (gs *GameState) NextDay() {
	gs.DayNumber++
	gs.WeekDay = gs.WeekDay.NextDay()
}

func (gs *GameState) BuildGame() {

	// Select country
	countryNames := lo.Map(nameapi.Countries, func(country nameapi.Country, i int) string { return country.Name })
	idx := MenuSelect(gs.Scanner, "Select a Country to begin your investigation:", countryNames)
	country := nameapi.Countries[idx]
	fmt.Printf("Travelling to %s...\n", country.Name)
	// Add locations and people to game
	results, err := nameapi.MakeHTTPGetRequest(country, 20)
	if err != nil {
		log.Fatalf("Error getting locations from API: %v", err)
	}
	// Create people
	for _, c := range results {
		gs.People = append(gs.People, characters.CreateRandomCharacter(c))
	}
	// Create locations
	apiLocations := lo.Map(results, func(character nameapi.Character, i int) nameapi.Location { return character.Location })
	gs.Places = gameobjects.CreateRandomLocations(apiLocations)

	// Set Work Targets
	for i := range gs.People {
		gs.People[i].FindTarget(gs)
	}

	cIdxs := []int{}
	for range 2 {
		num := -1
		for num < 0 || slices.Contains(cIdxs, num) {
			num = rand.IntN(len(gs.People))
		}
		gs.People[num].Role = characters.CreateBurglar() //characters.CriminalRoles[rand.IntN(len(characters.CriminalRoles))]
		gs.Criminals = append(gs.Criminals, gs.People[num])
		cIdxs = append(cIdxs, num)
	}
}

func (gs GameState) GetLocationsByType(locTypes []gameobjects.LocationType) []gameobjects.Location {
	filters := []func(gameobjects.Location, int) bool{}
	filters = append(filters, gameobjects.FilterLocationsByType(locTypes))
	return gs.GetLocations(filters)
}

func (gs GameState) GetLocationsByLootType(lootTypes []gameobjects.LootType) []gameobjects.Location {
	filters := []func(gameobjects.Location, int) bool{}
	filters = append(filters, gameobjects.FilterLocationsByLootType(lootTypes))
	return gs.GetLocations(filters)
}

func (gs GameState) GetLocations(filters []func(gameobjects.Location, int) bool) []gameobjects.Location {
	locations := gs.Places
	for _, filter := range filters {
		locations = lo.Filter(locations, filter)
	}
	return locations
}

func (gs *GameState) AddCharacterToLocation(location gameobjects.Location, character characters.Character) {
	for i := range gs.Places {
		if gs.Places[i].Equals(location) {
			gs.Places[i].Visitors = append(gs.Places[i].Visitors, character)
			break
		}
	}
}

func (gs *GameState) CreateCrime(location gameobjects.Location, name string, loot []gameobjects.Loot) {

	gs.Crimes = append(gs.Crimes, Crime{
		Day:        gs.DayNumber,
		TimeOfDay:  gs.TimeOfDay,
		Location:   location,
		Type:       name,
		StolenLoot: loot,
	})
}

func (gs *GameState) Update() {

	// Reset location visitors
	for i := range gs.Places {
		gs.Places[i].Visitors = nil
	}
	for i := range gs.People {
		gs.People[i].PerformAction(gs)
		// if gs.People[i].Goal.Progress >= gs.Criminals[i].Goal.Target {
		// 	fmt.Printf("%s has gathered enough loot and gone to ground.\n", gs.Criminals[i].Character.GetName())
		// 	fmt.Println("You have failed!")
		// 	os.Exit(0)
		// }
	}
	if gs.Player.CurrentLocation != nil {
		fmt.Printf("%s is currently at:\n", gs.Player.Name)
		gs.Player.CurrentLocation.Print()
	}
	for _, place := range gs.Places {
		if len(place.Visitors) > 0 {
			for _, visitor := range place.Visitors {
				if visitor != nil {
					fmt.Printf("%s is visiting %d %s\n", visitor.GetName(), place.Address.Number, place.Address.Name)
				} else {
					fmt.Printf("Visitor to %s is nil\n", place.GetAddress())
				}
			}
		}
	}
	gs.TimeOfDay = times.TransitionTimeOfDay(gs.TimeOfDay)
	if gs.TimeOfDay == times.Morning {
		gs.NextDay()
	}
}
