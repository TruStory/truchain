package db

import (
	"math"
	"testing"
	"time"

	ts "github.com/TruStory/truchain/x/truchain/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"

	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)

func Test_key(t *testing.T) {
	bz1 := key("stories", int64(5))
	bz2 := key("stories", int64(math.MaxInt64))

	assert.Equal(t, "stories:5", string(bz1), "should generate valid key")
	assert.Equal(t, "stories:9223372036854775807", string(bz2), "should generate valid key")
}

// ============================================================================
// Helper functions used for keeper tests

func createFakeStory(ms sdk.MultiStore, k TruKeeper) int64 {
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	body := "Body of story."
	category := ts.DEX
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := ts.Default

	storyID, _ := k.NewStory(ctx, body, category, creator, storyType)
	return storyID
}

func createFundedAccount(ctx sdk.Context, am auth.AccountMapper, coins sdk.Coins) sdk.AccAddress {
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

func makeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	ts.RegisterAmino(cdc)
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/BaseAccount", nil)
	return cdc
}

func mockDB() (sdk.Context, sdk.MultiStore, auth.AccountMapper, TruKeeper) {
	ms, accKey, storyKey, backingKey := setupMultiStore()
	cdc := makeCodec()
	am := auth.NewAccountMapper(cdc, accKey, auth.ProtoBaseAccount)
	ck := bank.NewBaseKeeper(am)
	k := NewTruKeeper(storyKey, backingKey, ck, cdc)

	// create fake context with fake block time in header
	// time := time.Date(2018, time.September, 14, 23, 0, 0, 0, time.UTC)
	time := time.Now().Add(5 * time.Hour)
	header := abci.Header{Time: time}
	ctx := sdk.NewContext(ms, header, false, log.NewNopLogger())

	return ctx, ms, am, k
}

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	backingKey := sdk.NewKVStoreKey("backings")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(backingKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, accKey, storyKey, backingKey
}
