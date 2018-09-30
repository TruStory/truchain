package db

import (
	"testing"
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
)

func Test_calculateInterest_MidAmountMidPeriod(t *testing.T) {
	category := ts.DEX
	// 500,000,000,000,000 nano / 10^9 = 500,000 trudex
	amount := sdk.NewCoin("trudex", sdk.NewInt(500000000000000))
	period := 45 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := calculateInterest(category, amount, period, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(25000000000000), "Interest is wrong")
}

func Test_calculateInterest_MaxAmountMinPeriod(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(1000000000000000))
	period := 3 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := calculateInterest(category, amount, period, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(35523333300000), "Interest is wrong")
}

func Test_calculateInterest_MinAmountMaxPeriod(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(0))
	period := 90 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := calculateInterest(category, amount, period, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(0), "Interest is wrong")
}

func Test_calculateInterest_MaxAmountMaxPeriod(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(1000000000000000))
	period := 90 * 24 * time.Hour
	params := ts.NewBackingParams()
	expected := sdk.NewDecFromInt(amount.Amount).Mul(params.MaxInterestRate)

	interest := calculateInterest(category, amount, period, params)
	assert.Equal(t, expected.RoundInt(), interest.Amount, "Interest is wrong")
}

func Test_calculateInterest_MinAmountMinPeriod(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(0))
	period := 3 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := calculateInterest(category, amount, period, params)
	assert.Equal(t, interest.String(), "0trudex", "Interest is wrong")
}
