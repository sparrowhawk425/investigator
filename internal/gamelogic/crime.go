package gamelogic

import (
	"github.com/sparrowhawk425/investigators/internal/characters"
	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/times"
)

// TODO: Collect witnesses when a crime is committed (need talk commands and logic)
type Crime struct {
	Day        int
	TimeOfDay  times.TimeOfDay
	Location   gameobjects.Location
	Type       string
	StolenLoot []gameobjects.Loot
	Witnesses  []characters.Character
}

func (c Crime) Print(printFunc func(string, ...any)) {
	printFunc("Day: %d, Time: %s\n", c.Day, c.TimeOfDay.GetName())
	printFunc("%s at %s\n", c.Type, c.Location.GetAddress())
	printFunc("Loot:\n")
	for _, loot := range c.StolenLoot {
		printFunc(" - %s: %d\n", loot.Type, loot.Quantity)
	}
	if len(c.Witnesses) > 0 {
		printFunc("Witnesses:\n")
		for _, w := range c.Witnesses {
			printFunc(" - %s\n", w.GetName())
		}
	}
}
