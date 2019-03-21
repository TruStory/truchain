package challenge

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"time"

	app "github.com/TruStory/truchain/types"
	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/backing"
	"github.com/TruStory/truchain/x/category"
	c "github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/trubank"
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

func mockDB() (sdk.Context, Keeper, story.Keeper, backing.Keeper, bank.Keeper) {
	db := dbm.NewMemDB()

	accKey := sdk.NewKVStoreKey("acc")
	argumentKey := sdk.NewKVStoreKey(argument.StoreKey)
	catKey := sdk.NewKVStoreKey(category.StoreKey)
	storyKey := sdk.NewKVStoreKey("stories")
	storyListKey := sdk.NewKVStoreKey(story.PendingListStoreKey)
	stakeKey := sdk.NewKVStoreKey(stake.StoreKey)
	truBankKey := sdk.NewKVStoreKey(trubank.StoreKey)
	expiredStoryQueueKey := sdk.NewKVStoreKey(story.ExpiringQueueStoreKey)
	challengeKey := sdk.NewKVStoreKey("challenges")
	votingStoryQueueKey := sdk.NewKVStoreKey("gameQueue")
	backingKey := sdk.NewKVStoreKey("backings")
	paramsKey := sdk.NewKVStoreKey(sdkparams.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(sdkparams.TStoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(argumentKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(truBankKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(stakeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(expiredStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(votingStoryQueueKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
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

	categoryKeeper := c.NewKeeper(catKey, codec)
	c.InitGenesis(ctx, categoryKeeper, c.DefaultCategories())

	sk := story.NewKeeper(
		storyKey,
		storyListKey,
		expiredStoryQueueKey,
		votingStoryQueueKey,
		categoryKeeper,
		pk.Subspace(story.StoreKey),
		codec)
	story.InitGenesis(ctx, sk, story.DefaultGenesisState())

	truBankKeeper := trubank.NewKeeper(
		truBankKey,
		bankKeeper,
		categoryKeeper,
		codec)

	stakeKeeper := stake.NewKeeper(
		sk,
		truBankKeeper,
		pk.Subspace(stake.StoreKey),
	)
	stake.InitGenesis(ctx, stakeKeeper, stake.DefaultGenesisState())

	argumentKeeper := argument.NewKeeper(
		argumentKey,
		sk,
		pk.Subspace(argument.StoreKey),
		codec)

	backingKeeper := backing.NewKeeper(
		backingKey,
		argumentKeeper,
		stakeKeeper,
		sk,
		bankKeeper,
		truBankKeeper,
		categoryKeeper,
		codec,
	)

	k := NewKeeper(
		challengeKey,
		argumentKeeper,
		stakeKeeper,
		backingKeeper,
		truBankKeeper,
		bankKeeper,
		sk,
		pk.Subspace(StoreKey),
		codec)

	InitGenesis(ctx, k, DefaultGenesisState())

	return ctx, k, sk, backingKeeper, bankKeeper
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

	// give user some category coins
	amount := sdk.NewCoin(app.StakeDenom, sdk.NewInt(2000000000000))
	k.AddCoins(ctx, creator, sdk.Coins{amount})

	return creator
}
