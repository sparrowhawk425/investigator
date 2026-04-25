package characters

import (
	"slices"

	"github.com/sparrowhawk425/investigators/internal/gameobjects"
)

type Behavior struct {
	Name               string
	Desc               string
	QualityPreference  []gameobjects.Quality
	LocationPreference []gameobjects.LocationType
	LootAmountModifier int //percent
	RiskModifier       int
	ReconModifier      int
}

func (b Behavior) FindTarget(findTarget func([]gameobjects.Location) *gameobjects.Location) func([]gameobjects.Location) *gameobjects.Location {
	return func(locations []gameobjects.Location) *gameobjects.Location {
		options := []gameobjects.Location{}
		for _, location := range locations {
			options = append(options, location)
			// Add additional hit if the location preference matches
			if slices.Contains(b.LocationPreference, location.Type) {
				options = append(options, location)
			}
			// Add additional hit if the quality preference matches
			if slices.Contains(b.QualityPreference, location.GetQuality()) {
				options = append(options, location)
			}
		}
		return findTarget(options)
	}
}

func (b Behavior) GetLootAmount(getLootAmt func(int) int) func(int) int {
	return func(maxAmt int) int {
		modAmt := maxAmt * (50 + b.LootAmountModifier) / 100
		// Always return min of 1
		return getLootAmt(max(1, modAmt))
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
		Name:               "Frugal",
		Desc:               "Prone to conserving money and prefer cheap locations",
		QualityPreference:  []gameobjects.Quality{gameobjects.Cheap},
		LootAmountModifier: -20,
	}
}

func CreateProfligate() Behavior {
	return Behavior{
		Name:               "Profligate",
		Desc:               "Prone to spending money and prefer expensive locations",
		QualityPreference:  []gameobjects.Quality{gameobjects.Expensive},
		LootAmountModifier: 20,
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
