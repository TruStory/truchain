package story

import (
	"testing"

	c "github.com/TruStory/truchain/x/category"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/stretchr/testify/assert"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func TestAddGetStory(t *testing.T) {
	ctx, sk, ck := mockDB()

	// test getting a non-existant story
	_, err := sk.GetStory(ctx, int64(5))
	assert.NotNil(t, err)

	storyID := createFakeStory(ctx, sk, ck)

	// test getting an existing story
	savedStory, err := sk.GetStory(ctx, storyID)
	assert.Nil(t, err)

	story := Story{
		ID:           storyID,
		Body:         "Body of story.",
		CategoryID:   int64(1),
		CreatedBlock: int64(0),
		Creator:      sdk.AccAddress([]byte{1, 2}),
		State:        Created,
		Kind:         Default,
	}

	assert.Equal(t, story, savedStory, "Story received from store does not match expected value")

	// test incrementing id by adding another story
	body := "Body of story 2."
	// category := fakeCategory(ctx, k)
	creator := sdk.AccAddress([]byte{3, 4})
	kind := Default

	storyID, _ = sk.NewStory(ctx, body, int64(1), creator, kind)
	assert.Equal(t, int64(2), storyID, "Story ID did not increment properly")
}

// mockDB returns a mock DB to test the module
func mockDB() (sdk.Context, Keeper, c.Keeper) {
	ms, accKey, storyKey, catKey := setupMultiStore()
	cdc := makeCodec()

	am := auth.NewAccountMapper(cdc, accKey, auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(am)

	ck := c.NewKeeper(catKey, storyKey, cdc)
	sk := NewKeeper(storyKey, ck, bk, cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	return ctx, sk, ck
}

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	accKey := sdk.NewKVStoreKey("acc")
	storyKey := sdk.NewKVStoreKey("stories")
	catKey := sdk.NewKVStoreKey("categories")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(accKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(catKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, accKey, storyKey, catKey
}

func makeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
	RegisterAmino(cdc)
	return cdc
}

// createFakeStory creates a fake story for testing
func createFakeStory(ctx sdk.Context, sk Keeper, ck c.ReadWriteKeeper) int64 {
	body := "Body of story."
	cat := createFakeCategory(ctx, ck)
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := Default

	storyID, _ := sk.NewStory(ctx, body, cat.ID, creator, storyType)

	return storyID
}

// createFakeCategory creates a fake dex category
func createFakeCategory(ctx sdk.Context, ck c.ReadWriteKeeper) c.Category {
	id, _ := ck.NewCategory(ctx, "decentralized exchanges", "trudex", "category for experts in decentralized exchanges")
	cat, _ := ck.GetCategory(ctx, id)
	return cat
}
