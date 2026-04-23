package characters

import (
	"math/rand/v2"
	"slices"

	"github.com/sparrowhawk425/investigators/internal/functions"
	"github.com/sparrowhawk425/investigators/internal/gameobjects"
)

type Behavior struct {
	Name               string
	Desc               string
	QualityPreference  []gameobjects.Quality
	LocationPreference []gameobjects.LocationType
	RiskModifier       int
	ReconModifier      int
}

func (b Behavior) FindTarget(findTarget func([]gameobjects.Location) *gameobjects.Location) func([]gameobjects.Location) *gameobjects.Location {
	return func(locations []gameobjects.Location) *gameobjects.Location {
		targets := locations
		prefChance := rand.IntN(100)
		if prefChance > 40 && len(b.QualityPreference) > 0 {
			qualMatches := functions.Filter(targets, func(loc gameobjects.Location, i int) bool {
				return slices.Contains(b.QualityPreference, loc.GetQuality())
			})
			if len(qualMatches) > 0 {
				targets = qualMatches
			}
		}
		prefChance = rand.IntN(100)
		if prefChance > 40 && len(b.LocationPreference) > 0 {
			locMatches := functions.Filter(targets, func(loc gameobjects.Location, i int) bool {
				return slices.Contains(b.LocationPreference, loc.Type)
			})
			if len(locMatches) > 0 {
				targets = locMatches
			}
		}
		// If we somehow end up with no matches, return original list
		if len(targets) == 0 {
			targets = locations
		}
		return findTarget(targets)
	}
}

func (b Behavior) GetRiskPercent(getRisk func() int) func() int {
	return func() int {
		return getRisk() + b.RiskModifier
	}
}

func (b Behavior) GetReconModifier(reconTimes func() int) func() int {
	return func() int {
		return reconTimes() + b.ReconModifier
	}
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

func CreateCautious() Behavior {
	return Behavior{
		Name:          "Cautious",
		Desc:          "Takes fewer risks and takes more time to reconnoiter",
		RiskModifier:  -10,
		ReconModifier: 3,
	}
}

func CreateReckless() Behavior {
	return Behavior{
		Name:          "Reckless",
		Desc:          "Takes more rists and takes less time to reconnoiter",
		RiskModifier:  10,
		ReconModifier: -1,
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
