package characters

import (
	"slices"

	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/gameobjects"
)

type Behavior struct {
	Name               string
	Desc               string
	QualityPreference  []gameobjects.Quality
	LocationPreference []gameobjects.LocationType
}

func (b Behavior) FilterLocations(locations []gameobjects.Location) []gameobjects.Location {
	targets := functions.Filter(locations, func(loc gameobjects.Location, i int) bool {
		return slices.Contains(b.QualityPreference, loc.GetQuality())
	})
	targets = functions.Filter(targets, func(loc gameobjects.Location, i int) bool {
		return slices.Contains(b.LocationPreference, loc.Type)
	})
	return targets
}

func CreateFrugal() Behavior {
	return Behavior{
		Name:              "Frugal",
		Desc:              "Prone to conserving money and prefer cheap locations",
		QualityPreference: []gameobjects.Quality{gameobjects.Cheap},
	}
}

func CreateProfligate() Behavior {
	return Behavior{
		Name:              "Profligate",
		Desc:              "Prone to spending money and prefer expensive locations",
		QualityPreference: []gameobjects.Quality{gameobjects.Expensive},
	}
}

func CreateGambler() Behavior {
	return Behavior{
		Name:               "Gambler",
		Desc:               "Tends to spend free time in Casinos. More willing to take risks",
		LocationPreference: []gameobjects.LocationType{gameobjects.Casino},
	}
}

var RegularBehaviors = []Behavior{
	CreateFrugal(), CreateProfligate(), CreateGambler(),
}

func CreateSquatter() Behavior {
	return Behavior{
		Name: "Squatter",
		Desc: "Lives in unoccupied buildings",
	}
}
