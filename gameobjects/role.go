package gameobjects

import (
	"github.com/sparrowhawk425/investigators/times"
)

type Role struct {
	Name         string
	ActiveDuring times.TimeOfDay
	SleepDuring  times.TimeOfDay

	targetLocations []LocationType
	preferredLoot   []LootType
	Solitary        bool
	Freelancer      bool

	target      *Location
	possessions []Loot
	Action      Action
}

func CreateNightGuard() Role {
	return Role{
		Name:         "Guard",
		ActiveDuring: times.Night,
		SleepDuring:  times.Morning,
		targetLocations: []LocationType{
			Bank, Casino,
		},
		Solitary:   false,
		Freelancer: false,
		Action:     CreateGuardAction(),
	}
}

// Criminal Roles

func CreateBurglar() Role {
	return Role{
		Name:         "Burglar",
		ActiveDuring: times.Night,
		SleepDuring:  times.Morning,
		targetLocations: []LocationType{
			Residence, Store, Museum,
		},
		preferredLoot: []LootType{
			Jewelry, Art, Money,
		},
		Solitary:   true,
		Freelancer: true,
		Action:     CreateBurgleAction(),
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
