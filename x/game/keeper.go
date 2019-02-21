package game

import (
	"fmt"
	"time"

	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/vote"

	"github.com/cosmos/cosmos-sdk/x/bank"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/story"
	list "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	amino "github.com/tendermint/go-amino"
)

const (
	// StoreKey is string representation of the store key for games
	StoreKey = "games"
	// QueueStoreKey is string representation of the store key for game queue.
	QueueStoreKey = "gameQueue"
	// PendingListStoreKey is string representation of the store key for the list of pending games.
	PendingListStoreKey = "pendingGameList"
)

// ReadKeeper defines a module interface that facilitates read only access to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	ChallengeThreshold(totalBackingAmount sdk.Coin) sdk.Coin
	Game(ctx sdk.Context, id int64) (game Game, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	Create(ctx sdk.Context, storyID int64, creator sdk.AccAddress) (
		int64, sdk.Error)
	AddToChallengePool(
		ctx sdk.Context, gameID int64, amount sdk.Coin) (err sdk.Error)
	Update(ctx sdk.Context, game Game)
	SetParams(ctx sdk.Context, params Params)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	storyQueueKey sdk.StoreKey

	storyKeeper     story.WriteKeeper
	backingKeeper   backing.WriteKeeper
	challengeKeeper challenge.WriteKeeper
	voteKeeper      vote.WriteKeeper
	bankKeeper      bank.Keeper
	paramStore      params.Subspace
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	storyQueueKey sdk.StoreKey,
	storyKeeper story.WriteKeeper,
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	voteKeeper vote.WriteKeeper,
	bankKeeper bank.Keeper,
	paramStore params.Subspace,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		storyQueueKey,
		storyKeeper,
		backingKeeper,
		challengeKeeper,
		voteKeeper,
		bankKeeper,
		paramStore.WithTypeTable(ParamTypeTable()),
	}
}

// ============================================================================

// Create starts a validation game on a story
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, creator sdk.AccAddress) (int64, sdk.Error) {

	// logger := ctx.Logger().With("module", "game")

	// get the story being challenged
	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// create new game type
	game := Game{
		ID:      k.GetNextID(ctx),
		StoryID: storyID,
		Creator: creator,
		// ChallengeExpireTime: ctx.BlockHeader().Time.Add(DefaultParams().Expires),
		VotingEndTime: time.Time{},
		ChallengePool: sdk.NewCoin(app.StakeDenom, sdk.ZeroInt()),
		Timestamp:     app.NewTimestamp(ctx.BlockHeader()),
	}

	// set game in KVStore
	k.set(ctx, game)

	// update story with gameID
	k.storyKeeper.UpdateStory(ctx, story)

	return game.ID, nil
}

// Game the game for the given id
func (k Keeper) Game(ctx sdk.Context, id int64) (game Game, err sdk.Error) {
	store := k.GetStore(ctx)
	bz := store.Get(k.GetIDKey(id))
	if bz == nil {
		return game, ErrNotFound(id)
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(bz, &game)

	return
}

// AddToChallengePool updates challenge pool and starts game if possible
func (k Keeper) AddToChallengePool(
	ctx sdk.Context, gameID int64, amount sdk.Coin) (err sdk.Error) {

	logger := ctx.Logger().With("module", "game")

	game, err := k.Game(ctx, gameID)
	if err != nil {
		return err
	}

	// add amount to challenge pool
	game.ChallengePool = game.ChallengePool.Plus(amount)
	k.Update(ctx, game)

	msg := "Added %s to challenge pool for game %d"
	logger.Info(fmt.Sprintf(msg, amount.String(), game.ID))

	// get the total of all backings on story
	totalBackingAmount, err := k.backingKeeper.TotalBackingAmount(ctx, game.StoryID)
	if err != nil {
		return err
	}

	threshold := k.ChallengeThreshold(totalBackingAmount)

	// start game if challenge pool is greater than OR equal to challenge threshold
	if game.ChallengePool.IsGTE(threshold) {
		err = k.start(ctx, &game)
		if err != nil {
			return err
		}

		msg := "Challenge threshold met, game started for story %d"
		logger.Info(fmt.Sprintf(msg, game.StoryID))
	}

	return nil
}

// ChallengeThreshold calculates the challenge threshold
func (k Keeper) ChallengeThreshold(totalBackingAmount sdk.Coin) sdk.Coin {
	params := DefaultParams()

	// we have backers
	// calculate challenge threshold amount (based on total backings)
	totalBackingDec := sdk.NewDecFromInt(totalBackingAmount.Amount)
	challengeThresholdAmount := totalBackingDec.
		Mul(params.ChallengeToBackingRatio).
		TruncateInt()

	// challenge threshold can't be less than min challenge stake
	if challengeThresholdAmount.LT(params.MinChallengeThreshold) {
		return sdk.NewCoin(totalBackingAmount.Denom, params.MinChallengeThreshold)
	}

	return sdk.NewCoin(totalBackingAmount.Denom, challengeThresholdAmount)
}

// Update updates the `Game` object
func (k Keeper) Update(ctx sdk.Context, game Game) {

	newGame := Game{
		ID:                  game.ID,
		StoryID:             game.StoryID,
		Creator:             game.Creator,
		ChallengePool:       game.ChallengePool,
		ChallengeExpireTime: game.ChallengeExpireTime,
		VotingEndTime:       game.VotingEndTime,
		Timestamp:           game.Timestamp,
	}

	k.set(ctx, newGame)
}

// ============================================================================

func (k Keeper) storyQueue(ctx sdk.Context) list.Queue {
	queueStore := ctx.KVStore(k.storyQueueKey)
	return list.NewQueue(k.GetCodec(), queueStore)
}

// set saves the `Game` in the KVStore
func (k Keeper) set(ctx sdk.Context, game Game) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(game.ID),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(game))
}

// start registers that a validation game has started
func (k Keeper) start(ctx sdk.Context, game *Game) (err sdk.Error) {
	logger := ctx.Logger().With("module", "game")

	msg := "Challenge threshold met for game %d"
	logger.Info(fmt.Sprintf(msg, game.ID))

	err = k.storyKeeper.StartGame(ctx, game.StoryID)
	if err != nil {
		return
	}

	// set end time = block time + voting period
	game.VotingEndTime = ctx.BlockHeader().Time.Add(DefaultParams().VotingPeriod)

	// update existing game in KVStore
	k.set(ctx, *game)

	// promote game from pending game list to game queue
	// k.removePendingGameList(ctx, game.ID)
	// k.pushGameQueue(ctx, game.ID)

	return
}
