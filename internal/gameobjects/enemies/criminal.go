package enemies

import (
	"github.com/sparrowhawk425/investigators/gameobjects"
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
	Personality []PersonalityTrait
	WorksFor    []Enemy

	Target *gameobjects.Location
}
