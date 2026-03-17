package gamelogic

import (
	"slices"

	"github.com/samber/lo"

	"github.com/sparrowhawk425/investigators/gameobjects"
	"github.com/sparrowhawk425/investigators/gameobjects/enemies"
)

type GameState struct {
	Places    []gameobjects.Location
	People    []gameobjects.Character
	Criminals []enemies.Enemy
}

func (gs GameState) GetLocationsByType(locTypes []gameobjects.LocationType) []gameobjects.Location {
	return lo.Filter(gs.Places, func(loc gameobjects.Location, i int) bool {
		return slices.Contains(locTypes, loc.Type)
	})
}

func (gs GameState) GetLocationsByLoot(loots []gameobjects.Loot) []gameobjects.Location {
	return lo.Filter(gs.Places, func(loc gameobjects.Location, i int) bool {
		for _, loot := range loots {
			if slices.Contains(loc.GetAvailableLoot(), loot) {
				return true
			}
		}
		return false
	})
}
