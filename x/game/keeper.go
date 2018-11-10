package game

import (
	"time"

	"github.com/cosmos/cosmos-sdk/x/bank"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// ReadKeeper defines a module interface that facilitates read only access to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	Get(ctx sdk.Context, id int64) (game Game, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	Create(
		ctx sdk.Context, storyID int64, creator sdk.AccAddress) (int64, sdk.Error)

	Update(
		ctx sdk.Context, gameID int64, amount sdk.Coin) (int64, sdk.Error)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	queueKey    sdk.StoreKey // queue of games currently active
	storyKeeper story.WriteKeeper
	bankKeeper  bank.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey, queueKey sdk.StoreKey, storyKeeper story.WriteKeeper, bankKeeper bank.Keeper,
	codec *amino.Codec) Keeper {

	return Keeper{app.NewKeeper(codec, storeKey), queueKey, storyKeeper, bankKeeper}
}

// ============================================================================

// Create starts a validation game on a story
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, creator sdk.AccAddress) (int64, sdk.Error) {

	// get the story being challenged
	story, err := k.storyKeeper.GetStory(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// check if a game already exists on story
	if story.GameID != 0 {
		return 0, ErrExists(storyID)
	}

	// create an initial empty challenge pool
	coinName, err := k.storyKeeper.GetCoinName(ctx, storyID)
	if err != nil {
		return 0, err
	}
	emptyPool := sdk.NewCoin(coinName, sdk.ZeroInt())

	// create new game type
	game := Game{
		k.GetNextID(ctx),
		storyID,
		creator,
		ctx.BlockHeader().Time.Add(DefaultParams().Expires),
		time.Time{},
		emptyPool,
		thresholdAmount(story),
		app.NewTimestamp(ctx.BlockHeader()),
	}

	// push game id onto queue that will get checked
	// on each block tick for expired games
	q := queue.NewQueue(k.GetCodec(), k.GetStore(ctx))
	q.Push(game.ID)

	// set game in KVStore
	k.set(ctx, game)

	// update story with gameID
	story.GameID = game.ID
	k.storyKeeper.UpdateStory(ctx, story)

	return game.ID, nil
}

// Get the game for the given id
func (k Keeper) Get(ctx sdk.Context, id int64) (game Game, err sdk.Error) {
	store := k.GetStore(ctx)
	bz := store.Get(k.GetIDKey(id))
	if bz == nil {
		return game, ErrNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinary(bz, &game)

	return
}

// Update the challenge pool
func (k Keeper) Update(
	ctx sdk.Context, gameID int64, amount sdk.Coin) (int64, sdk.Error) {

	game, err := k.Get(ctx, gameID)
	if err != nil {
		return 0, err
	}

	// add amount to challenge pool
	game.Pool = game.Pool.Plus(amount)

	// if threshold is reached, start challenge and allow voting to begin
	if game.Pool.Amount.GT(game.ThresholdAmount) {
		err = k.storyKeeper.StartGame(ctx, game.StoryID)
		if err != nil {
			return 0, err
		}
		game.EndTime = ctx.BlockHeader().Time.Add(DefaultParams().Period)
	}

	// update existing challenge in KVStore
	k.set(ctx, game)

	return game.ID, nil
}

// ============================================================================

// saves the `Game` in the KVStore
func (k Keeper) set(ctx sdk.Context, game Game) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(game.ID),
		k.GetCodec().MustMarshalBinary(game))
}

// ============================================================================

// [Shane] TODO: https://github.com/TruStory/truchain/issues/50
func thresholdAmount(s story.Story) sdk.Int {
	return sdk.NewInt(10)
}
