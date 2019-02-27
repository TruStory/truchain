package stake

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func Test_interest_MidAmountMidPeriod(t *testing.T) {
	ctx, k := mockDB()

	amount := sdk.NewCoin("crypto", sdk.NewInt(500000000000000))
	period := (7 / 2) * 24 * time.Hour
	categoryID := int64(1)

	interest := k.interest(ctx, amount, categoryID, period)
	assert.Equal(t, sdk.NewInt(11660000000000).String(), interest.String())
}

func Test_interest_MaxAmountMinPeriod(t *testing.T) {
	ctx, k := mockDB()
	amount := sdk.NewCoin("crypto", sdk.NewInt(1000000000000000))
	period := 7 * 24 * time.Hour
	categoryID := int64(1)

	interest := k.interest(ctx, amount, categoryID, period)
	assert.Equal(t, interest.String(), sdk.NewInt(48863333333333).String())
}

func Test_interest_MinAmountMaxPeriod(t *testing.T) {
	ctx, k := mockDB()
	amount := sdk.NewCoin("crypto", sdk.NewInt(0))
	period := 30 * 24 * time.Hour
	categoryID := int64(1)

	interest := k.interest(ctx, amount, categoryID, period)
	assert.Equal(t, interest.String(), sdk.NewInt(0).String())
}

func Test_interest_MaxAmountMaxPeriod(t *testing.T) {
	ctx, k := mockDB()
	amount := sdk.NewCoin("crypto", sdk.NewInt(1000000000000000))
	period := 30 * 24 * time.Hour
	categoryID := int64(1)
	maxInterestRate := k.GetParams(ctx).MaxInterestRate
	expected := sdk.NewDecFromInt(amount.Amount).Mul(maxInterestRate)

	interest := k.interest(ctx, amount, categoryID, period)
	assert.Equal(t, expected.RoundInt().String(), interest.String())
}
