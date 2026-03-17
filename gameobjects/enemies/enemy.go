package enemies

import (
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
	Profiles    []Profile
	Personality []PersonalityTrait
	Freelancer  bool
	WorksFor    []Enemy
}

func CreateEnemy(char nameapi.Character) Enemy {
	return Enemy{
		Character: gameobjects.CreateRandomCharacter(char),
		Address:   gameobjects.CreateLocation(char.Location, gameobjects.Hotel),
		Profiles: []Profile{
			createBurglar(),
		},
		Personality: []PersonalityTrait{
			createProfligate(),
		},
		Freelancer: true,
	}
}
