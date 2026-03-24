package gamelogic

import (
	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/times"
)

type Crime struct {
	Day        int
	TimeOfDay  times.TimeOfDay
	StolenLoot []gameobjects.Loot
}
