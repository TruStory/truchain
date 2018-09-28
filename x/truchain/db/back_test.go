package db

import (
	"testing"
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
)

func Test_calculateInterest(t *testing.T) {
	// ctx, _, _, _ := mockDB()

	category := ts.DEX
	amount := sdk.NewCoin("cat", sdk.NewInt(5))
	period := 5 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := calculateInterest(category, amount, period, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(5), "Interest is wrong")
}
