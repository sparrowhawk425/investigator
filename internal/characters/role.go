package characters

import (
	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/times"
)

type Goal struct {
	Progress int
	Target   int
}

func (g Goal) IsComplete() bool {
	return g.Progress >= g.Target
}

type Role struct {
	Name         string
	ActiveDuring times.TimeOfDay
	SleepDuring  times.TimeOfDay

	targetLocations []gameobjects.LocationType
	preferredLoot   []gameobjects.LootType
	Solitary        bool
	Freelancer      bool

	target *gameobjects.Location
	Action Action
}

var RegularRoles = []Role{
	CreateDayGuard(), CreateNightGuard(), CreateBanker(), CreateManager(),
}

var CriminalRoles = []Role{
	CreateBurglar(), CreateRobber(), CreateVandal(), CreateFence(),
}

func CreateDayGuard() Role {
	return Role{
		Name:         "Guard",
		ActiveDuring: times.Afternoon,
		SleepDuring:  times.Night,
		targetLocations: []gameobjects.LocationType{
			gameobjects.Bank, gameobjects.Casino, gameobjects.Store, gameobjects.Business, gameobjects.Hotel,
		},
		preferredLoot: []gameobjects.LootType{
			gameobjects.Money,
		},
		Solitary:   false,
		Freelancer: false,
		Action:     CreateGuardAction(),
	}
}

func CreateNightGuard() Role {
	return Role{
		Name:         "Guard",
		ActiveDuring: times.Night,
		SleepDuring:  times.Morning,
		targetLocations: []gameobjects.LocationType{
			gameobjects.Bank, gameobjects.Casino, gameobjects.Hotel,
		},
		preferredLoot: []gameobjects.LootType{
			gameobjects.Money,
		},
		Solitary:   false,
		Freelancer: false,
		Action:     CreateGuardAction(),
	}
}

func CreateBanker() Role {
	return Role{
		Name:         "Banker",
		ActiveDuring: times.Afternoon,
		SleepDuring:  times.Night,
		targetLocations: []gameobjects.LocationType{
			gameobjects.Bank,
		},
		preferredLoot: []gameobjects.LootType{
			gameobjects.Money,
		},
		Solitary:   false,
		Freelancer: false,
		Action:     CreateBankingAction(),
	}
}

func CreateManager() Role {
	return Role{
		Name:         "Manager",
		ActiveDuring: times.Afternoon,
		SleepDuring:  times.Night,
		targetLocations: []gameobjects.LocationType{
			gameobjects.Bank, gameobjects.Casino, gameobjects.Store, gameobjects.Business,
		},
		preferredLoot: []gameobjects.LootType{
			gameobjects.Money,
		},
		Solitary:   false,
		Freelancer: false,
		Action:     CreateManagingAction(),
	}
}

// Criminal Roles

func CreateBurglar() Role {
	return Role{
		Name:         "Burglar",
		ActiveDuring: times.Night,
		SleepDuring:  times.Morning,
		targetLocations: []gameobjects.LocationType{
			gameobjects.Residence, gameobjects.Store, gameobjects.Museum, gameobjects.Business, gameobjects.Bank,
		},
		preferredLoot: []gameobjects.LootType{
			gameobjects.Jewelry, gameobjects.Art, gameobjects.Money,
		},
		Solitary:   true,
		Freelancer: true,
		Action:     CreateBurgleAction(),
	}
}

func CreateRobber() Role {
	return Role{
		Name:         "Robber",
		ActiveDuring: times.Afternoon,
		SleepDuring:  times.Night,
		targetLocations: []gameobjects.LocationType{
			gameobjects.Bank, gameobjects.Casino, gameobjects.Store,
		},
		preferredLoot: []gameobjects.LootType{
			gameobjects.Money, gameobjects.Jewelry, gameobjects.Electronics, gameobjects.Cars,
		},
		Solitary:   false,
		Freelancer: true,
		Action:     CreateRobAction(),
	}
}

func CreateVandal() Role {
	return Role{
		Name:         "Vandal",
		ActiveDuring: times.Night,
		SleepDuring:  times.Afternoon,
		targetLocations: []gameobjects.LocationType{
			gameobjects.Store, gameobjects.Business, gameobjects.Residence,
		},
		preferredLoot: []gameobjects.LootType{
			gameobjects.Electronics, gameobjects.Money, gameobjects.Cars,
		},
		Solitary:   true,
		Freelancer: true,
		Action:     CreateVandalizeAction(),
	}
}

func CreateFence() Role {
	return Role{
		Name:         "Fence",
		ActiveDuring: times.Afternoon,
		SleepDuring:  times.Night,
		targetLocations: []gameobjects.LocationType{
			gameobjects.Store, gameobjects.PawnShop,
		},
		preferredLoot: []gameobjects.LootType{
			gameobjects.Jewelry, gameobjects.Art, gameobjects.Cars, gameobjects.Electronics,
		},
		Solitary:   true,
		Freelancer: true,
		Action:     CreateFenceAction(),
	}
}

// Hacker
// Bruiser
// Hitman
// Cleaner
// Ghost?
