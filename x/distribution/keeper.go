package distribution

import (
	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/story"
	queue "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
)

const (
	// StoreKey is string representation of the store key for distribution
	StoreKey = "distribution"
)

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	expiredStoryQueueKey sdk.StoreKey

	storyKeeper     story.WriteKeeper
	backingKeeper   backing.WriteKeeper
	challengeKeeper challenge.WriteKeeper
	bankKeeper      bank.Keeper
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	expiredStoryQueueKey sdk.StoreKey,
	storyKeeper story.WriteKeeper,
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	bankKeeper bank.Keeper,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		expiredStoryQueueKey,
		storyKeeper,
		backingKeeper,
		challengeKeeper,
		bankKeeper,
	}
}

func (k Keeper) expiredStoryQueue(ctx sdk.Context) queue.Queue {
	store := ctx.KVStore(k.expiredStoryQueueKey)
	return queue.NewQueue(k.GetCodec(), store)
}
