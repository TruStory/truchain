package backing

import (
	"math"
	"testing"
	"time"

	app "github.com/TruStory/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var fiver = sdk.Coin{
	Amount: sdk.NewInt(5),
	Denom:  app.StakeDenom,
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

	_, err := bk.Backing(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ErrNotFound(id).Code(), err.Code(), "Should get error")
}

func TestToggleVote(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount, _ := sdk.ParseCoin("5trudex")
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 7 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	backingID, _ := bk.Create(ctx, storyID, amount, argument, creator, duration)

	bk.ToggleVote(ctx, backingID)
	b, _ := bk.Backing(ctx, backingID)
	assert.False(t, b.VoteChoice())
}

func TestGetBacking(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew.."
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 7 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	backingID, _ := bk.Create(ctx, storyID, amount, argument, creator, duration)

	b, err := bk.Backing(ctx, backingID)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), b.ID())
}

func TestBackingsByStoryID(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	duration := 7 * time.Hour

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	creator2 := sdk.AccAddress([]byte{2, 3})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})

	bk.Create(ctx, storyID, amount, argument, creator, duration)
	bk.Create(ctx, storyID, amount, argument, creator2, duration)

	backings, _ := bk.BackingsByStoryID(ctx, storyID)
	assert.Equal(t, 2, len(backings))
}

func TestBackingsByStoryIDAndCreator(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	duration := 7 * time.Hour

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	bk.Create(ctx, storyID, amount, argument, creator, duration)

	backing, _ := bk.BackingByStoryIDAndCreator(ctx, storyID, creator)
	assert.Equal(t, int64(1), backing.ID())
}

func TestTally(t *testing.T) {
	ctx, k, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	duration := 7 * time.Hour

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	k.Create(ctx, storyID, amount, argument, creator, duration)

	creator2 := sdk.AccAddress([]byte{2, 3})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})
	k.Create(ctx, storyID, amount, argument, creator2, duration)

	yes, _, _ := k.Tally(ctx, storyID)

	assert.Equal(t, 2, len(yes))
}

func TestTotalBacking(t *testing.T) {
	ctx, k, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	duration := 7 * time.Hour

	creator := sdk.AccAddress([]byte{1, 2})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	k.Create(ctx, storyID, amount, argument, creator, duration)

	creator2 := sdk.AccAddress([]byte{2, 3})
	bankKeeper.AddCoins(ctx, creator2, sdk.Coins{amount})
	k.Create(ctx, storyID, amount, argument, creator2, duration)

	total, _ := k.TotalBackingAmount(ctx, storyID)

	assert.Equal(t, "10000000trusteak", total.String())
}

func TestNewBacking_ErrInsufficientFunds(t *testing.T) {
	ctx, bk, sk, ck, _, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 7 * time.Hour

	_, err := bk.Create(ctx, storyID, amount, argument, creator, duration)
	assert.NotNil(t, err)
	assert.Equal(t, sdk.ErrInsufficientFunds("blah").Code(), err.Code(), "Should get error")
}

func TestNewBacking(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 7 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	backingID, _ := bk.Create(ctx, storyID, amount, argument, creator, duration)
	assert.NotNil(t, backingID)
}

func TestDuplicateBacking(t *testing.T) {
	ctx, bk, sk, ck, bankKeeper, _ := mockDB()
	storyID := createFakeStory(ctx, sk, ck)
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(5000000))
	argument := "cool story brew"
	creator := sdk.AccAddress([]byte{1, 2})
	duration := 7 * time.Hour
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})
	bankKeeper.AddCoins(ctx, creator, sdk.Coins{amount})

	backingID, _ := bk.Create(ctx, storyID, amount, argument, creator, duration)
	assert.NotNil(t, backingID)

	_, err := bk.Create(ctx, storyID, amount, argument, creator, duration)
	assert.Equal(t, ErrDuplicate(storyID, creator).Code(), err.Code())
}

func Test_getInterest_MidAmountMidPeriod(t *testing.T) {
	ctx, k, _, _, _, _ := mockDB()
	// 500,000,000,000,000 nano / 10^9 = 500,000 trudex
	cred := "trudex"
	amount := sdk.NewCoin(cred, sdk.NewInt(500000000000000))
	period := (7 / 2) * 24 * time.Hour
	maxPeriod := 30 * 24 * time.Hour

	interest := k.getInterest(ctx, amount, period, maxPeriod, cred)
	assert.Equal(t, "trudex", interest.Denom)
	assert.Equal(t, sdk.NewInt(11660000000000).String(), interest.Amount.String())
}

func Test_getInterest_MaxAmountMinPeriod(t *testing.T) {
	ctx, k, _, _, _, _ := mockDB()
	cred := "trudex"
	amount := sdk.NewCoin(cred, sdk.NewInt(1000000000000000))
	period := 7 * 24 * time.Hour

	interest := k.getInterest(
		ctx, amount, period, 30*24*time.Hour, cred)
	assert.Equal(t, interest.Amount.String(), sdk.NewInt(48863333333333).String())
}

func Test_getInterest_MinAmountMaxPeriod(t *testing.T) {
	ctx, k, _, _, _, _ := mockDB()
	cred := "trudex"
	amount := sdk.NewCoin(cred, sdk.NewInt(0))
	maxPeriod := 30 * 24 * time.Hour

	interest := k.getInterest(
		ctx, amount, maxPeriod, maxPeriod, cred)
	assert.Equal(t, interest.Amount.String(), sdk.NewInt(0).String())
}

func Test_getInterest_MaxAmountMaxPeriod(t *testing.T) {
	ctx, k, _, _, _, _ := mockDB()

	cred := "trudex"
	amount := sdk.NewCoin("trudex", sdk.NewInt(1000000000000000))
	maxPeriod := 30 * 24 * time.Hour
	maxInterestRate := k.stakeKeeper.GetParams(ctx).MaxInterestRate
	expected := sdk.NewDecFromInt(amount.Amount).Mul(maxInterestRate)

	interest := k.getInterest(
		ctx, amount, maxPeriod, maxPeriod, cred)
	assert.Equal(t, expected.RoundInt().String(), interest.Amount.String())
}
