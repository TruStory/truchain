package db

import (
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// denomPrefix is the prefix for category coins (i.e: trubtc, trustablecoins)
const denomPrefix = "tru"

// NewBacking adds a new vote to the vote store
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

	// mint category coin from trustake
	coins, err := mintCategoryCoins(k, ctx, story.Category, amount, duration, creator)
	if err != nil {
		return -1, err
	}

	// create new backing type
	backing := ts.NewBacking(
		k.newID(ctx, k.backingKey),
		storyID,
		coins,
		time.Now().Add(duration),
		creator)

	// get handle for backing store
	store := ctx.KVStore(k.backingKey)

	// save it in the store
	store.Set(
		generateKey(k.backingKey.String(), backing.ID),
		k.cdc.MustMarshalBinary(backing))

	return backing.ID, nil
}

// mintCategoryCoin mints new category coins by burning trustake and using a formula
// based on the amount of trustake and backing duration
func mintCategoryCoins(
	k TruKeeper,
	ctx sdk.Context,
	cat ts.StoryCategory,
	amount sdk.Coin,
	duration time.Duration,
	addr sdk.AccAddress) (sdk.Coin, sdk.Error) {

	// naive implementation: 1 trustake = 1 category coin
	coin := sdk.NewCoin(cat.CoinDenom(), amount.Amount)

	// burn trustake
	if _, _, err := k.ck.SubtractCoins(ctx, addr, sdk.Coins{amount}); err != nil {
		return sdk.Coin{}, err
	}

	return coin, nil
}
