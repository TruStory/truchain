package challenge

import (
	"crypto/rand"
	"net/url"
	"time"

	"github.com/TruStory/truchain/x/backing"

	params "github.com/TruStory/truchain/parameters"
	c "github.com/TruStory/truchain/x/category"
	game "github.com/TruStory/truchain/x/game"
	"github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdkparams "github.com/cosmos/cosmos-sdk/x/params"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (sdk.Context, Keeper, story.Keeper, c.Keeper, bank.Keeper) {
	db := dbm.NewMemDB()

	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	catKey := sdk.NewKVStoreKey("categories")
	challengeKey := sdk.NewKVStoreKey("challenges")
	gameKey := sdk.NewKVStoreKey("games")
	pendingGameListKey := sdk.NewKVStoreKey("pendingGameList")
	gameQueueKey := sdk.NewKVStoreKey("gameQueue")
	backingKey := sdk.NewKVStoreKey("backings")
	backingListKey := sdk.NewKVStoreKey("backingList")
	paramsKey := sdk.NewKVStoreKey(sdkparams.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(sdkparams.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(gameKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(pendingGameListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(gameQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	// fake block time in the future
	header := abci.Header{Time: time.Now().Add(50 * 24 * time.Hour)}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	RegisterAmino(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	pk := sdkparams.NewKeeper(codec, paramsKey, transientParamsKey)
	am := auth.NewAccountKeeper(codec, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(am,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	ck := c.NewKeeper(catKey, codec)
	sk := story.NewKeeper(storyKey, ck, codec)
	backingKeeper := backing.NewKeeper(
		backingKey,
		backingListKey,
		pendingGameListKey,
		gameQueueKey,
		sk,
		bankKeeper,
		ck,
		codec,
	)
	gameKeeper := game.NewKeeper(gameKey, pendingGameListKey, gameQueueKey, sk, backingKeeper, bankKeeper, codec)

	k := NewKeeper(challengeKey, pendingGameListKey, backingKeeper, bankKeeper, gameKeeper, sk, codec)

	return ctx, k, sk, ck, bankKeeper
}

func createFakeStory(ctx sdk.Context, sk story.Keeper, ck c.WriteKeeper) int64 {
	body := "Body of story."
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := story.Default
	source := url.URL{}
	argument := "fake argument"

	storyID, _ := sk.Create(ctx, argument, body, cat.ID, creator, source, storyType)

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

func fakeFundedCreator(ctx sdk.Context, k bank.Keeper) sdk.AccAddress {
	bz := make([]byte, 4)
	rand.Read(bz)
	creator := sdk.AccAddress(bz)

	// give user some category coins
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(2000000000000))
	k.AddCoins(ctx, creator, sdk.Coins{amount})

	return creator
}

func fakePendingGameQueue() (ctx sdk.Context, k Keeper) {
	ctx, k, storyKeeper, catKeeper, _ := mockDB()

	storyID := createFakeStory(ctx, storyKeeper, catKeeper)
	amount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(1000000000000))
	argument := "test argument"

	creator1 := fakeFundedCreator(ctx, k.bankKeeper)
	creator2 := fakeFundedCreator(ctx, k.bankKeeper)
	creator3 := fakeFundedCreator(ctx, k.bankKeeper)
	creator4 := fakeFundedCreator(ctx, k.bankKeeper)

	// fake backings
	// needed to get a decent challenge threshold
	duration := 1 * time.Hour
	// need to type assert for testing
	// because the backing keeper inside the challenge keeper is read-only
	bk, _ := k.backingKeeper.(backing.WriteKeeper)
	bk.Create(ctx, storyID, amount, argument, creator1, duration)
	bk.Create(ctx, storyID, amount, argument, creator2, duration)
	bk.Create(ctx, storyID, amount, argument, creator3, duration)
	bk.Create(ctx, storyID, amount, argument, creator4, duration)

	// fake challenges
	challengeAmount := sdk.NewCoin(params.StakeDenom, sdk.NewInt(10000000000))

	k.Create(ctx, storyID, challengeAmount, argument, creator1)
	k.Create(ctx, storyID, challengeAmount, argument, creator2)

	return
}
