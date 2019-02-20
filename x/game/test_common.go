package game

import (
	"net/url"
	"time"

	"github.com/TruStory/truchain/x/challenge"
	"github.com/TruStory/truchain/x/vote"

	"github.com/TruStory/truchain/x/backing"
	c "github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/story"
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

func mockDB() (sdk.Context, Keeper, c.Keeper) {
	db := dbm.NewMemDB()

	accKey := sdk.NewKVStoreKey(auth.StoreKey)
	storyKey := sdk.NewKVStoreKey(story.StoreKey)
	storyQueueKey := sdk.NewKVStoreKey(story.QueueStoreKey)
	expiredStoryQueueKey := sdk.NewKVStoreKey(story.ExpiredQueueStoreKey)
	catKey := sdk.NewKVStoreKey(c.StoreKey)
	challengeKey := sdk.NewKVStoreKey(challenge.StoreKey)
	gameKey := sdk.NewKVStoreKey(StoreKey)
	votingStoryQueueKey := sdk.NewKVStoreKey(QueueStoreKey)
	backingKey := sdk.NewKVStoreKey(backing.StoreKey)
	voteKey := sdk.NewKVStoreKey(vote.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(expiredStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(gameKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(votingStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(voteKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	header := abci.Header{Time: time.Now().Add(50 * 24 * time.Hour)}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	pk := params.NewKeeper(codec, paramsKey, transientParamsKey)
	am := auth.NewAccountKeeper(codec, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(am,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	ck := c.NewKeeper(catKey, codec)
	sk := story.NewKeeper(
		storyKey,
		storyQueueKey,
		expiredStoryQueueKey,
		votingStoryQueueKey,
		ck,
		pk.Subspace(story.DefaultParamspace),
		codec)

	story.InitGenesis(ctx, sk, story.DefaultGenesisState())

	backingKeeper := backing.NewKeeper(
		backingKey,
		sk,
		bankKeeper,
		ck,
		codec,
	)

	challengeKeeper := challenge.NewKeeper(
		challengeKey,
		backingKeeper,
		bankKeeper,
		sk,
		codec,
	)

	voteKeeper := vote.NewKeeper(
		voteKey,
		votingStoryQueueKey,
		am,
		backingKeeper,
		challengeKeeper,
		sk,
		bankKeeper,
		codec,
	)

	k := NewKeeper(
		gameKey,
		storyQueueKey,
		sk,
		backingKeeper,
		challengeKeeper,
		voteKeeper,
		bankKeeper,
		codec,
	)

	return ctx, k, ck
}

func createFakeCategory(ctx sdk.Context, ck c.WriteKeeper) c.Category {
	existing, err := ck.GetCategory(ctx, 1)
	if err == nil {
		return existing
	}
	id := ck.Create(ctx, "decentralized exchanges", "trudex", "category for experts in decentralized exchanges")
	cat, _ := ck.GetCategory(ctx, id)
	return cat
}

func createFakeStory(ctx sdk.Context, sk story.WriteKeeper, ck c.WriteKeeper) int64 {
	body := "TruStory can be goverened by it's stakeholders."
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := story.Default
	source := url.URL{}

	storyID, _ := sk.Create(ctx, body, cat.ID, creator, source, storyType)

	return storyID
}
