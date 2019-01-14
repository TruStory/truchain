package game

import (
	"fmt"
	"time"

	"github.com/TruStory/truchain/x/backing"

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

	ChallengeThreshold(totalBackingAmount sdk.Coin) sdk.Coin
	Game(ctx sdk.Context, id int64) (game Game, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	Create(ctx sdk.Context, storyID int64, creator sdk.AccAddress) (
		int64, sdk.Error)
	RegisterChallenge(
		ctx sdk.Context, gameID int64, amount sdk.Coin) (err sdk.Error)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	queueKey       sdk.StoreKey // queue of unexpired active
	activeQueueKey sdk.StoreKey // queue of started games
	storyKeeper    story.WriteKeeper
	backingKeeper  backing.WriteKeeper
	bankKeeper     bank.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	queueKey sdk.StoreKey,
	activeQueueKey sdk.StoreKey,
	storyKeeper story.WriteKeeper,
	backingKeeper backing.WriteKeeper,
	bankKeeper bank.Keeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		queueKey,
		activeQueueKey,
		storyKeeper,
		backingKeeper,
		bankKeeper,
	}
}

// ============================================================================

// Create starts a validation game on a story
func (k Keeper) Create(
	ctx sdk.Context, storyID int64, creator sdk.AccAddress) (int64, sdk.Error) {

	// get the story being challenged
	story, err := k.storyKeeper.Story(ctx, storyID)
	if err != nil {
		return 0, err
	}

	// check if a game already exists on story
	if story.GameID != 0 {
		return 0, ErrExists(storyID)
	}

	// create an initial empty challenge pool
	coinName, err := k.storyKeeper.CategoryDenom(ctx, storyID)
	if err != nil {
		return 0, err
	}

	params := DefaultParams()
	emptyPool := sdk.NewCoin(coinName, sdk.ZeroInt())

	// create new game type
	game := Game{
		ID:                  k.GetNextID(ctx),
		StoryID:             storyID,
		Creator:             creator,
		ExpiresTime:         ctx.BlockHeader().Time.Add(params.Expires),
		VotingPeriodEndTime: time.Time{},
		ChallengePool:       emptyPool,
		Timestamp:           app.NewTimestamp(ctx.BlockHeader()),
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

// RegisterChallenge updates threshold pool and starts game if possible
func (k Keeper) RegisterChallenge(
	ctx sdk.Context, gameID int64, amount sdk.Coin) (err sdk.Error) {

	logger := ctx.Logger().With("module", "x/game")

	game, err := k.Game(ctx, gameID)
	if err != nil {
		return err
	}

	// add amount to threshold pool
	game.ChallengePool = game.ChallengePool.Plus(amount)
	k.update(ctx, game)

	// get the total of all backings on story
	totalBackingAmount, err := k.backingKeeper.TotalBackingAmount(ctx, game.StoryID)
	if err != nil {
		return err
	}

	threshold := k.ChallengeThreshold(totalBackingAmount)

	// start game if challenge pool is greater than OR equal to challenge threshold
	if game.ChallengePool.IsGTE(threshold) {
		err = k.start(ctx, &game)
		logger.Info(fmt.Sprintf(
			"Challenge threshold met, game started for story %d", game.StoryID))
	}
	if err != nil {
		return err
	}

	return nil
}

// ChallengeThreshold calculates the challenge threshold
func (k Keeper) ChallengeThreshold(totalBackingAmount sdk.Coin) sdk.Coin {
	params := DefaultParams()

	// we have backers
	// calculate challenge threshold amount (based on total backings)
	totalBackingDec := sdk.NewDecFromInt(totalBackingAmount.Amount)
	challengeThresholdAmount := totalBackingDec.Mul(params.ChallengeToBackingRatio).RoundInt()

	// challenge threshold can't be less than min challenge stake
	if challengeThresholdAmount.LT(params.MinChallengeStake) {
		return sdk.NewCoin(totalBackingAmount.Denom, params.MinChallengeStake)
	}

	return sdk.NewCoin(totalBackingAmount.Denom, challengeThresholdAmount)
}

// ============================================================================

// set saves the `Game` in the KVStore
func (k Keeper) set(ctx sdk.Context, game Game) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(game.ID),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(game))
}

// start registers that a validation game has started
func (k Keeper) start(ctx sdk.Context, game *Game) (err sdk.Error) {
	err = k.storyKeeper.StartGame(ctx, game.StoryID)
	if err != nil {
		return
	}

	// set end time = block time + voting period
	game.VotingPeriodEndTime = ctx.BlockHeader().Time.Add(DefaultParams().VotingPeriod)
	game.Started = true

	// push game id onto active game queue that will get checked on each tick
	activeQueueStore := ctx.KVStore(k.activeQueueKey)
	q := queue.NewQueue(k.GetCodec(), activeQueueStore)
	q.Push(game.ID)

	// update existing game in KVStore
	k.set(ctx, *game)

	return
}

// update updates the `Game` object
func (k Keeper) update(ctx sdk.Context, game Game) {

	newGame := Game{
		ID:                  game.ID,
		StoryID:             game.StoryID,
		Creator:             game.Creator,
		ExpiresTime:         game.ExpiresTime,
		VotingPeriodEndTime: game.VotingPeriodEndTime,
		ChallengePool:       game.ChallengePool,
		Started:             game.Started,
		Timestamp:           game.Timestamp,
	}

	k.set(ctx, newGame)
}
