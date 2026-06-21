package characters

import (
	"slices"
	"testing"

	"github.com/sparrowhawk425/investigators/internal/gameobjects"
	"github.com/stretchr/testify/assert"
)

func getAmt(amt int) int {
	return amt
}

func TestGetLootAmount(t *testing.T) {

	// No mod should return 50%
	behavior := Behaviors{LootAmountModifier: 0}
	finalAmt := behavior.GetLootAmount(getAmt)(100)
	assert.Equal(t, 50, finalAmt)

	// Mod 20 should return 70%
	behavior.LootAmountModifier = 20
	finalAmt = behavior.GetLootAmount(getAmt)(100)
	assert.Equal(t, 70, finalAmt)

	// Mod -20 should return 30%
	behavior.LootAmountModifier = -20
	finalAmt = behavior.GetLootAmount(getAmt)(100)
	assert.Equal(t, 30, finalAmt)

	// No mod for 1 amt
	behavior.LootAmountModifier = 0
	finalAmt = behavior.GetLootAmount(getAmt)(1)
	assert.Equal(t, 1, finalAmt)
}

// TODO: get the most common from the whole list (return any of ties)?
func testFindTarget(locations []gameobjects.Location) *gameobjects.Location {
	places := []gameobjects.Location{}

	for _, loc := range locations {
		if slices.ContainsFunc(places, func(l gameobjects.Location) bool { return l.Address.Name == loc.Address.Name }) {
			return &loc
		}
		places = append(places, loc)
	}
	return nil
}

func TestFindTargetQuality(t *testing.T) {

	expensive := "Expensive Casino"
	moderate := "Moderate Casino"
	cheap := "Cheap Casino"

	locations := []gameobjects.Location{
		{Type: gameobjects.Casino, Address: gameobjects.Address{Number: 100, Name: expensive}, Quality: gameobjects.Expensive},
		{Type: gameobjects.Casino, Address: gameobjects.Address{Number: 200, Name: moderate}, Quality: gameobjects.Moderate},
		{Type: gameobjects.Casino, Address: gameobjects.Address{Number: 300, Name: cheap}, Quality: gameobjects.Cheap},
	}

	// Test Expensive
	behavior := Behaviors{QualityPreference: []gameobjects.Quality{gameobjects.Expensive}}
	target := behavior.FindTarget(testFindTarget)(locations)
	assert.NotNil(t, target)
	assert.Equal(t, expensive, target.Address.Name)

	// Test Moderate
	behavior = Behaviors{QualityPreference: []gameobjects.Quality{gameobjects.Moderate}}
	target = behavior.FindTarget(testFindTarget)(locations)
	assert.NotNil(t, target)
	assert.Equal(t, moderate, target.Address.Name)

	// Test Cheap
	behavior = Behaviors{QualityPreference: []gameobjects.Quality{gameobjects.Cheap}}
	target = behavior.FindTarget(testFindTarget)(locations)
	assert.NotNil(t, target)
	assert.Equal(t, cheap, target.Address.Name)
}

func TestFindTargetLocationType(t *testing.T) {

	residence := "Residence"
	store := "Store"
	casino := "Casino"

	locations := []gameobjects.Location{
		{Type: gameobjects.Residence, Address: gameobjects.Address{Number: 100, Name: residence}, Quality: gameobjects.Expensive},
		{Type: gameobjects.Store, Address: gameobjects.Address{Number: 200, Name: store}, Quality: gameobjects.Moderate},
		{Type: gameobjects.Casino, Address: gameobjects.Address{Number: 300, Name: casino}, Quality: gameobjects.Cheap},
	}

	// Test Residence
	behavior := Behaviors{LocationPreference: []gameobjects.LocationType{gameobjects.Residence}}
	target := behavior.FindTarget(testFindTarget)(locations)
	assert.NotNil(t, target)
	assert.Equal(t, residence, target.Address.Name)

	// Test Store
	behavior = Behaviors{LocationPreference: []gameobjects.LocationType{gameobjects.Store}}
	target = behavior.FindTarget(testFindTarget)(locations)
	assert.NotNil(t, target)
	assert.Equal(t, store, target.Address.Name)

	// Test Casino
	behavior = Behaviors{LocationPreference: []gameobjects.LocationType{gameobjects.Casino}}
	target = behavior.FindTarget(testFindTarget)(locations)
	assert.NotNil(t, target)
	assert.Equal(t, casino, target.Address.Name)

	// Test no pref
	behavior = Behaviors{LocationPreference: []gameobjects.LocationType{}}
	target = behavior.FindTarget(testFindTarget)(locations)
	assert.Nil(t, target)
}

func TestFindTargetQualityAndLocationType(t *testing.T) {

}
