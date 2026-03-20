package enemies

import (
	"fmt"

	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/times"
)

// Interface for the GameState to avoid circular import
type HasLocations interface {
	GetTimeOfDay() times.TimeOfDay
	GetLocationsByType(locTypes []gameobjects.LocationType) []gameobjects.Location
	GetLocationsByLoot(loots []gameobjects.Loot) []gameobjects.Location
	AddCharacterToLocation(location gameobjects.Location, character gameobjects.Character)
}
type Action struct {
	Name string
	Risk int //percent
	Act  func(*HasLocations, *Enemy)
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

func SleepAction(locations *HasLocations, enemy *Enemy) {
	fmt.Println("Sleeping...")
}

func LieLowAction(locations *HasLocations, enemy *Enemy) {
	fmt.Println("Lying low...")
}

func ReconAction(locations *HasLocations, enemy *Enemy) {
	fmt.Println("Performing recon...")

	(*locations).AddCharacterToLocation(*(enemy.Target), enemy.Character)
}

func BurgleAction(locations *HasLocations, enemy *Enemy) {
	fmt.Println("Burgling...")

	// Enemy needs new target
	enemy.Target = nil
}
