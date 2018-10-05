package db

import (
	"testing"
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
)

func TestGetBacking_ErrBackingNotFound(t *testing.T) {
	ctx, _, _, k := MockDB()
	id := int64(5)

	_, err := k.GetBacking(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ts.ErrBackingNotFound(id).Code(), err.Code(), "Should get error")
}

func TestGetBacking(t *testing.T) {
	ctx, ms, _, k := MockDB()
	storyID := CreateFakeStory(ms, k)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := ts.NewBackingParams().MinPeriod
	k.ck.AddCoins(ctx, creator, sdk.Coins{amount})
	backingID, _ := k.NewBacking(ctx, storyID, amount, creator, duration)

	_, err := k.GetBacking(ctx, backingID)
	assert.Nil(t, err)
}

func TestNewBacking_ErrInsufficientFunds(t *testing.T) {
	ctx, ms, _, k := MockDB()
	storyID := CreateFakeStory(ms, k)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := ts.NewBackingParams().MinPeriod

	_, err := k.NewBacking(ctx, storyID, amount, creator, duration)
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInsufficientFunds("blah").Code(), err.Code(), "Should get error")
}

func TestNewBacking(t *testing.T) {
	ctx, ms, _, k := MockDB()
	storyID := CreateFakeStory(ms, k)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := ts.NewBackingParams().MinPeriod
	k.ck.AddCoins(ctx, creator, sdk.Coins{amount})

	backingID, _ := k.NewBacking(ctx, storyID, amount, creator, duration)
	assert.NotNil(t, backingID)
}

func Test_getPrincipal_InCategoryCoins(t *testing.T) {
	ctx, _, _, k := MockDB()
	cat := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(5))
	userAddr := sdk.AccAddress([]byte{1, 2})

	// give fake user some fake category coins
	k.ck.AddCoins(ctx, userAddr, sdk.Coins{amount})

	coin, err := k.getPrincipal(ctx, cat, amount, userAddr)
	assert.Nil(t, err)
	assert.Equal(t, amount, coin, "Incorrect principal calculation")
	assert.Equal(t, "trudex", amount.Denom, "Incorrect principal coin")
}

func Test_getPrincipal_InTrustake(t *testing.T) {
	ctx, _, _, k := MockDB()
	cat := ts.DEX
	userAddr := sdk.AccAddress([]byte{1, 2})

	// give fake user some fake trustake
	amount, _ := sdk.ParseCoin("5trustake")
	k.ck.AddCoins(ctx, userAddr, sdk.Coins{amount})

	// back with trustake, get principal in cat coins
	coin, err := k.getPrincipal(ctx, cat, amount, userAddr)
	assert.Nil(t, err)
	assert.Equal(t, amount.Amount, coin.Amount, "Incorrect principal calculation")
	assert.Equal(t, "trudex", coin.Denom, "Incorrect principal coin")
}

func Test_getPrincipal_ErrInvalidCoin(t *testing.T) {
	ctx, _, _, k := MockDB()
	cat := ts.DEX
	amount := sdk.NewCoin("trubtc", sdk.NewInt(5))
	userAddr := sdk.AccAddress([]byte{1, 2})

	// give fake user some fake coins
	fakeCoin, _ := sdk.ParseCoin("5trustake")
	k.ck.AddCoins(ctx, userAddr, sdk.Coins{fakeCoin})

	_, err := k.getPrincipal(ctx, cat, amount, userAddr)
	assert.NotNil(t, err)
	assert.Equal(t, ts.ErrInvalidBackingCoin().Code(), err.Code(), "invalid error")
}

func Test_getInterest_MidAmountMidPeriod(t *testing.T) {
	category := ts.DEX
	// 500,000,000,000,000 nano / 10^9 = 500,000 trudex
	amount := sdk.NewCoin("trudex", sdk.NewInt(500000000000000))
	period := 45 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := getInterest(category, amount, period, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(25000000000000), "Interest is wrong")
}

func Test_getInterest_MaxAmountMinPeriod(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(1000000000000000))
	period := 3 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := getInterest(category, amount, period, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(35523333300000), "Interest is wrong")
}

func Test_getInterest_MinAmountMaxPeriod(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(0))
	period := 90 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := getInterest(category, amount, period, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(0), "Interest is wrong")
}

func Test_getInterest_MaxAmountMaxPeriod(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(1000000000000000))
	period := 90 * 24 * time.Hour
	params := ts.NewBackingParams()
	expected := sdk.NewDecFromInt(amount.Amount).Mul(params.MaxInterestRate)

	interest := getInterest(category, amount, period, params)
	assert.Equal(t, expected.RoundInt(), interest.Amount, "Interest is wrong")
}

func Test_getInterest_MinAmountMinPeriod(t *testing.T) {
	category := ts.DEX
	amount := sdk.NewCoin("trudex", sdk.NewInt(0))
	period := 3 * 24 * time.Hour
	params := ts.NewBackingParams()

	interest := getInterest(category, amount, period, params)
	assert.Equal(t, interest.String(), "0trudex", "Interest is wrong")
}
