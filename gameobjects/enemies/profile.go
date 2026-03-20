package enemies

import (
	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/times"
)

type Profile struct {
	Name            string
	ActiveDuring    times.TimeOfDay
	SleepDuring     times.TimeOfDay
	TargetLocations []gameobjects.LocationType
	PreferredLoot   []gameobjects.Loot
	Solitary        bool

	Action Action
}

// Burglar
func createBurglar() Profile {
	return Profile{
		Name:         "Burglar",
		ActiveDuring: times.Night,
		SleepDuring:  times.Morning,
		TargetLocations: []gameobjects.LocationType{
			gameobjects.Residence, gameobjects.Store, gameobjects.Museum,
		},
		PreferredLoot: []gameobjects.Loot{
			gameobjects.Jewelry, gameobjects.Art,
		},
		Solitary: true,
		Action:   CreateBurgleAction(),
	}
}

// Robber
// Hacker
// Bruiser
// Vandal
// Fence
// Hitman
// Cleaner
// Ghost?
