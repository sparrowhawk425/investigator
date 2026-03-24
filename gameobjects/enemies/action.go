package enemies

import (
	"fmt"
	"math/rand/v2"

	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/times"
)

// Interface for the GameState to avoid circular import
type HasLocations interface {
	GetTimeOfDay() times.TimeOfDay
	GetLocationsByType(locTypes []gameobjects.LocationType) []gameobjects.Location
	GetLocationsByLootType(loots []gameobjects.LootType) []gameobjects.Location
	AddCharacterToLocation(location gameobjects.Location, character gameobjects.Character)
	CreateCrime(location gameobjects.Location, loot []gameobjects.Loot)
}
type Action struct {
	Name string
	Risk int //percent
	Act  func(HasLocations, *Enemy)
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

func SleepAction(locations HasLocations, enemy *Enemy) {
	fmt.Println("Sleeping...")
}

func LieLowAction(locations HasLocations, enemy *Enemy) {
	fmt.Println("Lying low...")
}

func ReconAction(locations HasLocations, enemy *Enemy) {
	fmt.Println("Performing recon...")

	locations.AddCharacterToLocation(*(enemy.Target), enemy.Character)
}

// TODO: There seem to be an inordinate amount of zeros returned from rand...
func BurgleAction(locations HasLocations, enemy *Enemy) {
	fmt.Println("Burgling...")

	locations.AddCharacterToLocation(*(enemy.Target), enemy.Character)
	stolenLoot := []gameobjects.Loot{}
	for _, lootType := range enemy.GetPreferredLoot() {
		maxLoot := enemy.Target.GetLootAmount(lootType)
		// Take a random amount of the available loot
		if maxLoot > 0 {
			amt := rand.IntN(maxLoot + 1)
			enemy.Target.UpdateLoot(lootType, -1*amt)
			enemy.UpdateLoot(lootType, amt)
			stolenLoot = append(stolenLoot, gameobjects.Loot{
				Type:     lootType,
				Quantity: amt,
			})
		}
	}
	locations.CreateCrime(*enemy.Target, stolenLoot)
	// Enemy needs new target
	enemy.Target = nil
}
