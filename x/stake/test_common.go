package stake

import (
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/trubank"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper) {
	db := dbm.NewMemDB()

	accKey := sdk.NewKVStoreKey(auth.StoreKey)
	storyKey := sdk.NewKVStoreKey(story.StoreKey)
	storyListKey := sdk.NewKVStoreKey(story.PendingListStoreKey)
	expiredStoryQueueKey := sdk.NewKVStoreKey(story.ExpiringQueueStoreKey)
	votingStoryQueueKey := sdk.NewKVStoreKey(story.ChallengedQueueStoreKey)
	catKey := sdk.NewKVStoreKey(category.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(expiredStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(votingStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	ck := category.NewKeeper(catKey, codec)

	pk := params.NewKeeper(codec, paramsKey, transientParamsKey)
	am := auth.NewAccountKeeper(codec, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(am,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	storyKeeper := story.NewKeeper(
		storyKey,
		storyListKey,
		expiredStoryQueueKey,
		votingStoryQueueKey,
		ck,
		pk.Subspace(story.StoreKey),
		codec)

	story.InitGenesis(ctx, storyKeeper, story.DefaultGenesisState())

	truBankKey := sdk.NewKVStoreKey(trubank.StoreKey)
	ms.MountStoreWithDB(truBankKey, sdk.StoreTypeIAVL, db)

	truBankKeeper := trubank.NewKeeper(
		truBankKey,
		bankKeeper,
		ck,
		codec)

	stakeKeeper := NewKeeper(
		storyKeeper,
		truBankKeeper,
		pk.Subspace(StoreKey),
	)
	InitGenesis(ctx, stakeKeeper, DefaultGenesisState())

	return ctx, stakeKeeper
}
