package gamelogic

import (
	"fmt"

	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/sparrowhawk425/investigators/internal/times"
)

type Crime struct {
	Day        int
	TimeOfDay  times.TimeOfDay
	Location   gameobjects.Location
	Type       string
	StolenLoot []gameobjects.Loot
}

func (c Crime) Print() {
	fmt.Printf("Day: %d, Time: %s\n", c.Day, c.TimeOfDay.GetName())
	fmt.Printf("%s at %s\n", c.Type, c.Location.GetAddress())
	fmt.Println("Loot:")
	for _, loot := range c.StolenLoot {
		fmt.Printf(" - %s: %d\n", loot.Type, loot.Quantity)
	}
}
