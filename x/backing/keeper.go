package backing

import (
	"fmt"
	"time"

	param "github.com/TruStory/truchain/parameters"
	app "github.com/TruStory/truchain/types"
	cat "github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
)

const (
	// StoreKey is string representation of the store key for backings
	StoreKey = "backings"
	// ListStoreKey is string representation of the store key for backing list
	ListStoreKey = "backingList"
)

// ReadKeeper defines a module interface that facilitates read only access
type ReadKeeper interface {
	app.ReadKeeper

	Backing(ctx sdk.Context, id int64) (backing Backing, err sdk.Error)

	BackingByStoryIDAndCreator(
		ctx sdk.Context,
		storyID int64,
		creator sdk.AccAddress) (backing Backing, err sdk.Error)

	BackingsByStoryID(
		ctx sdk.Context, storyID int64) (backings []Backing, err sdk.Error)

	Tally(ctx sdk.Context, storyID int64) (
		trueVotes []Backing, falseVotes []Backing, err sdk.Error)

	TotalBackingAmount(
		ctx sdk.Context, storyID int64) (totalAmount sdk.Coin, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context,
		storyID int64,
		amount sdk.Coin,
		argument string,
		creator sdk.AccAddress,
		duration time.Duration) (int64, sdk.Error)

	Update(ctx sdk.Context, backing Backing)

	ToggleVote(ctx sdk.Context, backingID int64) (int64, sdk.Error)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	storyKeeper    story.WriteKeeper // read-write access to story store
	bankKeeper     bank.Keeper       // read-write access coin store
	categoryKeeper cat.ReadKeeper    // read access to category store

	backingStoryList app.UserList // backings <-> story mappings
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
	argument string,
	creator sdk.AccAddress,
	duration time.Duration) (id int64, err sdk.Error) {

	logger := ctx.Logger().With("module", "backing")

	if amount.Denom != param.StakeDenom {
		return 0, sdk.ErrInvalidCoins("Invalid backing token.")
	}

	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for backing.")
	}

	// check if user has already backed
	if k.backingStoryList.Includes(ctx, k, storyID, creator) {
		return 0, ErrDuplicate(storyID, creator)
	}

	// subtract principal from user
	_, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{amount})
	if err != nil {
		return
	}

	params := DefaultParams()

	credDenom, err := k.storyKeeper.CategoryDenom(ctx, storyID)
	if err != nil {
		return
	}

	// mint category coin from interest earned
	interest := getInterest(
		amount, duration, DefaultMsgParams().MaxPeriod, credDenom, params)

	vote := app.Vote{
		ID:        k.GetNextID(ctx),
		StoryID:   storyID,
		Amount:    amount,
		Argument:  argument,
		Creator:   creator,
		Vote:      true,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	backing := Backing{
		Vote:        vote,
		Interest:    interest,
		MaturesTime: time.Now().Add(duration),
		Params:      params,
		Period:      duration,
	}
	k.setBacking(ctx, backing)

	// add backing <-> story mapping
	k.backingStoryList.Append(ctx, k, storyID, creator, backing.ID())

	logger.Info(fmt.Sprintf(
		"Backed story %d by user %s", storyID, creator.String()))

	return backing.ID(), nil
}

// Update updates an existing backing
func (k Keeper) Update(ctx sdk.Context, backing Backing) {
	newBacking := Backing{
		Vote:        backing.Vote,
		Interest:    backing.Interest,
		MaturesTime: backing.MaturesTime,
		Params:      backing.Params,
		Period:      backing.Period,
	}

	k.setBacking(ctx, newBacking)
}

// ToggleVote changes a true vote to false and vice versa
func (k Keeper) ToggleVote(ctx sdk.Context, backingID int64) (int64, sdk.Error) {
	backing, err := k.Backing(ctx, backingID)
	if err != nil {
		return 0, err
	}

	backing.Vote.Vote = !backing.VoteChoice()
	k.Update(ctx, backing)

	return backingID, nil
}

// Backing gets the backing at the current index from the KVStore
func (k Keeper) Backing(ctx sdk.Context, id int64) (backing Backing, err sdk.Error) {
	store := k.GetStore(ctx)
	key := k.GetIDKey(id)
	val := store.Get(key)
	if val == nil {
		return backing, ErrNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(val, &backing)

	return
}

// BackingsByStoryID returns backings for a given story id
func (k Keeper) BackingsByStoryID(
	ctx sdk.Context, storyID int64) (backings []Backing, err sdk.Error) {

	// iterate over backing list and get backings
	err = k.backingStoryList.Map(ctx, k, storyID, func(backingID int64) sdk.Error {
		backing, err := k.Backing(ctx, backingID)
		if err != nil {
			return err
		}
		backings = append(backings, backing)

		return nil
	})

	return
}

// BackingByStoryIDAndCreator returns backings for a given story id and creator
func (k Keeper) BackingByStoryIDAndCreator(
	ctx sdk.Context,
	storyID int64,
	creator sdk.AccAddress) (backing Backing, err sdk.Error) {

	backingID := k.backingStoryList.Get(ctx, k, storyID, creator)
	backing, err = k.Backing(ctx, backingID)

	return
}

// Tally backings for voting
func (k Keeper) Tally(
	ctx sdk.Context, storyID int64) (
	trueVotes []Backing, falseVotes []Backing, err sdk.Error) {

	err = k.backingStoryList.Map(ctx, k, storyID, func(backingID int64) sdk.Error {
		backing, err := k.Backing(ctx, backingID)
		if err != nil {
			return err
		}
		if backing.VoteChoice() == true {
			trueVotes = append(trueVotes, backing)
		} else {
			falseVotes = append(falseVotes, backing)
		}

		return nil
	})

	return
}

// TotalBackingAmount returns the total of all backings
func (k Keeper) TotalBackingAmount(ctx sdk.Context, storyID int64) (
	totalCoin sdk.Coin, err sdk.Error) {

	totalAmount := sdk.NewCoin(param.StakeDenom, sdk.ZeroInt())

	err = k.backingStoryList.Map(ctx, k, storyID, func(backingID int64) sdk.Error {
		backing, err := k.Backing(ctx, backingID)
		if err != nil {
			return err
		}
		totalAmount = totalAmount.Plus(backing.Amount())

		return nil
	})

	if err != nil {
		return
	}

	return totalAmount, nil
}

// ============================================================================

// setBacking stores a `Backing` type in the KVStore
func (k Keeper) setBacking(ctx sdk.Context, backing Backing) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(backing.ID()),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(backing))
}

// ============================================================================

// getInterest calcuates the interest for the backing
func getInterest(
	amount sdk.Coin,
	period time.Duration,
	maxPeriod time.Duration,
	credDenom string,
	params Params) sdk.Coin {

	// TODO: keep track of total supply
	// https://github.com/TruStory/truchain/issues/22

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
	if interestRate.LT(params.MinInterestRate) {
		interestRate = params.MinInterestRate
	}
	interest := amountDec.Mul(interestRate)

	// return cred coin with rounded interest
	cred := sdk.NewCoin(credDenom, interest.RoundInt())

	return cred
}
