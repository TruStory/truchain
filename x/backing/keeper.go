package backing

import (
	"fmt"
	"time"

	params "github.com/TruStory/truchain/parameters"
	app "github.com/TruStory/truchain/types"
	cat "github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
	list "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
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

	RemoveFromList(ctx sdk.Context, backingID int64) sdk.Error

	Update(ctx sdk.Context, backing Backing)

	ToggleVote(ctx sdk.Context, backingID int64) (int64, sdk.Error)

	NewResponseEndBlock(ctx sdk.Context) sdk.Tags
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	// list of unmatured backings
	backingListKey sdk.StoreKey
	// list of games in the challenged state
	pendingGameListKey sdk.StoreKey
	// queue of games in the voting state
	gameQueueKey sdk.StoreKey

	storyKeeper    story.WriteKeeper // read-write access to story store
	bankKeeper     bank.Keeper       // read-write access coin store
	categoryKeeper cat.ReadKeeper    // read access to category store

	backingStoryList app.UserList // backings <-> story mappings
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	backingListKey sdk.StoreKey,
	pendingGameListKey sdk.StoreKey,
	gameQueueKey sdk.StoreKey,
	storyKeeper story.WriteKeeper,
	bankKeeper bank.Keeper,
	categoryKeeper cat.ReadKeeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		backingListKey,
		pendingGameListKey,
		gameQueueKey,
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

	// check if user has enough cat coins or trustake to back
	trustake := sdk.NewCoin(params.StakeDenom, amount.Amount)
	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) &&
		!k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{trustake}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for backing.")
	}

	// check if user has already backed
	if k.backingStoryList.Includes(ctx, k, storyID, creator) {
		return 0, ErrDuplicate(storyID, creator)
	}

	// get story value from story id
	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return
	}

	// get category value from category id
	cat, err := k.categoryKeeper.GetCategory(ctx, story.CategoryID)
	if err != nil {
		return
	}

	// set principal, converting from trustake if needed
	principal, err := k.getPrincipal(ctx, cat.Denom(), amount, creator)
	if err != nil {
		return
	}

	// subtract principal from user
	_, _, err = k.bankKeeper.SubtractCoins(ctx, creator, sdk.Coins{principal})
	if err != nil {
		return
	}

	// load default backing parameters
	params := DefaultParams()

	// mint category coin from interest earned
	interest := getInterest(
		cat, amount, duration, DefaultMsgParams().MaxPeriod, params)

	// create new implicit true vote type
	vote := app.Vote{
		ID:        k.GetNextID(ctx),
		Amount:    principal,
		Argument:  argument,
		Creator:   creator,
		Vote:      true,
		Timestamp: app.NewTimestamp(ctx.BlockHeader()),
	}

	// create new backing type with embedded vote
	backing := Backing{
		Vote:        vote,
		StoryID:     storyID,
		Interest:    interest,
		MaturesTime: time.Now().Add(duration),
		Params:      params,
		Period:      duration,
	}

	// store backing
	k.setBacking(ctx, backing)

	// add backing id to the backing list for future processing
	k.backingList(ctx).Push(backing.ID())

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
		StoryID:     backing.StoryID,
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

// RemoveFromList removes a backing from the backing list
func (k Keeper) RemoveFromList(ctx sdk.Context, backingID int64) sdk.Error {
	var ID int64
	var indexToDelete uint64
	var found bool
	backingList := k.backingList(ctx)
	backingList.Iterate(&ID, func(index uint64) bool {
		var tempBackingID int64
		err := backingList.Get(index, &tempBackingID)
		if err != nil {
			panic(err)
		}
		if tempBackingID == backingID {
			indexToDelete = index
			found = true
			return true
		}
		return false
	})

	if found == false {
		return ErrNotFound(backingID)
	}

	backingList.Delete(indexToDelete)

	return nil
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

	totalAmount := sdk.ZeroInt()

	err = k.backingStoryList.Map(ctx, k, storyID, func(backingID int64) sdk.Error {
		backing, err := k.Backing(ctx, backingID)
		if err != nil {
			return err
		}

		// totalAmount = totalAmount.Add(backing.Amount().Amount)
		totalAmount = totalAmount.Add(backing.Amount().Amount)

		return nil
	})
	if err != nil {
		return
	}

	denom, err := k.storyKeeper.CategoryDenom(ctx, storyID)
	if err != nil {
		return
	}

	return sdk.NewCoin(denom, totalAmount), nil
}

// ============================================================================

// getPrincipal calculates the principal, the amount the user gets back
// after the backing matures. Returns a coin.
func (k Keeper) getPrincipal(
	ctx sdk.Context,
	denom string,
	amount sdk.Coin,
	userAddr sdk.AccAddress) (principal sdk.Coin, err sdk.Error) {

	// check which type of coin user wants to back in
	switch amount.Denom {
	case denom:
		// check and return amount if user has enough category coins
		if k.bankKeeper.HasCoins(ctx, userAddr, sdk.Coins{amount}) {
			return amount, nil
		}

	case params.StakeDenom:
		// mint category coins from trustake
		principal, err = app.SwapForCategoryCoin(
			ctx, k.bankKeeper, amount, denom, userAddr)

	default:
		return principal, sdk.ErrInvalidCoins("Invalid backing token")

	}

	return
}

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
	coin := sdk.NewCoin(category.Denom(), interest.RoundInt())

	return coin
}

// ExportState returns the state for a given context
func ExportState() {
	fmt.Println("Backing State")
	backing := Backing{}
	fmt.Printf("%+v\n", backing)
}

func (k Keeper) backingList(ctx sdk.Context) list.List {
	store := ctx.KVStore(k.backingListKey)
	return list.NewList(k.GetCodec(), store)
}

func (k Keeper) pendingGameList(ctx sdk.Context) list.List {
	return list.NewList(
		k.GetCodec(),
		ctx.KVStore(k.pendingGameListKey))
}

func (k Keeper) gameQueue(ctx sdk.Context) list.Queue {
	store := ctx.KVStore(k.gameQueueKey)
	return list.NewQueue(k.GetCodec(), store)
}
