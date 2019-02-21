package game

import (
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
)

// ReadKeeper defines a module interface that facilitates read only access to truchain data
type ReadKeeper interface {
	app.ReadKeeper
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

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

func (k Keeper) storyQueue(ctx sdk.Context) list.Queue {
	queueStore := ctx.KVStore(k.storyQueueKey)
	return list.NewQueue(k.GetCodec(), queueStore)
}
