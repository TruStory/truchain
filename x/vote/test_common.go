package vote

import (
	"crypto/rand"
	"net/url"
	"time"

	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/stake"

	"github.com/TruStory/truchain/x/challenge"

	app "github.com/TruStory/truchain/types"
	c "github.com/TruStory/truchain/x/category"
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

func mockDB() (sdk.Context, Keeper, c.Keeper) {

	db := dbm.NewMemDB()
	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	storyQueueKey := sdk.NewKVStoreKey(story.QueueStoreKey)
	expiredStoryQueueKey := sdk.NewKVStoreKey(story.ExpiredQueueStoreKey)
	catKey := sdk.NewKVStoreKey("categories")
	challengeKey := sdk.NewKVStoreKey("challenges")
	gameKey := sdk.NewKVStoreKey("games")
	pendingGameListKey := sdk.NewKVStoreKey("pendingGameList")
	votingStoryQueueKey := sdk.NewKVStoreKey("gameQueue")
	voteKey := sdk.NewKVStoreKey("vote")
	backingKey := sdk.NewKVStoreKey("backing")
	paramsKey := sdk.NewKVStoreKey(sdkparams.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(sdkparams.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(expiredStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(gameKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(pendingGameListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(votingStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(voteKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.LoadLatestVersion()

	header := abci.Header{Time: time.Now().Add(50 * 24 * time.Hour)}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	pk := sdkparams.NewKeeper(codec, paramsKey, transientParamsKey)
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
		pk.Subspace(story.StoreKey),
		codec)

	story.InitGenesis(ctx, sk, story.DefaultGenesisState())

	stakeKeeper := stake.NewKeeper(pk.Subspace(stake.StoreKey))
	stake.InitGenesis(ctx, stakeKeeper, stake.DefaultGenesisState())

	backingKeeper := backing.NewKeeper(
		backingKey,
		stakeKeeper,
		sk,
		bankKeeper,
		ck,
		codec,
	)

	challengeKeeper := challenge.NewKeeper(
		challengeKey,
		stakeKeeper,
		backingKeeper,
		bankKeeper,
		sk,
		pk.Subspace(challenge.StoreKey),
		codec,
	)
	challenge.InitGenesis(ctx, challengeKeeper, challenge.DefaultGenesisState())

	k := NewKeeper(
		voteKey,
		votingStoryQueueKey,
		stakeKeeper,
		am,
		backingKeeper,
		challengeKeeper,
		sk,
		bankKeeper,
		pk.Subspace(StoreKey),
		codec)
	InitGenesis(ctx, k, DefaultGenesisState())

	return ctx, k, ck
}

func createFakeStory(ctx sdk.Context, sk story.WriteKeeper, ck c.WriteKeeper) int64 {
	body := "TruStory validators can be bootstrapped with a single genesis file."
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := story.Default
	source := url.URL{}

	storyID, _ := sk.Create(ctx, body, cat.ID, creator, source, storyType)

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

	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(2000000000000))
	k.AddCoins(ctx, creator, sdk.Coins{amount})

	return creator
}
