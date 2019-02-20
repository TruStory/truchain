package distribution

import (
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (
	sdk.Context,
	Keeper,
	story.Keeper,
	category.Keeper,
	bank.Keeper,
	auth.AccountKeeper) {

	db := dbm.NewMemDB()

	accKey := sdk.NewKVStoreKey(auth.StoreKey)
	storyKey := sdk.NewKVStoreKey(story.StoreKey)
	storyQueueKey := sdk.NewKVStoreKey(story.QueueStoreKey)
	expiredStoryQueueKey := sdk.NewKVStoreKey(story.ExpiredQueueStoreKey)
	votingStoryQueueKey := sdk.NewKVStoreKey(story.VotingQueueStoreKey)
	catKey := sdk.NewKVStoreKey(category.StoreKey)
	backingKey := sdk.NewKVStoreKey(StoreKey)
	challengeKey := sdk.NewKVStoreKey(challenge.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(expiredStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(votingStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := amino.NewCodec()
	// cryptoAmino.RegisterAmino(codec)
	// RegisterAmino(codec)
	// codec.RegisterInterface((*auth.Account)(nil), nil)
	// codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	categoryKeeper := category.NewKeeper(catKey, codec)

	pk := params.NewKeeper(codec, paramsKey, transientParamsKey)
	am := auth.NewAccountKeeper(codec, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(am,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	storyKeeper := story.NewKeeper(
		storyKey,
		storyQueueKey,
		expiredStoryQueueKey,
		votingStoryQueueKey,
		categoryKeeper,
		pk.Subspace(story.DefaultParamspace),
		codec)

	story.InitGenesis(ctx, storyKeeper, story.DefaultGenesisState())

	backingKeeper := backing.NewKeeper(
		backingKey,
		expiredStoryQueueKey,
		votingStoryQueueKey,
		storyKeeper,
		bankKeeper,
		categoryKeeper,
		codec,
	)

	return ctx, backingKeeper, storyKeeper, categoryKeeper, bankKeeper, am
}
