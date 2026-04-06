package characters

import (
	"fmt"
	"math/rand/v2"

	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/times"
)

// Interface for the GameState to avoid circular import
type HasLocations interface {
	GetTimeOfDay() times.TimeOfDay
	GetLocationsByType(locTypes []gameobjects.LocationType) []gameobjects.Location
	GetLocationsByLootType(loots []gameobjects.LootType) []gameobjects.Location
	AddCharacterToLocation(location gameobjects.Location, character Character)
	CreateCrime(location gameobjects.Location, name string, loot []gameobjects.Loot)
}
type Action struct {
	Name string
	Risk int //percent
	Act  func(HasLocations, *Character)
}

func ActionFiniteStateMachine(role Role, trait Behavior) Action {

	return Action{}
}

func CreateSleepAction() Action {
	return Action{
		Name: "Sleeping",
		Risk: 5,
		Act:  SleepAction,
	}
}

// Regular Actions

func CreateGuardAction() Action {
	return Action{
		Name: "Guarding",
		Risk: 0,
		Act:  GuardAction,
	}
}

func CreateBankingAction() Action {
	return Action{
		Name: "Banking",
		Risk: 0,
		Act:  BankingAction,
	}
}

func CreateManagingAction() Action {
	return Action{
		Name: "Managing",
		Risk: 0,
		Act:  ManagingAction,
	}
}

// Criminal Actions

func CreateLieLowAction() Action {
	return Action{
		Name: "Lying Low",
		Risk: 5,
		Act:  LieLowAction,
	}
}

func CreateReconAction() Action {
	return Action{
		Name: "Performing Recon",
		Risk: 20,
		Act:  ReconAction,
	}
}

func CreateBurgleAction() Action {
	return Action{
		Name: "Burglary",
		Risk: 30,
		Act:  BurgleAction,
	}
}

func CreateRobAction() Action {
	return Action{
		Name: "Robbery",
		Risk: 40,
		Act:  RobAction,
	}
}

func CreateVandalizeAction() Action {
	return Action{
		Name: "Vandalism",
		Risk: 20,
		Act:  VandalizeAction,
	}
}

func CreateFenceAction() Action {
	return Action{
		Name: "Fencing",
		Risk: 15,
		Act:  FenceAction,
	}
}

// Perform Actions

func SleepAction(locations HasLocations, person *Character) {
	fmt.Println("Sleeping...")
	locations.AddCharacterToLocation(person.Address, *person)
}

func GuardAction(locations HasLocations, person *Character) {
	fmt.Println("Guarding...")
	locations.AddCharacterToLocation(*person.GetTarget(), *person)

}

func BankingAction(locations HasLocations, person *Character) {
	fmt.Println("Banking...")
	locations.AddCharacterToLocation(*person.GetTarget(), *person)
}

func ManagingAction(locations HasLocations, person *Character) {
	fmt.Println("Managing...")
	locations.AddCharacterToLocation(*person.GetTarget(), *person)
}

func LieLowAction(locations HasLocations, person *Character) {
	fmt.Println("Lying low...")
}

func ReconAction(locations HasLocations, person *Character) {
	fmt.Println("Performing recon...")

	locations.AddCharacterToLocation(*person.GetTarget(), *person)
}

func BurgleAction(locations HasLocations, person *Character) {
	fmt.Println("Burgling...")

	takeLoot(locations, "Burglary", person)

	// Enemy needs new target
	person.SetTarget(nil)
}

// TODO: There seem to be an inordinate amount of zeros returned from rand...
func takeLoot(locations HasLocations, crime string, person *Character) {

	locations.AddCharacterToLocation(*person.GetTarget(), *person)
	stolenLoot := []gameobjects.Loot{}
	for _, lootType := range person.GetPreferredLoot() {
		maxLoot := person.GetTarget().GetLootAmount(lootType)
		// Take a random amount of the available loot
		if maxLoot > 0 {
			amt := rand.IntN(maxLoot + 1)
			if amt > 0 {
				person.GetTarget().UpdateLoot(lootType, -1*amt)
				person.UpdatePossessions(lootType, amt)
				stolenLoot = append(stolenLoot, gameobjects.Loot{
					Type:     lootType,
					Quantity: amt,
				})
			}
		}
	}
	if len(stolenLoot) > 0 {
		locations.CreateCrime(*person.GetTarget(), crime, stolenLoot)
		riskPct := person.Role.Action.Risk + person.GetTarget().GetRiskPercent()
		num := rand.IntN(100) + 1
		if riskPct > num {
			person.GetTarget().AddClue(person.CreateClue())
		}
	}
}

func RobAction(locations HasLocations, person *Character) {
	fmt.Println("Robbing...")

	takeLoot(locations, "Robbery", person)

	person.SetTarget(nil)
}

func VandalizeAction(locations HasLocations, person *Character) {
	fmt.Println("Vandalizing...")

}

func FenceAction(locations HasLocations, person *Character) {
	fmt.Println("Fencing...")
}
