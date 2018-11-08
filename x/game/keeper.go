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
		ctx sdk.Context, gameID int64) (int64, sdk.Error)

	NewResponseEndBlock(ctx sdk.Context) sdk.Tags

	QueueStore(ctx sdk.Context) sdk.KVStore
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

// Create adds a new challenge on a story
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, creator sdk.AccAddress) (int64, sdk.Error) {

	// get the story being challenged
	story, err := k.storyKeeper.GetStory(ctx, storyID)
	if err != nil {
		return 0, err
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
		false,
		thresholdAmount(story),
		app.NewTimestamp(ctx.BlockHeader()),
	}

	// story.GameID = game.ID
	k.storyKeeper.UpdateStory(ctx, story)

	// push game id onto queue that will get checked
	// on each block tick for expired games
	q := queue.NewQueue(k.GetCodec(), k.GetStore(ctx))
	q.Push(game.ID)

	// set game in KVStore
	k.set(ctx, game)

	return game.ID, nil
}

// Get the game for the given id
func (k Keeper) Get(ctx sdk.Context, id int64) (game Game, err sdk.Error) {
	store := k.GetStore(ctx)
	bz := store.Get(k.GetIDKey(id))
	// if bz == nil {
	// 	return game, ErrNotFound(id)
	// }
	k.GetCodec().MustUnmarshalBinary(bz, &game)

	return
}

// Update mutates an existing challenge, adding a new challenger and updating the pool
func (k Keeper) Update(
	ctx sdk.Context, id int64, creator sdk.AccAddress) (int64, sdk.Error) {

	game, err := k.Get(ctx, id)
	if err != nil {
		return 0, err
	}

	// update existing challenge in KVStore
	k.set(ctx, game)

	return game.ID, nil
}

func (k Keeper) QueueStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.queueKey)
}

// ============================================================================

// Delete removes a challenge from the KVStore
// func (k Keeper) delete(ctx sdk.Context, id int64) sdk.Error {
// 	store := k.GetStore(ctx)
// 	key := k.GetIDKey(id)
// 	bz := store.Get(key)
// 	if bz == nil {
// 		return ErrNotFound(id)
// 	}
// 	store.Delete(key)

// 	return nil
// }

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
