package enemies

import (
	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/times"
)

type Profile struct {
	Name            string
	ActiveDuring    times.TimeOfDay
	TargetLocations []gameobjects.LocationType
	PreferredLoot   []gameobjects.Loot
	Solitary        bool
}

// Burglar
func createBurglar() Profile {
	return Profile{
		Name:         "Burglar",
		ActiveDuring: times.Night,
		TargetLocations: []gameobjects.LocationType{
			gameobjects.Residence, gameobjects.Store, gameobjects.Museum,
		},
		PreferredLoot: []gameobjects.Loot{
			gameobjects.Jewelry, gameobjects.Art,
		},
		Solitary: true,
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
