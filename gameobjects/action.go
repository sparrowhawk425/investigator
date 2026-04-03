package gameobjects

import (
	"fmt"
	"math/rand/v2"

	"github.com/sparrowhawk425/investigators/times"
)

// Interface for the GameState to avoid circular import
type HasLocations interface {
	GetTimeOfDay() times.TimeOfDay
	GetLocationsByType(locTypes []LocationType) []Location
	GetLocationsByLootType(loots []LootType) []Location
	AddCharacterToLocation(location Location, character Character)
	CreateCrime(location Location, loot []Loot)
}
type Action struct {
	Name string
	Risk int //percent
	Act  func(HasLocations, *Character)
}

func CreateSleepAction() Action {
	return Action{
		Name: "Sleep",
		Risk: 5,
		Act:  SleepAction,
	}
}

func CreateLieLowAction() Action {
	return Action{
		Name: "Lie Low",
		Risk: 5,
		Act:  LieLowAction,
	}
}

func CreateReconAction() Action {
	return Action{
		Name: "Recon",
		Risk: 20,
		Act:  ReconAction,
	}
}

func CreateBurgleAction() Action {
	return Action{
		Name: "Burgle",
		Risk: 30,
		Act:  BurgleAction,
	}
}

func CreateGuardAction() Action {
	return Action{
		Name: "Guard",
		Risk: 0,
		Act:  GuardAction,
	}
}

// PerformActions

func SleepAction(locations HasLocations, person *Character) {
	fmt.Println("Sleeping...")
}

func GuardAction(locations HasLocations, person *Character) {
	fmt.Println("Guarding...")
}

func LieLowAction(locations HasLocations, person *Character) {
	fmt.Println("Lying low...")
}

func ReconAction(locations HasLocations, person *Character) {
	fmt.Println("Performing recon...")

	locations.AddCharacterToLocation(*person.GetTarget(), *person)
}

// TODO: There seem to be an inordinate amount of zeros returned from rand...
func BurgleAction(locations HasLocations, person *Character) {
	fmt.Println("Burgling...")

	locations.AddCharacterToLocation(*person.GetTarget(), *person)
	stolenLoot := []Loot{}
	for _, lootType := range person.GetPreferredLoot() {
		maxLoot := person.GetTarget().GetLootAmount(lootType)
		// Take a random amount of the available loot
		if maxLoot > 0 {
			amt := rand.IntN(maxLoot + 1)
			if amt > 0 {
				person.GetTarget().UpdateLoot(lootType, -1*amt)
				person.UpdatePossessions(lootType, amt)
				stolenLoot = append(stolenLoot, Loot{
					Type:     lootType,
					Quantity: amt,
				})
			}
		}
	}
	if len(stolenLoot) > 0 {
		locations.CreateCrime(*person.GetTarget(), stolenLoot)
		riskPct := person.Role.Action.Risk + person.GetTarget().GetRiskPercent()
		num := rand.IntN(100) + 1
		if riskPct > num {
			person.GetTarget().AddClue(person.CreateClue())
		}
	}

	// Enemy needs new target
	person.SetTarget(nil)
}
