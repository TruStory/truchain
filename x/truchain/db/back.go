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
		// TODO: change to backing error
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
	// TODO: handle precision, write test
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

	// 1,000,000,000 nano = 10^9 = 1 trustake
	// 1,000,000,000,000,000 nano = 10^15 / 10^9 = 10^6 = 1,000,000 trustake
	totalSupply := sdk.NewDec(1000000000000000)

	// inputs
	maxAmount := totalSupply
	maxPeriod := params.MaxPeriod
	amountWeight := params.AmountWeight
	periodWeight := params.PeriodWeight
	maxInterestRate := params.MaxInterestRate

	// type cast values to unitless decimals for math operations
	periodDec := sdk.NewDec(int64(period))
	maxPeriodDec := sdk.NewDec(int64(maxPeriod))
	amountDec := sdk.NewDecFromInt(amount.Amount)

	// normalize amount and period to 0 - 1
	normalizedAmount := amountDec.Quo(maxAmount)
	normalizedPeriod := periodDec.Quo(maxPeriodDec)

	// apply weights to normalized amount and period
	weightedAmount := normalizedAmount.Mul(amountWeight)
	weightedPeriod := normalizedPeriod.Mul(periodWeight)

	// calculate interest
	interestRate := maxInterestRate.Mul(weightedAmount.Add(weightedPeriod))
	// convert rate to a value
	interest := amountDec.Mul(interestRate)

	// debugging...
	fmt.Println(normalizedAmount.String())
	fmt.Println(normalizedPeriod.String())
	fmt.Println(interest.String())
	fmt.Println(sdk.NewCoin(category.CoinDenom(), interest.RoundInt()))

	// output: coin with rounded interest
	coin := sdk.NewCoin(category.CoinDenom(), interest.RoundInt())

	return coin
}
