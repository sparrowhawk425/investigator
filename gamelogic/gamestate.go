package gamelogic

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	"github.com/samber/lo"

	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/gameobjects/enemies"
	"github.com/sparrowhawk425/investigators/times"
)

type GameState struct {
	Scanner   *bufio.Scanner
	Day       int
	TimeOfDay times.TimeOfDay
	Player    Player
	Places    []gameobjects.Location
	People    []gameobjects.Character
	Criminals []enemies.Enemy
	Crimes    []Crime
}

func (gs GameState) PrintDay() {
	fmt.Printf("Day: %d Time: %s\n", gs.Day, times.GetTimeOfDayName(gs.TimeOfDay))
}

func (gs GameState) GetTimeOfDay() times.TimeOfDay {
	return gs.TimeOfDay
}

func (gs GameState) GetLocationsByType(locTypes []gameobjects.LocationType) []gameobjects.Location {
	return lo.Filter(gs.Places, func(loc gameobjects.Location, i int) bool {
		return slices.Contains(locTypes, loc.Type)
	})
}

func (gs GameState) GetLocationsByLootType(loots []gameobjects.LootType) []gameobjects.Location {
	return lo.Filter(gs.Places, func(loc gameobjects.Location, i int) bool {
		for _, loot := range loots {
			if slices.Contains(loc.GetAvailableLoot(), loot) {
				return true
			}
		}
		return false
	})
}

func (gs *GameState) AddCharacterToLocation(location gameobjects.Location, character gameobjects.Character) {
	for i := range gs.Places {
		if gs.Places[i].Equals(location) {
			gs.Places[i].Visitors = append(gs.Places[i].Visitors, character)
			break
		}
	}
}

func (gs *GameState) CreateCrime(location gameobjects.Location, loot []gameobjects.Loot) {

	gs.Crimes = append(gs.Crimes, Crime{
		Day:        gs.Day,
		TimeOfDay:  gs.TimeOfDay,
		StolenLoot: loot,
	})
}

func (gs *GameState) Update() {

	// Reset location visitors
	for i := range gs.Places {
		gs.Places[i].Visitors = []gameobjects.Character{}
	}
	for i := range gs.Criminals {
		gs.Criminals[i].PerformAction(gs)
		if gs.Criminals[i].Goal.Progress >= gs.Criminals[i].Goal.Target {
			fmt.Printf("%s %s has gathered enough loot and gone to ground.\n", gs.Criminals[i].Character.Name.First, gs.Criminals[i].Character.Name.Last)
			fmt.Println("You have failed!")
			os.Exit(0)
		}
	}
	for _, place := range gs.Places {
		if len(place.Visitors) > 0 {
			for _, visitor := range place.Visitors {
				fmt.Printf("%s %s is visiting %d %s\n", visitor.Name.First, visitor.Name.Last, place.Address.Number, place.Address.Name)
			}
		}
	}
	gs.TimeOfDay = times.TransitionTimeOfDay(gs.TimeOfDay)
	if gs.TimeOfDay == times.Morning {
		gs.Day++
	}
}
