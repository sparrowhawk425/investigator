package characters

import (
	"slices"

	"github.com/sparrowhawk425/investigators/internal/gameobjects"
)

type Behaviors struct {
	Name               string
	Desc               string
	QualityPreference  []gameobjects.Quality
	LocationPreference []gameobjects.LocationType
	LootAmountModifier int //percent
	RiskModifier       int
	ReconModifier      int

	mutuallyExclusive []string
}

func (b Behaviors) FindTarget(findTarget func([]gameobjects.Location) *gameobjects.Location) func([]gameobjects.Location) *gameobjects.Location {
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

func (b Behaviors) GetLootAmount(getLootAmt func(int) int) func(int) int {
	return func(maxAmt int) int {
		modAmt := maxAmt * (50 + b.LootAmountModifier) / 100
		// Always return min of 1
		return getLootAmt(max(1, modAmt))
	}
}

func (b Behaviors) GetRiskPercent(getRisk func() int) func() int {
	return func() int {
		return getRisk() + b.RiskModifier
	}
}

func (b Behaviors) GetReconModifier(reconTimes func() int) func() int {
	return func() int {
		return reconTimes() + b.ReconModifier
	}
}

func FindTarget(behaviors []Behaviors, targetFunc func([]gameobjects.Location) *gameobjects.Location) func([]gameobjects.Location) *gameobjects.Location {
	if len(behaviors) == 0 {
		return targetFunc
	}
	return FindTarget(behaviors[1:], behaviors[0].FindTarget(targetFunc))
}

func GetLootAmount(behaviors []Behaviors, lootFunc func(int) int) func(int) int {
	if len(behaviors) == 0 {
		return lootFunc
	}
	return GetLootAmount(behaviors[1:], behaviors[0].GetLootAmount(lootFunc))
}

func GetRiskPercent(behaviors []Behaviors, riskFunc func() int) func() int {
	if len(behaviors) == 0 {
		return riskFunc
	}
	return GetRiskPercent(behaviors[1:], behaviors[0].GetRiskPercent(riskFunc))
}

func GetReconModifier(behaviors []Behaviors, reconFunc func() int) func() int {
	if len(behaviors) == 0 {
		return reconFunc
	}
	return GetReconModifier(behaviors[1:], behaviors[0].GetReconModifier(reconFunc))
}

func CreateFrugal() Behaviors {
	return Behaviors{
		Name:               "Frugal",
		Desc:               "Prone to conserving money and prefer cheap locations",
		QualityPreference:  []gameobjects.Quality{gameobjects.Cheap},
		LootAmountModifier: -20,
		mutuallyExclusive:  []string{"Profligate"},
	}
}

func CreateProfligate() Behaviors {
	return Behaviors{
		Name:               "Profligate",
		Desc:               "Prone to spending money and prefer expensive locations",
		QualityPreference:  []gameobjects.Quality{gameobjects.Expensive},
		LootAmountModifier: 20,
		mutuallyExclusive:  []string{"Frugal"},
	}
}

func CreateGambler() Behaviors {
	return Behaviors{
		Name:               "Gambler",
		Desc:               "Tends to spend free time in Casinos. More willing to take risks",
		LocationPreference: []gameobjects.LocationType{gameobjects.Casino},
		RiskModifier:       20,
	}
}

func CreateCautious() Behaviors {
	return Behaviors{
		Name:              "Cautious",
		Desc:              "Takes fewer risks and takes more time to reconnoiter",
		RiskModifier:      -10,
		ReconModifier:     3,
		mutuallyExclusive: []string{"Reckless"},
	}
}

func CreateReckless() Behaviors {
	return Behaviors{
		Name:              "Reckless",
		Desc:              "Takes more risks and takes less time to reconnoiter",
		RiskModifier:      10,
		ReconModifier:     -1,
		mutuallyExclusive: []string{"Cautious"},
	}
}

var RegularBehaviors = []Behaviors{
	CreateFrugal(), CreateProfligate(), CreateGambler(), CreateCautious(), CreateReckless(),
}

func CreateSquatter() Behaviors {
	return Behaviors{
		Name: "Squatter",
		Desc: "Lives in unoccupied buildings",
	}
}

func CreateLawful() Behaviors {
	return Behaviors{
		Name:              "Lawful",
		Desc:              "Implicitly trusts law enforcement and follows the rules",
		mutuallyExclusive: []string{"Chaotic"},
	}
}

func CreateChaotic() Behaviors {
	return Behaviors{
		Name:              "Chaotic",
		Desc:              "Distrusts law enforcement and breaks the rules",
		mutuallyExclusive: []string{"Lawful"},
	}
}
