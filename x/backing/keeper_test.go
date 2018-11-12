package backing

import (
	"math"
	"testing"
	"time"

	params "github.com/TruStory/truchain/parameters"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var fiver = sdk.Coin{
	Amount: sdk.NewInt(5),
	Denom:  params.StakeDenom,
}

func Test_key(t *testing.T) {
	_, bk, _, _, _, _ := mockDB()

	bz1 := bk.GetIDKey(5)
	bz2 := bk.GetIDKey(math.MaxInt64)

	assert.Equal(t, "backings:id:5", string(bz1), "should generate valid key")
	assert.Equal(t, "backings:id:9223372036854775807", string(bz2), "should generate valid key")
}

func TestGetBacking_ErrBackingNotFound(t *testing.T) {
	ctx, bk, _, _, _, _ := mockDB()
	id := int64(5)

	_, err := bk.GetBacking(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound(id).Code(), err.Code(), "Should get error")
}

func TestGetBacking(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := DefaultMsgParams().MinPeriod
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	backingID, _ := bk.NewBacking(ctx, storyID, amount, creator, duration)

	_, err := bk.GetBacking(ctx, backingID)
	assert.Nil(t, err)
}

func TestNewBacking_ErrInsufficientFunds(t *testing.T) {
	ctx, bk, sk, ck, _, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := DefaultMsgParams().MinPeriod

	_, err := bk.NewBacking(ctx, storyID, amount, creator, duration)
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInsufficientFunds("blah").Code(), err.Code(), "Should get error")
}

func TestNewBacking(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	creator := sdk.AccAddress([]byte{1, 2})
	duration := DefaultMsgParams().MinPeriod
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	backingID, _ := bk.NewBacking(ctx, storyID, amount, creator, duration)
	assert.NotNil(t, backingID)
}

func Test_getPrincipal_InCategoryCoins(t *testing.T) {
	ctx, bk, _, ck, bankKeeper, _ := mockDB()
	cat := createFakeCategory(ctx, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	userAddr := sdk.AccAddress([]byte{1, 2})

	// give fake user some fake category coins
	bankKeeper.AddCoins(ctx, userAddr, sdk.Coins{amount})

	coin, err := bk.getPrincipal(ctx, cat, amount, userAddr)
	assert.Nil(t, err)
	assert.Equal(t, amount, coin, "Incorrect principal calculation")
	assert.Equal(t, "trudex", amount.Denom, "Incorrect principal coin")
}

func Test_getPrincipal_InTrustake(t *testing.T) {
	ctx, bk, _, ck, bankKeeper, _ := mockDB()
	cat := createFakeCategory(ctx, ck)
	userAddr := sdk.AccAddress([]byte{1, 2})

	// give fake user some fake trustake
	bankKeeper.AddCoins(ctx, userAddr, sdk.Coins{fiver})

	// back with trustake, get principal in cat coins
	coin, err := bk.getPrincipal(ctx, cat, fiver, userAddr)
	assert.Nil(t, err)
	assert.Equal(t, fiver.Amount, coin.Amount, "Incorrect principal calculation")
	assert.Equal(t, "trudex", coin.Denom, "Incorrect principal coin")
}

func Test_getPrincipal_ErrInvalidCoin(t *testing.T) {
	ctx, bk, _, ck, bankKeeper, _ := mockDB()
	cat := createFakeCategory(ctx, ck)
	amount := sdk.NewCoin("trubtc", sdk.NewInt(5))
	userAddr := sdk.AccAddress([]byte{1, 2})

	// give fake user some fake coins
	bankKeeper.AddCoins(ctx, userAddr, sdk.Coins{fiver})

	_, err := bk.getPrincipal(ctx, cat, amount, userAddr)
	assert.NotNil(t, err)
	// assert.Equal(t, sdk.ErrInsufficientCoins().Code(), err.Code(), "invalid error")
}

func Test_getInterest_MidAmountMidPeriod(t *testing.T) {
	ctx, _, _, ck, _, _ := mockDB()
	cat := createFakeCategory(ctx, ck)
	// 500,000,000,000,000 nano / 10^9 = 500,000 trudex
	amount := sdk.NewCoin("trudex", sdk.NewInt(500000000000000))
	period := 45 * 24 * time.Hour
	params := DefaultParams()
	maxPeriod := DefaultMsgParams().MaxPeriod

	interest := getInterest(cat, amount, period, maxPeriod, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(25000000000000), "Interest is wrong")
}

func Test_getInterest_MaxAmountMinPeriod(t *testing.T) {
	ctx, _, _, ck, _, _ := mockDB()
	cat := createFakeCategory(ctx, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(1000000000000000))
	period := 3 * 24 * time.Hour
	params := DefaultParams()

	interest := getInterest(
		cat, amount, period, DefaultMsgParams().MaxPeriod, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(35523333300000), "Interest is wrong")
}

func Test_getInterest_MinAmountMaxPeriod(t *testing.T) {
	ctx, _, _, ck, _, _ := mockDB()
	cat := createFakeCategory(ctx, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(0))
	period := 90 * 24 * time.Hour
	params := DefaultParams()

	interest := getInterest(
		cat, amount, period, DefaultMsgParams().MaxPeriod, params)
	assert.Equal(t, interest.Amount, sdk.NewInt(0), "Interest is wrong")
}

func Test_getInterest_MaxAmountMaxPeriod(t *testing.T) {
	ctx, _, _, ck, _, _ := mockDB()
	cat := createFakeCategory(ctx, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(1000000000000000))
	period := 90 * 24 * time.Hour
	params := DefaultParams()
	expected := sdk.NewDecFromInt(amount.Amount).Mul(params.MaxInterestRate)

	interest := getInterest(
		cat, amount, period, DefaultMsgParams().MaxPeriod, params)
	assert.Equal(t, expected.RoundInt(), interest.Amount, "Interest is wrong")
}

func Test_getInterest_MinAmountMinPeriod(t *testing.T) {
	ctx, _, _, ck, _, _ := mockDB()
	cat := createFakeCategory(ctx, ck)
	amount := sdk.NewCoin("trudex", sdk.NewInt(0))
	period := 3 * 24 * time.Hour
	params := DefaultParams()

	interest := getInterest(
		cat, amount, period, DefaultMsgParams().MaxPeriod, params)
	assert.Equal(t, interest.String(), "0trudex", "Interest is wrong")
}
