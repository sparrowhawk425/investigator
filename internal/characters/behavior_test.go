package characters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getAmt(amt int) int {
	return amt
}

func TestGetLootAmount(t *testing.T) {

	// No mod should return 50%
	behavior := Behavior{LootAmountModifier: 0}
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
