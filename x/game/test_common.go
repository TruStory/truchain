package game

import (
	"crypto/rand"
	"net/url"
	"time"

	"github.com/TruStory/truchain/x/category"
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

func mockDB() (
	sdk.Context,
	Keeper,
	story.Keeper,
	backing.Keeper,
	challenge.Keeper,
	bank.Keeper) {

	db := dbm.NewMemDB()

	accKey := sdk.NewKVStoreKey(auth.StoreKey)
	storyKey := sdk.NewKVStoreKey(story.StoreKey)
	storyQueueKey := sdk.NewKVStoreKey(story.QueueStoreKey)
	expiredStoryQueueKey := sdk.NewKVStoreKey(story.ExpiredQueueStoreKey)
	catKey := sdk.NewKVStoreKey(c.StoreKey)
	challengeKey := sdk.NewKVStoreKey(challenge.StoreKey)
	gameKey := sdk.NewKVStoreKey(StoreKey)
	votingStoryQueueKey := sdk.NewKVStoreKey(story.VotingQueueStoreKey)
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
	category.InitGenesis(ctx, ck, category.DefaultCategories())

	sk := story.NewKeeper(
		storyKey,
		storyQueueKey,
		expiredStoryQueueKey,
		votingStoryQueueKey,
		ck,
		pk.Subspace(story.StoreKey),
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
		pk.Subspace(challenge.StoreKey),
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
		pk.Subspace(StoreKey),
		codec,
	)
	InitGenesis(ctx, k, DefaultGenesisState())

	return ctx, k, sk, backingKeeper, challengeKeeper, bankKeeper
}

func createFakeStory(ctx sdk.Context, sk story.WriteKeeper) int64 {
	body := "These bits are going inside a key-value store. Woo hoo!"
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := story.Default
	source := url.URL{}

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})
	catID := int64(1)
	storyID, _ := sk.Create(ctx, body, catID, creator, source, storyType)

	return storyID
}

func fakeFundedCreator(ctx sdk.Context, k bank.Keeper) sdk.AccAddress {
	bz := make([]byte, 4)
	rand.Read(bz)
	creator := sdk.AccAddress(bz)
	amount := sdk.NewCoin("trusteak", sdk.NewInt(2000000000000))
	k.AddCoins(ctx, creator, sdk.Coins{amount})

	return creator
}
