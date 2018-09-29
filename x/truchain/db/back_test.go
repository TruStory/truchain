package db

import (
	"fmt"
	"testing"
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
)

func Test_calculateInterest_5days(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(50))
	period := 5 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := calculateInterest(category, amount, period, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(2), "Interest is wrong")
}

func Test_calculateInterest_50days(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(50))
	period := 50 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := calculateInterest(category, amount, period, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(5), "Interest is wrong")
}

func Test_calculateInterest_Max(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(100))
	period := 90 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := calculateInterest(category, amount, period, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(10), "Interest is wrong")
}

func Test_calculateInterest_Min(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(0))
	period := 3 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := calculateInterest(category, amount, period, params)
	fmt.Println(interest.Amount)
	assert.Equal(t, interest.String(), "0trudex", "Interest is wrong")
}
