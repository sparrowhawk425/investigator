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

func (e Enemy) GetPreferredLoot() []gameobjects.Loot {
	loot := []gameobjects.Loot{}
	for _, profLoot := range e.Profile.PreferredLoot {
		if !slices.Contains(loot, profLoot) {
			loot = append(loot, profLoot)
		}

	}
	return loot
}

func (e *Enemy) PerformAction(gs HasLocations) {
	if gs.GetTimeOfDay() == e.Profile.SleepDuring {
		e.Action = CreateSleepAction()
	} else if gs.GetTimeOfDay() != e.Profile.ActiveDuring {
		e.Action = CreateLieLowAction()
	} else if !e.HasTarget() {
		// If no target is currently selected, find a desirable target and assign it
		e.Action = CreateReconAction()
		targets := gs.GetLocationsByLoot(e.GetPreferredLoot())
		target := targets[rand.IntN(len(targets))]
		e.Target = &target
	} else {
		e.Action = e.Profile.Action
	}
	// Perform the selected action
	e.Action.Act(&gs, e)
}

func CreateEnemy(char nameapi.Character) Enemy {
	return Enemy{
		Character: gameobjects.CreateRandomCharacter(char),
		Address:   gameobjects.CreateLocation(char.Location, gameobjects.Hotel),
		Profile:   createBurglar(),
		Personality: []PersonalityTrait{
			createProfligate(),
		},
		Freelancer: true,
	}
}
