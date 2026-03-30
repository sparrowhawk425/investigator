package gamelogic

import (
	"fmt"

	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/times"
)

type Crime struct {
	Day        int
	TimeOfDay  times.TimeOfDay
	Location   gameobjects.Location
	StolenLoot []gameobjects.Loot
}

func (c Crime) Print() {
	fmt.Printf("Day: %d, Time: %s\n", c.Day, c.TimeOfDay.GetName())
	fmt.Println(c.Location.GetAddress())
	fmt.Println("Loot:")
	for _, loot := range c.StolenLoot {
		fmt.Printf("\t%s: %d\n", loot.Type, loot.Quantity)
	}
}
