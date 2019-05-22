package backing

import (
	"fmt"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/argument"
	cat "github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	trubank "github.com/TruStory/truchain/x/trubank"
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

	BackersByStoryID(
		ctx sdk.Context, storyID int64) (backers []sdk.AccAddress, err sdk.Error)

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
		argumentID int64,
		argument string,
		creator sdk.AccAddress) (int64, sdk.Error)
	LikeArgument(ctx sdk.Context, argumentID int64, creator sdk.AccAddress, amount sdk.Coin) (*stake.LikeResult, sdk.Error)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	argumentKeeper argument.Keeper
	stakeKeeper    stake.Keeper
	storyKeeper    story.WriteKeeper // read-write access to story store
	bankKeeper     bank.Keeper       // read-write access bank store
	trubankKeeper  trubank.Keeper    // read-write access trubank store
	categoryKeeper cat.ReadKeeper    // read access to category store

	backingStoryList app.UserList // backings <-> story mappings
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	argumentKeeper argument.Keeper,
	stakeKeeper stake.Keeper,
	storyKeeper story.WriteKeeper,
	bankKeeper bank.Keeper,
	trubankKeeper trubank.Keeper,
	categoryKeeper cat.ReadKeeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		argumentKeeper,
		stakeKeeper,
		storyKeeper,
		bankKeeper,
		trubankKeeper,
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
	argumentID int64,
	argument string,
	creator sdk.AccAddress) (id int64, err sdk.Error) {

	logger := ctx.Logger().With("module", StoreKey)

	err = k.stakeKeeper.ValidateAmount(ctx, amount)
	if err != nil {
		return 0, err
	}

	err = k.stakeKeeper.ValidateStoryState(ctx, storyID)
	if err != nil {
		return 0, err
	}

	if amount.Denom != app.StakeDenom {
		return 0, sdk.ErrInvalidCoins("Invalid backing token.")
	}

	if !k.bankKeeper.HasCoins(ctx, creator, sdk.Coins{amount}) {
		return 0, sdk.ErrInsufficientFunds("Insufficient funds for backing.")
	}

	// check if user has already backed
	if k.backingStoryList.Includes(ctx, k, storyID, creator) {
		return 0, ErrDuplicate(storyID, creator)
	}

	stakeID := k.GetNextID(ctx)

	argumentID, err = k.argumentKeeper.Create(ctx, stakeID, storyID, argumentID, argument, creator)
	if err != nil {
		return 0, err
	}

	vote := stake.Vote{
		ID:         stakeID,
		StoryID:    storyID,
		Amount:     amount,
		ArgumentID: argumentID,
		Creator:    creator,
		Vote:       true,
		Timestamp:  app.NewTimestamp(ctx.BlockHeader()),
	}

	backing := Backing{
		Vote: &vote,
	}
	k.setBacking(ctx, backing)

	// add backing <-> story mapping
	k.backingStoryList.Append(ctx, k, storyID, creator, backing.ID())

	// subtract principal from user
	_, err = k.trubankKeeper.SubtractCoin(ctx, creator, amount, storyID, trubank.Backing, backing.ID())
	if err != nil {
		return
	}

	logger.Info(fmt.Sprintf("Backed story %d by user %s", storyID, creator))

	return backing.ID(), nil
}

// LikeArgument likes and argument
func (k Keeper) LikeArgument(ctx sdk.Context, argumentID int64, creator sdk.AccAddress, amount sdk.Coin) (*stake.LikeResult, sdk.Error) {
	arg, err := k.argumentKeeper.Argument(ctx, argumentID)

	if err != nil {
		return nil, err
	}

	err = k.argumentKeeper.RegisterLike(ctx, argumentID, creator)
	if err != nil {
		return nil, err
	}

	backing, err := k.Backing(ctx, arg.StakeID)
	if err != nil {
		return nil, err
	}

	story, err := k.storyKeeper.Story(ctx, backing.StoryID())
	if err != nil {
		return nil, err
	}

	// new backing based on the liked argument that doesn't create a new argument
	backingID, err := k.Create(ctx, story.ID, amount, argumentID, "", creator)
	if err != nil {
		return nil, err
	}

	stakeToCredRatio := k.stakeKeeper.GetParams(ctx).StakeToCredRatio
	likeCredAmount := amount.Amount.Quo(stakeToCredRatio)

	_, err = k.trubankKeeper.MintAndAddCoin(
		ctx,
		backing.Creator(),
		story.CategoryID,
		story.ID,
		trubank.BackingLike,
		arg.StakeID,
		likeCredAmount)
	if err != nil {
		return nil, err
	}
	cat, err := k.categoryKeeper.GetCategory(ctx, story.CategoryID)
	if err != nil {
		return nil, err
	}
	return &stake.LikeResult{
		StakeID:         backingID,
		ArgumentID:      argumentID,
		ArgumentCreator: backing.Creator(),
		CredEarned:      sdk.NewCoin(cat.Denom(), likeCredAmount),
		StoryID:         story.ID,
	}, nil
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

// BackersByStoryID returns a list of addresses that backed a specific story.
func (k Keeper) BackersByStoryID(ctx sdk.Context, storyID int64) (backers []sdk.AccAddress, err sdk.Error) {
	// iterate over backing list and get backings
	err = k.backingStoryList.Map(ctx, k, storyID, func(backingID int64) sdk.Error {
		backing, err := k.Backing(ctx, backingID)
		if err != nil {
			return err
		}
		backers = append(backers, backing.Creator())
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

// TotalBackingAmount returns the total of all backings
func (k Keeper) TotalBackingAmount(ctx sdk.Context, storyID int64) (
	totalCoin sdk.Coin, err sdk.Error) {

	totalAmount := sdk.NewCoin(app.StakeDenom, sdk.ZeroInt())

	err = k.backingStoryList.Map(ctx, k, storyID, func(backingID int64) sdk.Error {
		backing, err := k.Backing(ctx, backingID)
		if err != nil {
			return err
		}
		totalAmount = totalAmount.Add(backing.Amount())

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
