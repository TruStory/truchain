package db

import (
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

	// Check if user has enough cat coins or trustake to back
	trustake := sdk.NewCoin(ts.NativeTokenName, amount.Amount)
	if !k.ck.HasCoins(ctx, creator, sdk.Coins{amount}) &&
		!k.ck.HasCoins(ctx, creator, sdk.Coins{trustake}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for backing.")
	}

	// get story from story id
	story, err := k.GetStory(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// load default backing parameters
	params := ts.NewBackingParams()

	// set principal, converting from trustake if needed
	principal, err := k.getPrincipal(ctx, story.Category, amount, creator)
	if err != nil {
		return 0, err
	}

	// mint category coin from interest earned
	interest := getInterest(story.Category, amount, duration, params)

	// create new backing type
	backing := ts.NewBacking(
		k.id(ctx, k.backingKey),
		storyID,
		principal,
		interest,
		time.Now().Add(duration),
		params,
		duration,
		creator)

	// store backing
	k.setBacking(ctx, backing)

	// add backing to the backing queue for processing
	k.BackingQueuePush(ctx, backing.ID)

	return backing.ID, nil
}

// GetBacking gets the backing at the current index from the KVStore
func (k TruKeeper) GetBacking(ctx sdk.Context, id int64) (ts.Backing, sdk.Error) {
	store := ctx.KVStore(k.backingKey)
	key := key(k.backingKey.String(), id)
	val := store.Get(key)
	if val == nil {
		return ts.Backing{}, ts.ErrBackingNotFound(id)
	}
	backing := &ts.Backing{}
	k.cdc.MustUnmarshalBinary(val, backing)

	return *backing, nil
}

// ============================================================================

// getPrincipal calculates the principal, the amount the user gets back
// after the backing expires/matures. Returns a coin.
func (k TruKeeper) getPrincipal(
	ctx sdk.Context,
	cat ts.StoryCategory,
	amount sdk.Coin,
	userAddr sdk.AccAddress) (sdk.Coin, sdk.Error) {

	// check and return amount if user has enough category coins
	if k.ck.HasCoins(ctx, userAddr, sdk.Coins{amount}) {
		return amount, nil
	}

	// naïve implementaion: 1 trustake = 1 category coin
	// https://github.com/TruStory/truchain/issues/21
	conversionRate := sdk.NewDec(1)

	// mint new category coins
	principal := sdk.NewCoin(
		cat.CoinDenom(),
		sdk.NewDecFromInt(amount.Amount).Mul(conversionRate).RoundInt())

	// burn equivalent trustake
	trustake := sdk.Coins{sdk.NewCoin(ts.NativeTokenName, principal.Amount)}
	if _, _, err := k.ck.SubtractCoins(ctx, userAddr, trustake); err != nil {
		return sdk.Coin{}, err
	}

	return principal, nil
}

// getInterest calcuates the interest for the backing
func getInterest(
	category ts.StoryCategory,
	amount sdk.Coin,
	period time.Duration,
	params ts.BackingParams) sdk.Coin {

	// TODO: keep track of total supply
	// https://github.com/TruStory/truchain/issues/22

	// 1,000,000,000 preethi = 10^9 = 1 trustake
	// 1,000,000,000,000,000 preethi = 10^15 / 10^9 = 10^6 = 1,000,000 trustake
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

	// return coin with rounded interest
	coin := sdk.NewCoin(category.CoinDenom(), interest.RoundInt())

	return coin
}

// setBacking stores a `Backing` type in the KVStore
func (k TruKeeper) setBacking(ctx sdk.Context, backing ts.Backing) {
	store := ctx.KVStore(k.backingKey)
	store.Set(
		key(k.backingKey.String(), backing.ID),
		k.cdc.MustMarshalBinary(backing))
}
