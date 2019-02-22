package voting

import (
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/story"
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

	votingStoryListKey sdk.StoreKey

	accountKeeper   auth.AccountKeeper
	backingKeeper   backing.WriteKeeper
	challengeKeeper challenge.WriteKeeper
	storyKeeper     story.WriteKeeper
	voteKeeper      vote.WriteKeeper
	bankKeeper      bank.Keeper
	paramStore      params.Subspace

	voterList app.UserList
}

// NewKeeper creates a new keeper with write and read access
func NewKeeper(
	storeKey sdk.StoreKey,
	votingStoryListKey sdk.StoreKey,
	accountKeeper auth.AccountKeeper,
	backingKeeper backing.WriteKeeper,
	challengeKeeper challenge.WriteKeeper,
	storyKeeper story.WriteKeeper,
	voteKeeper vote.WriteKeeper,
	bankKeeper bank.Keeper,
	paramStore params.Subspace,
	codec *amino.Codec) Keeper {

	return Keeper{
		app.NewKeeper(codec, storeKey),
		votingStoryListKey,
		accountKeeper,
		backingKeeper,
		challengeKeeper,
		storyKeeper,
		voteKeeper,
		bankKeeper,
		paramStore.WithTypeTable(ParamTypeTable()),
		app.NewUserList(storyKeeper.GetStoreKey()),
	}
}

func (k Keeper) votingStoryList(ctx sdk.Context) list.List {
	store := ctx.KVStore(k.votingStoryListKey)
	return list.NewList(k.GetCodec(), store)
}
