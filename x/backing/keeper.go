package backing

import (
	"time"

	params "github.com/TruStory/truchain/parameters"
	app "github.com/TruStory/truchain/types"
	cat "github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
)

// ReadKeeper defines a module interface that facilitates read only access
type ReadKeeper interface {
	app.ReadKeeper

	Backing(ctx sdk.Context, id int64) (backing Backing, err sdk.Error)
	Tally(ctx sdk.Context, storyID int64) (yes []Backing, no []Backing, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context,
		storyID int64,
		amount sdk.Coin,
		creator sdk.AccAddress,
		duration time.Duration,
	) (int64, sdk.Error)

	NewResponseEndBlock(ctx sdk.Context) sdk.Tags
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	storyKeeper    story.WriteKeeper // read-write access to story store
	bankKeeper     bank.Keeper       // read-write access coin store
	categoryKeeper cat.ReadKeeper    // read access to category store

	backingsList app.UserList // backings <-> story mappings
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	storyKeeper story.WriteKeeper,
	bankKeeper bank.Keeper,
	categoryKeeper cat.ReadKeeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		storyKeeper,
		bankKeeper,
		categoryKeeper,
		app.NewUserList(storyKeeper.GetStoreKey()),
	}
}

// ============================================================================

// Create adds a new backing to the backing store
func (k Keeper) Create(
	ctx sdk.Context,
	storyID int64,
	amount sdk.Coin,
	creator sdk.AccAddress,
	duration time.Duration,
) (id int64, err sdk.Error) {

	// Check if user has enough cat coins or trustake to back
	trustake := sdk.NewCoin(params.StakeDenom, amount.Amount)
	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) &&
		!k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{trustake}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for backing.")
	}

	// get story value from story id
	story, err := k.storyKeeper.GetStory(ctx, storyID)
	if err != nil {
		return
	}

	// get category value from category id
	cat, err := k.categoryKeeper.GetCategory(ctx, story.CategoryID)
	if err != nil {
		return
	}

	// load default backing parameters
	params := DefaultParams()

	// set principal, converting from trustake if needed
	principal, err := k.getPrincipal(ctx, cat, amount, creator)
	if err != nil {
		return
	}

	// mint category coin from interest earned
	interest := getInterest(
		cat, amount, duration, DefaultMsgParams().MaxPeriod, params)

	// create new implicit true vote type
	vote := app.NewVote(
		k.GetNextID(ctx), principal, creator, true, app.NewTimestamp(ctx.BlockHeader()))

	// create new backing type with embedded vote
	backing := Backing{
		Vote:     vote,
		StoryID:  storyID,
		Interest: interest,
		Expires:  time.Now().Add(duration),
		Params:   params,
		Period:   duration,
	}

	// store backing
	k.setBacking(ctx, backing)

	// add backing to the backing queue for processing
	k.QueuePush(ctx, backing.ID)

	// add backing <-> story mapping
	k.backingsList.Append(ctx, k, storyID, creator, backing.ID)

	return backing.ID, nil
}

// Backing gets the backing at the current index from the KVStore
func (k Keeper) Backing(ctx sdk.Context, id int64) (backing Backing, err sdk.Error) {
	store := k.GetStore(ctx)
	key := k.GetIDKey(id)
	val := store.Get(key)
	if val == nil {
		return backing, ErrNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinary(val, &backing)

	return
}

// BackingsByStory returns backings for a given story id
func (k Keeper) BackingsByStory(
	ctx sdk.Context, storyID int64) (backings []Backing, err sdk.Error) {

	// iterate over backing list and get backings
	err = k.backingsList.Map(ctx, k, storyID, func(backingID int64) sdk.Error {
		backing, err := k.Backing(ctx, backingID)
		if err != nil {
			return err
		}
		backings = append(backings, backing)

		return nil
	})

	if err != nil {
		return backings, err
	}

	return backings, nil
}

// Tally backings for voting
func (k Keeper) Tally(
	ctx sdk.Context, storyID int64) (yes []Backing, no []Backing, err sdk.Error) {

	err = k.backingsList.Map(ctx, k, storyID, func(backingID int64) sdk.Error {
		backing, err := k.Backing(ctx, backingID)
		if err != nil {
			return err
		}

		if backing.Vote.Vote == true {
			yes = append(yes, backing)
		} else {
			no = append(no, backing)
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}

// ============================================================================

// getPrincipal calculates the principal, the amount the user gets back
// after the backing expires/matures. Returns a coin.
func (k Keeper) getPrincipal(
	ctx sdk.Context,
	cat cat.Category,
	amount sdk.Coin,
	userAddr sdk.AccAddress) (principal sdk.Coin, err sdk.Error) {

	// check which type of coin user wants to back in
	switch amount.Denom {
	case cat.CoinName():
		// check and return amount if user has enough category coins
		if k.bankKeeper.HasCoins(ctx, userAddr, sdk.Coins{amount}) {
			return amount, nil
		}
	case params.StakeDenom:
		// mint category coins from trustake
		return mintFromNativeToken(ctx, k.bankKeeper, cat, amount, userAddr)
	default:
		return principal, sdk.ErrInvalidCoins("Invalid backing token")

	}

	return
}

// setBacking stores a `Backing` type in the KVStore
func (k Keeper) setBacking(ctx sdk.Context, backing Backing) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(backing.ID),
		k.GetCodec().MustMarshalBinary(backing))
}

// ============================================================================

// mintFromNativeToken creates category coins by burning trustake
func mintFromNativeToken(
	ctx sdk.Context,
	bankKeeper bank.Keeper,
	cat cat.Category,
	amount sdk.Coin,
	userAddr sdk.AccAddress) (principal sdk.Coin, err sdk.Error) {

	// naïve implementation: 1 trustake = 1 category coin
	// TODO [Shane]: https://github.com/TruStory/truchain/issues/21
	conversionRate := sdk.NewDec(1)

	// mint new category coins
	principal = sdk.NewCoin(
		cat.CoinName(),
		sdk.NewDecFromInt(amount.Amount).Mul(conversionRate).RoundInt())

	// burn equivalent trustake
	trustake := sdk.Coins{sdk.NewCoin(params.StakeDenom, principal.Amount)}
	_, _, err = bankKeeper.SubtractCoins(ctx, userAddr, trustake)
	if err != nil {
		return
	}

	return
}

// getInterest calcuates the interest for the backing
func getInterest(
	category cat.Category,
	amount sdk.Coin,
	period time.Duration,
	maxPeriod time.Duration,
	params Params) sdk.Coin {

	// TODO: keep track of total supply
	// https://github.com/TruStory/truchain/issues/22

	// 1,000,000,000 preethi = 10^9 = 1 trustake
	// 1,000,000,000,000,000 preethi = 10^15 / 10^9 = 10^6 = 1,000,000 trustake
	totalSupply := sdk.NewDec(1000000000000000)

	// inputs
	maxAmount := totalSupply
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
	coin := sdk.NewCoin(category.CoinName(), interest.RoundInt())

	return coin
}
