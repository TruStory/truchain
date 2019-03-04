package expiration

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	amino "github.com/tendermint/go-amino"
)

const (
	// StoreKey is string representation of the store key for expiration
	StoreKey = "expiration"
)

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	expiredStoryQueueKey sdk.StoreKey

	stakeKeeper     stake.Keeper
	storyKeeper     story.WriteKeeper
	backingKeeper   backing.WriteKeeper
	challengeKeeper challenge.WriteKeeper
	paramStore      params.Subspace
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	expiredStoryQueueKey sdk.StoreKey,
	stakeKeeper stake.Keeper,
	storyKeeper story.WriteKeeper,
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	paramStore params.Subspace,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		expiredStoryQueueKey,
		stakeKeeper,
		storyKeeper,
		backingKeeper,
		challengeKeeper,
		paramStore.WithTypeTable(ParamTypeTable()),
	}
}

func (k Keeper) expiringStoryQueue(ctx sdk.Context) queue.Queue {
	store := ctx.KVStore(k.expiredStoryQueueKey)
	return queue.NewQueue(k.GetCodec(), store)
}
