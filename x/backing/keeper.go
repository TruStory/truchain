package backing

import (
	"time"

	app "github.com/TruStory/truchain/types"
	cat "github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

// ReadKeeper defines a module interface that facilitates read only access
// to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	GetBacking(ctx sdk.Context, id int64) (backing Backing, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access
// to truchain data
type WriteKeeper interface {
	NewBacking(
		ctx sdk.Context,
		storyID int64,
		amount sdk.Coin,
		creator sdk.AccAddress,
		duration time.Duration,
	) (int64, sdk.Error)
	NewResponseEndBlock(ctx sdk.Context) abci.ResponseEndBlock
}

// ReadWriteKeeper defines a module interface that facilities read/write access
// to truchain data
type ReadWriteKeeper interface {
	ReadKeeper
	WriteKeeper
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	backingKey     sdk.StoreKey          // key to backing store
	baseKeeper     app.Keeper            // base keeper
	storyKeeper    story.ReadWriteKeeper // read-write access to story store
	bankKeeper     bank.Keeper           // read-write access coin store
	categoryKeeper cat.ReadKeeper        // read access to category store
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	backingKey sdk.StoreKey,
	sk story.ReadWriteKeeper,
	bk bank.Keeper,
	ck cat.ReadKeeper,
	codec *amino.Codec) Keeper {
	return Keeper{
		backingKey:     backingKey,
		baseKeeper:     app.NewKeeper(codec),
		storyKeeper:    sk,
		bankKeeper:     bk,
		categoryKeeper: ck,
	}
}

// ============================================================================

// GetCodec returns the base keeper's underlying codec
func (k Keeper) GetCodec() *amino.Codec {
	return k.baseKeeper.Codec
}

// NewBacking adds a new backing to the backing store
func (k Keeper) NewBacking(
	ctx sdk.Context,
	storyID int64,
	amount sdk.Coin,
	creator sdk.AccAddress,
	duration time.Duration,
) (int64, sdk.Error) {

	// Check if user has enough cat coins or trustake to back
	trustake := sdk.NewCoin(NativeTokenName, amount.Amount)
	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) &&
		!k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{trustake}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for backing.")
	}

	// get story value from story id
	story, err := k.storyKeeper.GetStory(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// get category value from category id
	cat, err := k.categoryKeeper.GetCategory(ctx, story.CategoryID)
	if err != nil {
		return 0, err
	}

	// load default backing parameters
	params := NewParams()

	// set principal, converting from trustake if needed
	principal, err := k.getPrincipal(ctx, cat, amount, creator)
	if err != nil {
		return 0, err
	}

	// mint category coin from interest earned
	interest := getInterest(cat, amount, duration, params)

	// create new backing type
	backing := NewBacking(
		k.baseKeeper.GetNextID(ctx, k.backingKey),
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
	k.QueuePush(ctx, backing.ID)

	return backing.ID, nil
}

// GetBacking gets the backing at the current index from the KVStore
func (k Keeper) GetBacking(ctx sdk.Context, id int64) (backing Backing, err sdk.Error) {
	store := ctx.KVStore(k.backingKey)
	key := getBackingIDKey(k, id)
	val := store.Get(key)
	if val == nil {
		return backing, ErrNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinary(val, &backing)

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
	case NativeTokenName:
		// mint category coins from trustake
		return mintFromNativeToken(ctx, k, cat, amount, userAddr)
	default:
		return principal, sdk.ErrInvalidCoins("Invalid backing token")

	}

	return
}

// setBacking stores a `Backing` type in the KVStore
func (k Keeper) setBacking(ctx sdk.Context, backing Backing) {
	store := ctx.KVStore(k.backingKey)
	store.Set(
		getBackingIDKey(k, backing.ID),
		k.GetCodec().MustMarshalBinary(backing))
}

// ============================================================================

// mintFromNativeToken creates category coins by burning trustake
func mintFromNativeToken(
	ctx sdk.Context,
	k Keeper,
	cat cat.Category,
	amount sdk.Coin,
	userAddr sdk.AccAddress) (principal sdk.Coin, err sdk.Error) {

	// naïve implementaion: 1 trustake = 1 category coin
	// https://github.com/TruStory/truchain/issues/21
	conversionRate := sdk.NewDec(1)

	// mint new category coins
	principal = sdk.NewCoin(
		cat.CoinName(),
		sdk.NewDecFromInt(amount.Amount).Mul(conversionRate).RoundInt())

	// burn equivalent trustake
	trustake := sdk.Coins{sdk.NewCoin(NativeTokenName, principal.Amount)}
	if _, _, err := k.bankKeeper.SubtractCoins(ctx, userAddr, trustake); err != nil {
		return principal, err
	}

	return
}

// getInterest calcuates the interest for the backing
func getInterest(
	category cat.Category,
	amount sdk.Coin,
	period time.Duration,
	params Params) sdk.Coin {

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
	coin := sdk.NewCoin(category.CoinName(), interest.RoundInt())

	return coin
}

// getBackingIDKey returns byte array for "backings:id:[ID]"
func getBackingIDKey(k Keeper, id int64) []byte {
	return app.GetIDKey(k.backingKey, id)
}
