package characters

import (
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/times"
)

// Interface for the GameState to avoid circular import
type GameStateI interface {
	GetTimeOfDay() times.TimeOfDay
	GetLocations() []gameobjects.Location
	GetLocationsByType(locTypes []gameobjects.LocationType) []gameobjects.Location
	GetLocationsByLootType(loots []gameobjects.LootType) []gameobjects.Location
	AddCharacterToLocation(location gameobjects.Location, character Character)
	CreateCrime(location gameobjects.Location, name string, loot []gameobjects.Loot)
	CreateClue(location gameobjects.Location, clue string)
	TransferItems(lootType gameobjects.LootType, amount int, isStolen bool, src gameobjects.ItemHolder, dest gameobjects.ItemHolder)
}

// TODO: Refactor actions to be more modular? Having a create method and second internal action function seems redundant. The name doesn't actually do anything
type Action struct {
	Name string
	Risk int //percent
	Act  func(GameStateI, *Character)
}

func CreateSleepAction() Action {
	return Action{
		Name: "Sleeping",
		Risk: 5,
		Act:  SleepAction,
	}
}

// Regular Actions

func CreateRestAction() Action {
	return Action{
		Name: "Resting",
		Risk: 0,
		Act:  RestAction,
	}
}

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

func CreateEatingAction() Action {
	return Action{
		Name: "Eating",
		Risk: 0,
		Act:  EatingAction,
	}
}

func CreateSellingAction() Action {
	return Action{
		Name: "Selling",
		Risk: 0,
		Act:  SellingAction,
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

func CreateEscapeAction() Action {
	return Action{
		Name: "Escaping",
		Risk: 0,
		Act:  EscapeAction,
	}
}

// Perform Actions

func SleepAction(gs GameStateI, person *Character) {
	fmt.Println("Sleeping...")
	gs.AddCharacterToLocation(person.Address, *person)
}

func RestAction(gs GameStateI, person *Character) {
	fmt.Println("Resting...")
	gs.AddCharacterToLocation(person.Address, *person)
}

func GuardAction(gs GameStateI, person *Character) {
	fmt.Println("Guarding...")
	gs.AddCharacterToLocation(*person.GetTarget(), *person)

}

func BankingAction(gs GameStateI, person *Character) {
	fmt.Println("Banking...")
	gs.AddCharacterToLocation(*person.GetTarget(), *person)
}

func ManagingAction(gs GameStateI, person *Character) {
	fmt.Println("Managing...")
	gs.AddCharacterToLocation(*person.GetTarget(), *person)
}

func EatingAction(gs GameStateI, person *Character) {
	fmt.Println("Eating...")
	gs.AddCharacterToLocation(*person.GetTarget(), *person)
	gs.TransferItems(gameobjects.Money, 3, false, person, person.GetTarget())
}

func SellingAction(gs GameStateI, person *Character) {
	fmt.Println("Selling...")
	// TODO: How to get the right target?
}

// Criminal Actions

func LieLowAction(gs GameStateI, person *Character) {
	fmt.Println("Lying low...")
}

func ReconAction(gs GameStateI, person *Character) {
	fmt.Println("Performing recon...")

	gs.AddCharacterToLocation(*person.GetTarget(), *person)
}

func BurgleAction(gs GameStateI, person *Character) {
	log.Println("Burgling...")

	takeLoot(gs, "Burglary", person)

	// Enemy needs new target
	person.SetTarget(nil)
}

func RobAction(gs GameStateI, person *Character) {
	fmt.Println("Robbing...")

	takeLoot(gs, "Robbery", person)

	person.SetTarget(nil)
}

func VandalizeAction(gs GameStateI, person *Character) {
	fmt.Println("Vandalizing...")

}

func FenceAction(gs GameStateI, person *Character) {
	fmt.Println("Fencing...")

	gs.AddCharacterToLocation(person.Address, *person)
}

func EscapeAction(gs GameStateI, person *Character) {
	fmt.Println("Escaping...")

	fmt.Println("A member of the Syndicate has left the area...")
}

// Helpers

// TODO: There seem to be an inordinate amount of zeros returned from rand...
func takeLoot(gs GameStateI, crime string, person *Character) {

	fmt.Printf("%s be coming to take the loots\n", person.GetName())
	gs.AddCharacterToLocation(*person.GetTarget(), *person)
	stolenLoot := []gameobjects.Loot{}
	for _, lootType := range person.GetPreferredLoot() {
		maxLoot := person.GetTarget().GetLootAmount(lootType)
		// Take a random amount of the available loot
		if maxLoot > 0 {
			amt := rand.IntN(maxLoot + 1)
			if amt > 0 {
				fmt.Printf("%d %s be stolen from %s\n", amt, lootType, person.GetTarget().GetAddress())
				gs.TransferItems(lootType, amt, true, person.GetTarget(), person)
				stolenLoot = append(stolenLoot, gameobjects.Loot{Type: lootType, Quantity: amt, Value: lootType.GetValue(), IsStolen: true})
			}
		}
	}
	if len(stolenLoot) > 0 {
		gs.CreateCrime(*person.GetTarget(), crime, stolenLoot)
		riskPct := person.Role.RoleAction.Risk + person.GetTarget().GetRiskPercent()
		percent := rand.IntN(101)
		fmt.Printf("Exposure Risk: %d, Actual: %d\n", riskPct, percent)
		if riskPct > percent {
			gs.CreateClue(*person.GetTarget(), person.CreateClue())
		}
	}
}
