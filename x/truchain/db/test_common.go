package db

import (
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

// MockDB returns a mock DB to test the app
func MockDB() (sdk.Context, sdk.MultiStore, auth.AccountMapper, TruKeeper) {
	ms, accKey, storyKey, catKey, backingKey := setupMultiStore()
	cdc := makeCodec()
	am := auth.NewAccountMapper(cdc, accKey, auth.ProtoBaseAccount)
	ck := bank.NewBaseKeeper(am)
	k := NewTruKeeper(storyKey, catKey, backingKey, ck, cdc)

	time := time.Now().Add(5 * time.Hour)
	header := abci.Header{Time: time}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	return ctx, ms, am, k
}

// CreateFakeFundedAccount creates a fake funded account for testing
func CreateFakeFundedAccount(ctx sdk.Context, am auth.AccountMapper, coins sdk.Coins) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	baseAcct := auth.NewBaseAccountWithAddress(addr)
	_ = baseAcct.SetCoins(coins)
	am.SetAccount(ctx, &baseAcct)

	return addr
}

// CreateFakeStory creates a fake story for testing
func CreateFakeStory(ms sdk.MultiStore, k TruKeeper) int64 {
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	body := "Body of story."
	cat := CreateFakeCategory(ctx, k)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := ts.Default

	storyID, _ := k.NewStory(ctx, body, cat.ID, creator, storyType)
	return storyID
}

// CreateFakeCategory creates a fake dex category
func CreateFakeCategory(ctx sdk.Context, k TruKeeper) ts.Category {
	id, _ := k.NewCategory(ctx, "decentralized exchanges", "trudex", "category for experts in decentralized exchanges")
	cat, _ := k.GetCategory(ctx, id)
	return cat
}

func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	key := ed25519.GenPrivKey()
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func makeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	ts.RegisterAmino(cdc)
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/BaseAccount", nil)
	return cdc
}

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	categoryKey := sdk.NewKVStoreKey("categories")
	backingKey := sdk.NewKVStoreKey("backings")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(categoryKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, accKey, storyKey, categoryKey, backingKey
}
