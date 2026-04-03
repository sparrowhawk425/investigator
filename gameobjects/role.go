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

	target *Location
	Action Action
}

type CreateRole func() Role

var RegularRoles = []CreateRole{
	CreateDayGuard, CreateNightGuard, CreateBanker, CreateManager,
}

var CriminalRoles = []CreateRole{
	CreateBurglar, CreateRobber, CreateVandal, CreateFence,
}

func CreateDayGuard() Role {
	return Role{
		Name:         "Guard",
		ActiveDuring: times.Afternoon,
		SleepDuring:  times.Night,
		targetLocations: []LocationType{
			Bank, Casino, Store, Business, Hotel,
		},
		preferredLoot: []LootType{
			Money,
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
		targetLocations: []LocationType{
			Bank, Casino, Hotel,
		},
		preferredLoot: []LootType{
			Money,
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
		targetLocations: []LocationType{
			Bank,
		},
		preferredLoot: []LootType{
			Money,
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
		targetLocations: []LocationType{
			Bank, Casino, Store, Business,
		},
		preferredLoot: []LootType{
			Money,
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
		targetLocations: []LocationType{
			Residence, Store, Museum, Business, Bank,
		},
		preferredLoot: []LootType{
			Jewelry, Art, Money,
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
		targetLocations: []LocationType{
			Bank, Casino, Store,
		},
		preferredLoot: []LootType{
			Money, Jewelry, Electronics, Cars,
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
		targetLocations: []LocationType{
			Store, Business, Residence,
		},
		preferredLoot: []LootType{
			Electronics, Money, Cars,
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
		targetLocations: []LocationType{
			Store, PawnShop,
		},
		preferredLoot: []LootType{
			Jewelry, Art, Cars, Electronics,
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
