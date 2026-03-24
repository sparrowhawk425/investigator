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
	PreferredLoot   []gameobjects.LootType
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
		PreferredLoot: []gameobjects.LootType{
			gameobjects.Jewelry, gameobjects.Art, gameobjects.Money,
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
