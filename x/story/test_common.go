package story

import (
	"net/url"
	"time"

	c "github.com/TruStory/truchain/x/category"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper, c.Keeper) {
	db := dbm.NewMemDB()

	storyKey := sdk.NewKVStoreKey("stories")
	storyListKey := sdk.NewKVStoreKey("pendingStoryList")
	expiredStoryQueueKey := sdk.NewKVStoreKey("expiringStoryQueue")
	votingStoryQueueKey := sdk.NewKVStoreKey("challengedStoryQueue")
	catKey := sdk.NewKVStoreKey("categories")
	challengeKey := sdk.NewKVStoreKey("challenges")
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(expiredStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(votingStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	RegisterAmino(codec)

	ck := c.NewKeeper(catKey, codec)
	pk := params.NewKeeper(codec, paramsKey, transientParamsKey)
	sk := NewKeeper(
		storyKey,
		storyListKey,
		ck,
		pk.Subspace(StoreKey),
		codec)
	InitGenesis(ctx, sk, DefaultGenesisState())

	return ctx, sk, ck
}

func createFakeStory(ctx sdk.Context, sk Keeper, ck c.WriteKeeper) int64 {
	body := "TruStory will add what's missing in social networks. Value and truth."

	ctx = ctx.WithBlockHeader(abci.Header{Time: time.Now().UTC()})
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := Default
	source := url.URL{}

	storyID, err := sk.Create(ctx, body, cat.ID, creator, source, storyType)
	if err != nil {
		panic(err)
	}

	return storyID
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

func fakeStories() (ctx sdk.Context, storyKeeper Keeper) {
	ctx, storyKeeper, catKeeper := mockDB()

	createFakeStory(ctx, storyKeeper, catKeeper)
	createFakeStory(ctx, storyKeeper, catKeeper)
	createFakeStory(ctx, storyKeeper, catKeeper)

	return
}
