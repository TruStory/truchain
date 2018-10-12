package backing

import (
	c "github.com/TruStory/truchain/x/category"
	s "github.com/TruStory/truchain/x/story"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func mockDB() (
	sdk.Context,
	Keeper,
	s.Keeper,
	c.Keeper,
	bank.Keeper,
	auth.AccountMapper) {

	db := dbm.NewMemDB()

	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	catKey := sdk.NewKVStoreKey("categories")
	backingKey := sdk.NewKVStoreKey("backings")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	codec := amino.NewCodec()
	cryptoAmino.RegisterAmino(codec)
	RegisterAmino(codec)
	codec.RegisterInterface((*auth.Account)(nil), nil)
	codec.RegisterConcrete(&auth.BaseAccount{}, "auth/Account", nil)

	ck := c.NewKeeper(catKey, storyKey, codec)
	am := auth.NewAccountMapper(codec, accKey, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(am)
	sk := s.NewKeeper(storyKey, ck, codec)
	bk := NewKeeper(backingKey, sk, bankKeeper, ck, codec)

	return ctx, bk, sk, ck, bankKeeper, am
}

func createFakeStory(ctx sdk.Context, sk s.Keeper, ck c.ReadWriteKeeper) int64 {
	body := "Body of story."
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := s.Default

	storyID, _ := sk.NewStory(ctx, body, cat.ID, creator, storyType)

	return storyID
}

func createFakeCategory(ctx sdk.Context, ck c.ReadWriteKeeper) c.Category {
	id, _ := ck.NewCategory(ctx, "decentralized exchanges", "trudex", "category for experts in decentralized exchanges")
	cat, _ := ck.GetCategory(ctx, id)
	return cat
}

func createFakeFundedAccount(ctx sdk.Context, am auth.AccountMapper, coins sdk.Coins) sdk.AccAddress {
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
