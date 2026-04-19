package characters

import (
	"math/rand/v2"

	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/times"
)

// Interface for the GameState to avoid circular import
type GameStateI interface {
	GetTimeOfDay() times.TimeOfDay
	GetDayOfTheWeek() times.DayOfTheWeek
	GetLocations() []gameobjects.Location
	GetLocationsByType(locTypes []gameobjects.LocationType) []gameobjects.Location
	GetLocationsByLootType(loots []gameobjects.LootType) []gameobjects.Location
	AddCharacterToLocation(location gameobjects.Location, character Character)
	CreateCrime(location gameobjects.Location, name string, loot []gameobjects.Loot)
	CreateClue(location gameobjects.Location, clue string)
	TransferItems(lootType gameobjects.LootType, amount int, src gameobjects.ItemHolder, dest gameobjects.ItemHolder)
	SetCriminalEscaping(person Character)
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

func CreateVisitingAction() Action {
	return Action{
		Name: "Visiting",
		Risk: 0,
		Act:  VisitAction,
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
		Risk: 15,
		Act:  ReconAction,
	}
}

func CreateBurgleAction() Action {
	return Action{
		Name: "Burglary",
		Risk: 20,
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

func CreateEscapeAction() Action {
	return Action{
		Name: "Escaping",
		Risk: 0,
		Act:  EscapeAction,
	}
}

// Perform Actions

func SleepAction(gs GameStateI, person *Character) {
	gs.AddCharacterToLocation(person.Address, *person)
}

func RestAction(gs GameStateI, person *Character) {
	gs.AddCharacterToLocation(person.Address, *person)
}

func GuardAction(gs GameStateI, person *Character) {
	gs.AddCharacterToLocation(*person.GetTarget(), *person)

	if gs.GetDayOfTheWeek() == times.Monday {
		person.AddItems(gameobjects.Money, 15)
	}
}

func BankingAction(gs GameStateI, person *Character) {
	gs.AddCharacterToLocation(*person.GetTarget(), *person)

	if gs.GetDayOfTheWeek() == times.Monday {
		person.AddItems(gameobjects.Money, 25)
	}
}

func ManagingAction(gs GameStateI, person *Character) {
	gs.AddCharacterToLocation(*person.GetTarget(), *person)

	if gs.GetDayOfTheWeek() == times.Monday {
		person.AddItems(gameobjects.Money, 50)
	}
}

func VisitAction(gs GameStateI, person *Character) {
	gs.AddCharacterToLocation(*person.GetIdleTarget(), *person)
	price := person.GetIdleTarget().GetAdmissionPrice()
	gs.TransferItems(gameobjects.Money, price, person, person.GetIdleTarget())
}

func SellingAction(gs GameStateI, person *Character) {

	gs.AddCharacterToLocation(*person.GetIdleTarget(), *person)
	money := person.GetIdleTarget().GetLootAmount(gameobjects.Money)
	for lootType, loot := range person.GetItems() {
		if lootType == gameobjects.Money {
			continue
		}
		quantity := getLootAmount(money, loot.Quantity, loot.Value)
		gs.TransferItems(lootType, loot.Quantity, person, person.GetIdleTarget())
		value := quantity * loot.Value
		gs.TransferItems(gameobjects.Money, value, person.GetIdleTarget(), person)
		money -= value
	}
}

func getLootAmount(max, quantity, value int) int {
	for quantity*value > max {
		quantity--
	}
	return quantity
}

// Criminal Actions

func LieLowAction(gs GameStateI, person *Character) {
	gs.AddCharacterToLocation(person.Address, *person)
}

func ReconAction(gs GameStateI, person *Character) {
	gs.AddCharacterToLocation(*person.GetTarget(), *person)
}

func BurgleAction(gs GameStateI, person *Character) {
	takeLoot(gs, "Burglary", person)

	// Enemy needs new target
	person.SetTarget(nil)
}

func RobAction(gs GameStateI, person *Character) {
	takeLoot(gs, "Robbery", person)

	person.SetTarget(nil)
}

func VandalizeAction(gs GameStateI, person *Character) {
}

func EscapeAction(gs GameStateI, person *Character) {
	gs.SetCriminalEscaping(*person)
}

// Helpers

// TODO: There seem to be an inordinate amount of zeros returned from rand...
func takeLoot(gs GameStateI, crime string, person *Character) {

	gs.AddCharacterToLocation(*person.GetTarget(), *person)
	stolenLoot := []gameobjects.Loot{}
	for _, lootType := range person.GetPreferredLoot() {
		maxLoot := person.GetTarget().GetLootAmount(lootType)
		// Take a random amount of the available loot
		if maxLoot > 0 {
			amt := rand.IntN(maxLoot + 1)
			if amt > 0 {
				gs.TransferItems(lootType, amt, person.GetTarget(), person)
				stolenLoot = append(stolenLoot, gameobjects.Loot{Type: lootType, Quantity: amt, Value: lootType.GetValue()})
			}
		}
	}
	if len(stolenLoot) > 0 {
		gs.CreateCrime(*person.GetTarget(), crime, stolenLoot)
		riskPct := person.Role.RoleAction.Risk + person.GetTarget().GetRiskPercent()
		percent := rand.IntN(101)
		if riskPct > percent {
			gs.CreateClue(*person.GetTarget(), person.CreateClue())
		}
	}
}
