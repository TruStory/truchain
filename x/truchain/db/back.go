package db

import (
	"fmt"
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewBacking adds a new backing to the backing store
func (k TruKeeper) NewBacking(
	ctx sdk.Context,
	storyID int64,
	amount sdk.Coin,
	creator sdk.AccAddress,
	duration time.Duration,
) (int64, sdk.Error) {

	// get story from story id
	story, err := k.GetStory(ctx, storyID)
	if err != nil {
		return -1, err
	}

	// naïve implementaion: 1 trustake = 1 category coin
	// https://github.com/TruStory/truchain/issues/21
	conversionRate := sdk.NewInt(int64(1))

	// mint category coin from trustake -- the principal, amount user gets back
	principal, err := convertCoins(k, ctx, story.Category, amount, duration, creator, conversionRate)
	if err != nil {
		return -1, err
	}

	// load default backing parameters
	params := ts.NewBackingParams()

	// mint category coin from interest earned
	interest := calculateInterest(story.Category, amount, duration, params)

	// create new backing type
	backing := ts.NewBacking(
		k.newID(ctx, k.backingKey),
		storyID,
		principal,
		interest,
		time.Now().Add(duration),
		params,
		duration,
		creator)

	// get handle for backing store
	store := ctx.KVStore(k.backingKey)

	// save backing in the store
	store.Set(
		generateKey(k.backingKey.String(), backing.ID),
		k.cdc.MustMarshalBinary(backing))

	// add backing to the backing queue for processing
	k.BackingQueuePush(ctx, backing.ID)

	return backing.ID, nil
}

// GetBacking gets the backing at the current index from the KVStore
func (k TruKeeper) GetBacking(ctx sdk.Context, id int64) (ts.Backing, sdk.Error) {
	store := ctx.KVStore(k.backingKey)
	key := generateKey(k.backingKey.String(), id)
	val := store.Get(key)
	if val == nil {
		return ts.Backing{}, ts.ErrVoteNotFound(id)
	}
	backing := &ts.Backing{}
	k.cdc.MustUnmarshalBinary(val, backing)

	return *backing, nil
}

// convertCoins mints new category coins by burning trustake
func convertCoins(
	k TruKeeper,
	ctx sdk.Context,
	cat ts.StoryCategory,
	amount sdk.Coin,
	duration time.Duration,
	addr sdk.AccAddress,
	conversionRate sdk.Int) (sdk.Coin, sdk.Error) {

	// mint new category coins
	coin := sdk.NewCoin(cat.CoinDenom(), amount.Amount.Mul(conversionRate))

	// burn trustake
	if _, _, err := k.ck.SubtractCoins(ctx, addr, sdk.Coins{amount}); err != nil {
		return sdk.Coin{}, err
	}

	return coin, nil
}

// calculateInterest calcuates the interest for the backing
func calculateInterest(
	category ts.StoryCategory,
	amount sdk.Coin,
	period time.Duration,
	params ts.BackingParams) sdk.Coin {

	// TODO: keep track of total supply
	// https://github.com/TruStory/truchain/issues/22
	coinBalance := int64(100)

	fmt.Println("testing..")

	maxAmount := coinBalance
	maxPeriod := int64(365 * 24 * time.Hour)
	normalizedAmount := amount.Amount.Div(sdk.NewInt(maxAmount))
	normalizedPeriod := sdk.NewInt(int64(period*time.Hour) / maxPeriod)
	amountWeight := normalizedAmount.Mul(sdk.NewInt(int64(1 / 3)))
	periodWeight := normalizedPeriod.Mul(sdk.NewInt(int64(2 / 3)))

	maxInterestRate := int64(10)
	minInterestRate := int64(0)

	interestRateRange := sdk.NewInt(maxInterestRate - minInterestRate)
	baseInterestRate := interestRateRange.Mul(amountWeight.Add(periodWeight))
	interest := baseInterestRate.Add(sdk.NewInt(minInterestRate))

	return sdk.NewCoin(category.CoinDenom(), interest)
}
