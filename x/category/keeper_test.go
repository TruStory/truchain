package category

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func TestGetCategory_ErrCategoryNotFound(t *testing.T) {
	ctx, k := mockDB()
	id := int64(5)

	_, err := k.GetCategory(ctx, id)
	assert.NotNil(t, err)
	assert.Equal(t, ErrCategoryNotFound(id).Code(), err.Code(), "should get error")
}

func TestGetCategory(t *testing.T) {
	ctx, k := mockDB()

	catID, _ := k.NewCategory(ctx, "dog memes", "doggo", "category for dog memes")
	cat, _ := k.GetCategory(ctx, catID)

	assert.Equal(t, cat.CoinName(), "doggo", "should return coin name")
}

// MockDB returns a mock DB to test the module
func mockDB() (sdk.Context, Keeper) {
	ms, storyKey, catKey := setupMultiStore()
	cdc := makeCodec()
	k := NewKeeper(catKey, storyKey, cdc)
	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	return ctx, k
}

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	storyKey := sdk.NewKVStoreKey("stories")
	categoryKey := sdk.NewKVStoreKey("categories")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(categoryKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, storyKey, categoryKey
}

func makeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	RegisterAmino(cdc)
	return cdc
}
