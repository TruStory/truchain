package backing

import (
	"net/url"

	"github.com/TruStory/truchain/x/argument"
	"github.com/TruStory/truchain/x/category"
	"github.com/TruStory/truchain/x/stake"
	"github.com/TruStory/truchain/x/story"
	"github.com/TruStory/truchain/x/trubank"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
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
	argumentKey := sdk.NewKVStoreKey(argument.StoreKey)
	storyKey := sdk.NewKVStoreKey(story.StoreKey)
	stroyListKey := sdk.NewKVStoreKey(story.QueueStoreKey)
	catKey := sdk.NewKVStoreKey(category.StoreKey)
	backingKey := sdk.NewKVStoreKey(StoreKey)
	pendingGameListKey := sdk.NewKVStoreKey("pendingGameList")
	challengeKey := sdk.NewKVStoreKey("challenges")
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	transientParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)
	truBankKey := sdk.NewKVStoreKey(trubank.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(argumentKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(stroyListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(pendingGameListKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(challengeKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(transientParamsKey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(truBankKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	RegisterAmino(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	ck := category.NewKeeper(catKey, codec)

	pk := params.NewKeeper(codec, paramsKey, transientParamsKey)
	am := auth.NewAccountKeeper(codec, accKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(am,
		pk.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	sk := story.NewKeeper(
		storyKey,
		stroyListKey,
		ck,
		pk.Subspace(story.StoreKey),
		codec)

	story.InitGenesis(ctx, sk, story.DefaultGenesisState())

	truBankKeeper := trubank.NewKeeper(
		truBankKey,
		bankKeeper,
		ck,
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
	argument.InitGenesis(ctx, argumentKeeper, argument.DefaultGenesisState())

	bk := NewKeeper(
		backingKey,
		argumentKeeper,
		stakeKeeper,
		sk,
		bankKeeper,
		truBankKeeper,
		ck,
		codec,
	)

	return ctx, bk, sk, ck, bankKeeper, am
}

func createFakeStory(ctx sdk.Context, sk story.Keeper, ck category.WriteKeeper) int64 {
	body := "TruStory has it's own programmable native currency."
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	storyType := story.Default
	source := url.URL{}

	storyID, _ := sk.Create(ctx, body, cat.ID, creator, source, storyType)

	return storyID
}

func createFakeCategory(ctx sdk.Context, ck category.WriteKeeper) category.Category {
	id := ck.Create(ctx, "decentralized exchanges", "trudex", "category for experts in decentralized exchanges")
	cat, _ := ck.GetCategory(ctx, id)
	return cat
}

func createFakeFundedAccount(ctx sdk.Context, am auth.AccountKeeper, coins sdk.Coins) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	_ = baseAcct.SetCoins(coins)
	am.SetAccount(ctx, &baseAcct)

	return addr
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}
