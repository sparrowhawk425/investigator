package gamelogic

import (
	"fmt"

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

func (c Crime) Print() {
	fmt.Printf("Day: %d, Time: %s\n", c.Day, c.TimeOfDay.GetName())
	fmt.Printf("%s at %s\n", c.Type, c.Location.GetAddress())
	fmt.Println("Loot:")
	for _, loot := range c.StolenLoot {
		fmt.Printf(" - %s: %d\n", loot.Type, loot.Quantity)
	}
	if len(c.Witnesses) > 0 {
		fmt.Println("Witnesses:")
		for _, w := range c.Witnesses {
			fmt.Printf(" - %s\n", w.GetName())
		}
	}
}
