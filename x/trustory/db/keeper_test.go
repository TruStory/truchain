package db

import (
	"testing"

	ts "github.com/TruStory/trucoin/x/trustory/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/stretchr/testify/assert"
	amino "github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func TestAddGetStory(t *testing.T) {
	ms, storyKey := setupMultiStore()
	cdc := makeCodec()

	keeper := NewStoryKeeper(storyKey, cdc)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	body := "Body of story."
	category := ts.DEX
	creator := sdk.AccAddress([]byte{1, 2})
	storyType := ts.Default

	storyID, err := keeper.AddStory(ctx, body, category, creator, storyType)
	assert.Nil(t, err)

	savedStory, err := keeper.GetStory(ctx, storyID)
	assert.Nil(t, err)

	story := ts.Story{
		ID:           storyID,
		Body:         body,
		Category:     ts.DEX,
		CreatedBlock: int64(0),
		Creator:      creator,
		State:        ts.Created,
		StoryType:    storyType,
	}

	assert.Equal(t, savedStory, story, "Story received from store does not match expected value")
}

// ============================================================================

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	storyKey := sdk.NewKVStoreKey("StoryKey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(storyKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, storyKey
}

func makeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	ts.RegisterAmino(cdc)
	crypto.RegisterAmino(cdc)
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "cosmos-sdk/BaseAccount", nil)
	return cdc
}
