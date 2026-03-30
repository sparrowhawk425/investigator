package enemies

import (
	"math/rand/v2"
	"slices"

	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/nameapi"
)

type Goal struct {
	Progress int
	Target   int
}

func (g Goal) IsComplete() bool {
	return g.Progress >= g.Target
}

type Enemy struct {
	Goal        Goal
	Character   gameobjects.Character
	Address     gameobjects.Location
	Profile     Profile
	Personality []PersonalityTrait
	Freelancer  bool
	WorksFor    []Enemy

	Target *gameobjects.Location
	Action Action
}

func (e Enemy) HasTarget() bool {
	return e.Target != nil
}

func (e Enemy) GetPreferredLoot() []gameobjects.LootType {
	loot := []gameobjects.LootType{}
	for _, profLoot := range e.Profile.PreferredLoot {
		if !slices.Contains(loot, profLoot) {
			loot = append(loot, profLoot)
		}
	}
	return loot
}

func (e *Enemy) UpdateLoot(lootType gameobjects.LootType, amt int) {
	value := lootType.GetValue() * amt
	e.Goal.Progress += value
}

func (e *Enemy) PerformAction(gs HasLocations) {
	if gs.GetTimeOfDay() == e.Profile.SleepDuring {
		e.Action = CreateSleepAction()
	} else if gs.GetTimeOfDay() != e.Profile.ActiveDuring {
		e.Action = CreateLieLowAction()
	} else if !e.HasTarget() {
		// If no target is currently selected, find a desirable target and assign it
		e.Action = CreateReconAction()
		e.findTarget(gs)
	} else {
		e.Action = e.Profile.Action
	}
	// Perform the selected action
	e.Action.Act(gs, e)
}

func (e *Enemy) findTarget(gs HasLocations) {
	targets := gs.GetLocationsByLootType(e.GetPreferredLoot())
	target := targets[rand.IntN(len(targets))]
	e.Target = &target
}

func CreateEnemy(char nameapi.Character) Enemy {
	return Enemy{
		Character: gameobjects.CreateRandomCharacter(char),
		Address:   gameobjects.CreateLocation(char.Location, gameobjects.Hotel, true),
		Profile:   createBurglar(),
		Personality: []PersonalityTrait{
			createProfligate(),
		},
		Freelancer: true,
		Goal: Goal{
			Progress: 0,
			Target:   500,
		},
	}
}
