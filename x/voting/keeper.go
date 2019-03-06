package voting

import (
	"fmt"

	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/trubank"
	"github.com/TruStory/truchain/x/vote"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"

	app "github.com/TruStory/truchain/types"
	list "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

const (
	// StoreKey is string representation of the store key for voting
	StoreKey = "voting"
)

// ReadKeeper defines a module interface that facilitates read only access to truchain data
type ReadKeeper interface {
	app.ReadKeeper

	GetParams(ctx sdk.Context) Params
	GetVoteResultsByStoryID(ctx sdk.Context, storyID int64) (results VoteResults, err sdk.Error)
}

// WriteKeeper defines a module interface that facilities write only access to truchain data
type WriteKeeper interface {
	ReadKeeper

	EndBlock(ctx sdk.Context) sdk.Tags
	SetParams(ctx sdk.Context, params Params)
}

// Keeper data type storing keys to the key-value store
type Keeper struct {
	app.Keeper

	votingStoryQueueKey sdk.StoreKey

	accountKeeper   auth.AccountKeeper
	backingKeeper   backing.WriteKeeper
	challengeKeeper challenge.WriteKeeper
	stakeKeeper     stake.Keeper
	storyKeeper     story.WriteKeeper
	voteKeeper      vote.Keeper
	bankKeeper      bank.Keeper
	truBankKeeper   trubank.WriteKeeper
	paramStore      params.Subspace
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	votingStoryListKey sdk.StoreKey,
	accountKeeper auth.AccountKeeper,
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	stakeKeeper stake.Keeper,
	storyKeeper story.WriteKeeper,
	voteKeeper vote.Keeper,
	bankKeeper bank.Keeper,
	truBankKeeper trubank.WriteKeeper,
	paramStore params.Subspace,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		votingStoryListKey,
		accountKeeper,
		backingKeeper,
		challengeKeeper,
		stakeKeeper,
		storyKeeper,
		voteKeeper,
		bankKeeper,
		truBankKeeper,
		paramStore.WithTypeTable(ParamTypeTable()),
	}
}

func (k Keeper) challengedStoryQueue(ctx sdk.Context) list.Queue {
	store := ctx.KVStore(k.votingStoryQueueKey)
	return list.NewQueue(k.GetCodec(), store)
}

// saves a `Vote` in the KVStore
func (k Keeper) set(ctx sdk.Context, results VoteResults) {
	store := k.GetStore(ctx)
	store.Set(
		k.GetIDKey(results.ID),
		k.GetCodec().MustMarshalBinaryLengthPrefixed(results))
}

// GetVoteResultsByStoryID gets the vote results for a story by a storyID
func (k Keeper) GetVoteResultsByStoryID(
	ctx sdk.Context, storyID int64) (results VoteResults, err sdk.Error) {

	store := k.GetStore(ctx)
	val := store.Get(k.GetIDKey(storyID))

	logger := ctx.Logger().With("module", StoreKey)
	logger.Info(fmt.Sprintf("Getting vote results for story %d", storyID))

	if val == nil {
		return results, ErrVoteResultsNotFound(storyID)
	}
	k.GetCodec().MustUnmarshalBinaryLengthPrefixed(val, &results)

	return
}
