package expiration

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"time"

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
	votingStoryQueueKey := sdk.NewKVStoreKey(story.VotingQueueStoreKey)
	catKey := sdk.NewKVStoreKey(category.StoreKey)
	backingKey := sdk.NewKVStoreKey(backing.StoreKey)
	challengeKey := sdk.NewKVStoreKey(challenge.StoreKey)
	distKey := sdk.NewKVStoreKey(StoreKey)
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
	ms.MountStoreWithDB(distKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	// RegisterAmino(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	categoryKeeper := category.NewKeeper(catKey, codec)
	category.InitGenesis(ctx, categoryKeeper, category.DefaultCategories())

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
		storyKeeper,
		bankKeeper,
		categoryKeeper,
		codec,
	)

	challengeKeeper := challenge.NewKeeper(
		challengeKey,
		votingStoryQueueKey,
		backingKeeper,
		bankKeeper,
		storyKeeper,
		codec,
	)

	distKeeper := NewKeeper(
		distKey,
		expiredStoryQueueKey,
		storyKeeper,
		backingKeeper,
		challengeKeeper,
		bankKeeper,
		pk.Subspace(DefaultParamspace),
		codec,
	)
	InitGenesis(ctx, distKeeper, DefaultGenesisState())

	return ctx, distKeeper, storyKeeper, backingKeeper, challengeKeeper, bankKeeper
}

func createFakeStory(ctx sdk.Context, sk story.WriteKeeper) int64 {
	body := "TruStory can be goverened by it's stakeholders."
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := story.Default
	source := url.URL{}

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})
	catID := int64(1)
	storyID, err := sk.Create(ctx, body, catID, creator, source, storyType)
	fmt.Println(err)

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
